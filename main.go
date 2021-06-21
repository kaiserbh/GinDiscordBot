package main

import (
	"github.com/kaiserbh/gin-bot-go/bot"
	"github.com/kaiserbh/gin-bot-go/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}
	bot.Start()

	<-make(chan struct{})
}
