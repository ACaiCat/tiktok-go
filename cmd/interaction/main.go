package main

import (
	"log"

	interaction "github.com/ACaiCat/tiktok-go/kitex_gen/interaction/interactionservice"
)

func main() {
	svr := interaction.NewServer(new(InteractionServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
