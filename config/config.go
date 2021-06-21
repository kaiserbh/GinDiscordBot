package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	log.Info("Reading from config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Error("reading config file: ", err)
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Unmarshal config file: ", err)
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix

	log.Info("Read successful.")
	return nil
}
