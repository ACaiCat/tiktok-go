package main

import (
	"log"

	social "github.com/ACaiCat/tiktok-go/kitex_gen/social/socialservice"
)

func main() {
	svr := social.NewServer(new(SocialServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
