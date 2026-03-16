FROM golang:1.26-alpine AS builder

WORKDIR /app
ENV TZ=Asia/Shanghai

ENV GOPROXY="https://goproxy.cn"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app
ENV TZ=Asia/Shanghai
COPY --from=builder /app/main .
RUN apk update --no-cache && apk --no-cache add ca-certificates tzdata ffmpeg

CMD ["./main"]