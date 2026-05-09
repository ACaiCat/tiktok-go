package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

type cliConfig struct {
	baseURL       string
	wsURL         string
	mode          string
	username      string
	password      string
	code          string
	token         string
	receiverID    int64
	message       string
	count         int
	concurrency   int
	interval      time.Duration
	timeout       time.Duration
	payloadSize   int
	showIncoming  bool
	compress      bool
	insecureLogin bool
}

type loginResponse struct {
	Base struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"base"`
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type inboundEnvelope struct {
	Message ws.Message
	Raw     []byte
	Err     error
}

type benchResult struct {
	Sent        int64
	Acked       int64
	Failed      int64
	StartedAt   time.Time
	FinishedAt  time.Time
	Durations   []time.Duration
	Min         time.Duration
	Max         time.Duration
	Average     time.Duration
	Throughput  float64
	SuccessRate float64
}

func main() {
	conf := parseFlags()

	if conf.token == "" {
		config.Init()
	}

	token, err := resolveToken(conf)
	if err != nil {
		log.Fatal(err)
	}

	wsURL, err := resolveWSURL(conf)
	if err != nil {
		log.Fatal(err)
	}

	conn, resp, err := dialWebsocket(wsURL, token, conf.compress, conf.insecureLogin)
	if err != nil {
		if resp != nil {
			log.Fatalf("connect websocket failed: %v (http status=%s)", err, resp.Status)
		}
		log.Fatal(err)
	}
	defer conn.Close()

	hlog.Error("connected to %s", wsURL)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	inbound := make(chan inboundEnvelope, 128)
	go readLoop(conn, inbound)

	switch conf.mode {
	case "bench":
		result, err := runBenchmark(conn, inbound, conf, interrupt)
		if err != nil {
			log.Fatal(err)
		}
		printBenchResult(result)
	case "interactive":
		if err := runInteractive(conn, inbound, conf, interrupt); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unsupported mode %q", conf.mode)
	}
}

func parseFlags() cliConfig {
	conf := cliConfig{}
	flag.StringVar(&conf.baseURL, "base-url", "http://127.0.0.1:13215", "HTTP base URL, used for login and default ws URL derivation")
	flag.StringVar(&conf.wsURL, "ws-url", "ws://127.0.0.1:13215/ws", "explicit WebSocket URL, for example ws://127.0.0.1:8080/ws")
	flag.StringVar(&conf.mode, "mode", "interactive", "client mode: interactive or bench")
	flag.StringVar(&conf.username, "username", "", "username used to call POST /user/login when token is empty")
	flag.StringVar(&conf.password, "password", "", "password used to call POST /user/login when token is empty")
	flag.StringVar(&conf.code, "code", "", "optional MFA code for login")
	flag.StringVar(&conf.token, "token", "", "access token used for Authorization: Bearer <token>")
	flag.Int64Var(&conf.receiverID, "receiver-id", 0, "chat target user ID; required for sending chat messages and benchmark mode")
	flag.StringVar(&conf.message, "message", "ping", "message prefix in benchmark mode or default chat content in interactive mode")
	flag.IntVar(&conf.count, "count", 100, "number of benchmark messages to send")
	flag.IntVar(&conf.concurrency, "concurrency", 1, "number of in-flight benchmark messages on one websocket connection")
	flag.DurationVar(&conf.interval, "interval", 0, "delay between benchmark sends per worker")
	flag.DurationVar(&conf.timeout, "timeout", 5*time.Second, "timeout for benchmark message echo and login HTTP requests")
	flag.IntVar(&conf.payloadSize, "payload-size", 0, "extra payload bytes appended to each benchmark message")
	flag.BoolVar(&conf.showIncoming, "show-incoming", true, "print incoming messages in interactive mode")
	flag.BoolVar(&conf.compress, "compress", true, "request websocket per-message compression")
	flag.BoolVar(&conf.insecureLogin, "insecure-login", false, "skip TLS verification for HTTPS login and WSS when using self-signed certificates")
	flag.Parse()

	if conf.count < 1 {
		conf.count = 1
	}
	if conf.concurrency < 1 {
		conf.concurrency = 1
	}
	if conf.timeout <= 0 {
		conf.timeout = 5 * time.Second
	}

	return conf
}

func resolveToken(conf cliConfig) (string, error) {
	if conf.token != "" {
		return conf.token, nil
	}
	if conf.username == "" || conf.password == "" {
		return "", errors.New("either -token or both -username and -password are required")
	}
	return login(conf)
}

func login(conf cliConfig) (string, error) {
	endpoint := strings.TrimRight(conf.baseURL, "/") + "/user/login"
	payload := map[string]string{
		"username": conf.username,
		"password": conf.password,
	}
	if conf.code != "" {
		payload["code"] = conf.code
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal login body: %w", err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if conf.insecureLogin {
		transport.TLSClientConfig = insecureTLSConfig()
	}

	client := &http.Client{
		Timeout:   conf.timeout,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read login response: %w", err)
	}

	accessToken := resp.Header.Get(constants.AccessTokenHeader)
	if accessToken == "" {
		return "", fmt.Errorf("login succeeded without %s header, body=%s", constants.AccessTokenHeader, strings.TrimSpace(string(respBody)))
	}

	var parsed loginResponse
	if len(respBody) > 0 {
		_ = json.Unmarshal(respBody, &parsed)
	}
	if parsed.Base.Code != 0 && parsed.Base.Code != 10000 {
		return "", fmt.Errorf("login returned code=%d msg=%s", parsed.Base.Code, parsed.Base.Msg)
	}

	hlog.Error("login ok, user_id=%s", parsed.Data.ID)
	return accessToken, nil
}

func resolveWSURL(conf cliConfig) (string, error) {
	if conf.wsURL != "" {
		return conf.wsURL, nil
	}

	parsed, err := url.Parse(conf.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}

	switch parsed.Scheme {
	case "http":
		parsed.Scheme = "ws"
	case "https":
		parsed.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported base URL scheme %q", parsed.Scheme)
	}

	parsed.Path = "/ws"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func dialWebsocket(wsURL, token string, compress bool, insecure bool) (*websocket.Conn, *http.Response, error) {
	dialer := *websocket.DefaultDialer
	dialer.EnableCompression = compress
	if insecure {
		dialer.TLSClientConfig = insecureTLSConfig()
	}

	headers := http.Header{}
	headers.Set(constants.AuthHeader, "Bearer "+token)

	conn, resp, err := dialer.Dial(wsURL, headers)
	if err != nil {
		return nil, resp, err
	}
	return conn, resp, nil
}

func readLoop(conn *websocket.Conn, inbound chan<- inboundEnvelope) {
	defer close(inbound)
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			inbound <- inboundEnvelope{Err: err}
			return
		}

		var message ws.Message
		if err := json.Unmarshal(data, &message); err != nil {
			inbound <- inboundEnvelope{Raw: data, Err: fmt.Errorf("decode ws envelope: %w", err)}
			continue
		}

		inbound <- inboundEnvelope{Message: message, Raw: data}
	}
}

func runInteractive(conn *websocket.Conn, inbound <-chan inboundEnvelope, conf cliConfig, interrupt <-chan os.Signal) error {
	log.Println("interactive mode commands: plain text sends chat, /history <user_id> [page] [page_size], /unread <user_id>, /raw <json>, /quit")

	stdinDone := make(chan error, 1)
	go func() {
		stdinDone <- interactiveInputLoop(conn, conf)
	}()

	for {
		select {
		case sig := <-interrupt:
			return fmt.Errorf("received signal %s", sig)
		case err := <-stdinDone:
			if err == nil || errors.Is(err, io.EOF) {
				return nil
			}
			return err
		case item, ok := <-inbound:
			if !ok {
				return nil
			}
			if item.Err != nil {
				if websocket.IsCloseError(item.Err, websocket.CloseNormalClosure) || errors.Is(item.Err, io.EOF) {
					return nil
				}
				return item.Err
			}
			if conf.showIncoming {
				printIncoming(item)
			}
		}
	}
}

func interactiveInputLoop(conn *websocket.Conn, conf cliConfig) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case line == "/quit" || line == "/exit":
			return nil
		case strings.HasPrefix(line, "/history "):
			message, err := buildHistoryCommand(line)
			if err != nil {
				log.Println(err)
				continue
			}
			if err := writeJSON(conn, message); err != nil {
				return err
			}
		case strings.HasPrefix(line, "/unread "):
			message, err := buildUnreadCommand(line)
			if err != nil {
				log.Println(err)
				continue
			}
			if err := writeJSON(conn, message); err != nil {
				return err
			}
		case strings.HasPrefix(line, "/raw "):
			raw := strings.TrimSpace(strings.TrimPrefix(line, "/raw "))
			if raw == "" {
				log.Println("/raw requires a JSON payload")
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(raw)); err != nil {
				return err
			}
		default:
			if conf.receiverID == 0 {
				log.Println("plain text sending requires -receiver-id")
				continue
			}
			message := ws.SendMessage{
				Type: ws.MessageTypeChat,
				Body: ws.ChatMessage{
					ReceiverID: conf.receiverID,
					Content:    line,
					Timestamp:  time.Now().Unix(),
				},
			}
			if err := writeJSON(conn, message); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func buildHistoryCommand(line string) (ws.SendMessage, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ws.SendMessage{}, errors.New("usage: /history <user_id> [page] [page_size]")
	}
	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return ws.SendMessage{}, fmt.Errorf("invalid history user_id: %w", err)
	}
	page := 1
	pageSize := 20
	if len(parts) > 2 {
		page, err = strconv.Atoi(parts[2])
		if err != nil {
			return ws.SendMessage{}, fmt.Errorf("invalid page: %w", err)
		}
	}
	if len(parts) > 3 {
		pageSize, err = strconv.Atoi(parts[3])
		if err != nil {
			return ws.SendMessage{}, fmt.Errorf("invalid page_size: %w", err)
		}
	}

	return ws.SendMessage{
		Type: ws.MessageTypeHistory,
		Body: ws.HistoryRequest{Sender: userID, Page: page, PageSize: pageSize},
	}, nil
}

func buildUnreadCommand(line string) (ws.SendMessage, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return ws.SendMessage{}, errors.New("usage: /unread <user_id>")
	}
	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return ws.SendMessage{}, fmt.Errorf("invalid unread user_id: %w", err)
	}

	return ws.SendMessage{
		Type: ws.MessageTypeUnread,
		Body: ws.UnreadRequest{Sender: userID},
	}, nil
}

func runBenchmark(conn *websocket.Conn, inbound <-chan inboundEnvelope, conf cliConfig, interrupt <-chan os.Signal) (*benchResult, error) {
	if conf.receiverID == 0 {
		return nil, errors.New("benchmark mode requires -receiver-id because the server validates friendship before echoing")
	}

	pending := sync.Map{}
	durations := make([]time.Duration, 0, conf.count)
	var durationsMu sync.Mutex
	var sentCount atomic.Int64
	var ackedCount atomic.Int64
	var failedCount atomic.Int64
	var sendMu sync.Mutex

	readerDone := make(chan error, 1)
	go func() {
		for item := range inbound {
			if item.Err != nil {
				readerDone <- item.Err
				return
			}

			if item.Message.Type != ws.MessageTypeChat {
				continue
			}

			var body ws.ChatMessage
			if err := json.Unmarshal(item.Message.Body, &body); err != nil {
				continue
			}

			value, ok := pending.LoadAndDelete(body.Content)
			if !ok {
				continue
			}

			started := value.(time.Time)
			latency := time.Since(started)
			durationsMu.Lock()
			durations = append(durations, latency)
			durationsMu.Unlock()
			ackedCount.Add(1)
		}
		readerDone <- nil
	}()

	jobs := make(chan int)
	workerErr := make(chan error, 1)
	var workerWG sync.WaitGroup

	startedAt := time.Now()
	for worker := 0; worker < conf.concurrency; worker++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			for idx := range jobs {
				content := buildBenchContent(conf, idx)
				message := ws.SendMessage{
					Type: ws.MessageTypeChat,
					Body: ws.ChatMessage{
						ReceiverID: conf.receiverID,
						Content:    content,
						Timestamp:  time.Now().Unix(),
					},
				}

				pending.Store(content, time.Now())
				sendMu.Lock()
				err := writeJSON(conn, message)
				sendMu.Unlock()
				if err != nil {
					pending.Delete(content)
					select {
					case workerErr <- err:
					default:
					}
					return
				}
				sentCount.Add(1)

				deadline := time.NewTimer(conf.timeout)
				acked := false
				for !acked {
					select {
					case <-deadline.C:
						pending.Delete(content)
						failedCount.Add(1)
						acked = true
						break
					default:
						if _, ok := pending.Load(content); !ok {
							acked = true
							break
						}
						time.Sleep(2 * time.Millisecond)
					}
				}
				if !deadline.Stop() {
					select {
					case <-deadline.C:
					default:
					}
				}

				if conf.interval > 0 {
					time.Sleep(conf.interval)
				}
			}
		}()
	}

	for idx := 1; idx <= conf.count; idx++ {
		select {
		case sig := <-interrupt:
			close(jobs)
			workerWG.Wait()
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, sig.String()), time.Now().Add(time.Second))
			return nil, fmt.Errorf("received signal %s", sig)
		case err := <-workerErr:
			close(jobs)
			workerWG.Wait()
			return nil, err
		default:
			jobs <- idx
		}
	}
	close(jobs)
	workerWG.Wait()

	select {
	case err := <-workerErr:
		if err != nil {
			return nil, err
		}
	default:
	}

	if err := drainPendingAcks(&pending, ackedCount.Load(), int64(conf.count), conf.timeout, interrupt); err != nil {
		return nil, err
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	select {
	case err := <-readerDone:
		if err != nil && !isTimeout(err) && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return nil, err
		}
	default:
	}

	finishedAt := time.Now()
	durationsMu.Lock()
	finalDurations := append([]time.Duration(nil), durations...)
	durationsMu.Unlock()

	result := summarizeBench(sentCount.Load(), ackedCount.Load(), failedCount.Load(), startedAt, finishedAt, finalDurations)
	return result, nil
}

func drainPendingAcks(pending *sync.Map, ackedCount int64, expected int64, timeout time.Duration, interrupt <-chan os.Signal) error {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	for {
		if ackedCount >= expected {
			return nil
		}

		empty := true
		pending.Range(func(_, _ any) bool {
			empty = false
			return false
		})
		if empty {
			return nil
		}

		select {
		case sig := <-interrupt:
			return fmt.Errorf("received signal %s", sig)
		case <-deadline.C:
			return nil
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func summarizeBench(sent, acked, failed int64, startedAt, finishedAt time.Time, durations []time.Duration) *benchResult {
	result := &benchResult{
		Sent:       sent,
		Acked:      acked,
		Failed:     failed,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		Durations:  durations,
	}

	if len(durations) == 0 {
		return result
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	result.Min = durations[0]
	result.Max = durations[len(durations)-1]

	var total time.Duration
	for _, duration := range durations {
		total += duration
	}
	result.Average = total / time.Duration(len(durations))

	elapsed := finishedAt.Sub(startedAt).Seconds()
	if elapsed > 0 {
		result.Throughput = float64(acked) / elapsed
	}
	if sent > 0 {
		result.SuccessRate = float64(acked) / float64(sent) * 100
	}

	return result
}

func printBenchResult(result *benchResult) {
	fmt.Printf("sent=%d acked=%d failed=%d success_rate=%.2f%% elapsed=%s throughput=%.2f msg/s\n",
		result.Sent,
		result.Acked,
		result.Failed,
		result.SuccessRate,
		result.FinishedAt.Sub(result.StartedAt).Round(time.Millisecond),
		result.Throughput,
	)

	if len(result.Durations) == 0 {
		fmt.Println("no chat echoes received")
		return
	}

	fmt.Printf("latency min=%s avg=%s p50=%s p95=%s p99=%s max=%s\n",
		result.Min.Round(time.Microsecond),
		result.Average.Round(time.Microsecond),
		percentile(result.Durations, 0.50).Round(time.Microsecond),
		percentile(result.Durations, 0.95).Round(time.Microsecond),
		percentile(result.Durations, 0.99).Round(time.Microsecond),
		result.Max.Round(time.Microsecond),
	)
}

func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	if p <= 0 {
		return durations[0]
	}
	if p >= 1 {
		return durations[len(durations)-1]
	}
	idx := int(math.Ceil(float64(len(durations))*p)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(durations) {
		idx = len(durations) - 1
	}
	return durations[idx]
}

func printIncoming(item inboundEnvelope) {
	switch item.Message.Type {
	case ws.MessageTypeChat:
		var body ws.ChatMessage
		if err := json.Unmarshal(item.Message.Body, &body); err == nil {
			fmt.Printf("[chat] from=%d to=%d ai=%t ts=%d content=%s\n", body.SenderID, body.ReceiverID, body.IsAI, body.Timestamp, body.Content)
			return
		}
	case ws.MessageTypeHistory:
		var body ws.HistoryMessage
		if err := json.Unmarshal(item.Message.Body, &body); err == nil {
			pretty, _ := json.MarshalIndent(body, "", "  ")
			fmt.Printf("[history] %s\n", pretty)
			return
		}
	case ws.MessageTypeUnread:
		var body ws.UnreadMessage
		if err := json.Unmarshal(item.Message.Body, &body); err == nil {
			pretty, _ := json.MarshalIndent(body, "", "  ")
			fmt.Printf("[unread] %s\n", pretty)
			return
		}
	case ws.MessageTypeError:
		var body ws.ErrorMessage
		if err := json.Unmarshal(item.Message.Body, &body); err == nil {
			fmt.Printf("[error] code=%d message=%s\n", body.Code, body.Message)
			return
		}
	}

	fmt.Printf("[raw] %s\n", strings.TrimSpace(string(item.Raw)))
}

func buildBenchContent(conf cliConfig, idx int) string {
	content := fmt.Sprintf("%s-%d-%d", conf.message, idx, time.Now().UnixNano())
	if conf.payloadSize <= 0 {
		return content
	}
	return content + "-" + strings.Repeat("x", conf.payloadSize)
}

func writeJSON(conn *websocket.Conn, payload any) error {
	return conn.WriteJSON(payload)
}

func isTimeout(err error) bool {
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}

func insecureTLSConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}
