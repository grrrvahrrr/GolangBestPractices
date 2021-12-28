package main

import (
	"CourseWork/config"
	"CourseWork/process"
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	//SELECT SNo, Province/State FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Country/Region = "Mainland China" AND Confirmed > 100 AND Deaths < 50 AND Recovered > 20
	//SELECT SNo, Country/Region FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000 AND Deaths < 500 AND Recovered > 5000
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000
	//SELECT SNo, Province/State FROM default WHERE Country/Region = "Mainland China" AND Confirmed > 100 AND Deaths < 50 AND Recovered > 20

	//Write tests with mocks

	//parse seach parameter name that has several words in it

	//Logging
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile("logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Load Config
	path, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}
	cfg, err := config.LoadConfig(path + "/config/config.env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//Get Request
	requestBody, err := process.GetRequest(cfg.DefaultFileName)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}

	//Process request
	var p process.Processer = &process.Request{}

	err = p.ParseRequest(requestBody, cfg.DefaultFileName)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}

	timeoutSec, err := strconv.Atoi(cfg.TimeoutSeconds)
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	chSig := make(chan os.Signal, 10)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			sig := <-chSig
			switch sig {
			case syscall.SIGTERM:
				log.Info("Signal SIGTERM caught")
				cancel()
			case syscall.SIGINT:
				log.Info("Signal SIGINT caught")
				cancel()
			}
		}
	}()

	err = p.ReadFile(ctx)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}

}
