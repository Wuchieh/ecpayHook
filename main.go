package main

import (
	_ "ecpayHook/redis"
	"ecpayHook/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sc := make(chan os.Signal, 1)
	go func() {
		err := server.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
