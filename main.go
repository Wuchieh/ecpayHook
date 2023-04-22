package main

import (
	"ecpayHook/database"
	_ "ecpayHook/redis"
	"ecpayHook/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sc := make(chan os.Signal, 1)
	if err := database.DatabaseInit(); err != nil {
		panic(err)
	}
	go func() {
		err := server.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
