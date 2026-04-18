package main

import (
	"log"

	video "github.com/ACaiCat/tiktok-go/kitex_gen/video/videoservice"
)

func main() {
	svr := video.NewServer(new(VideoServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
