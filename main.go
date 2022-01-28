package main

import (
	"CourseWork/config"
	"CourseWork/process"
	"CourseWork/request"
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
	//Request example
	//SELECT SNo, Province/State FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Country/Region = "Mainland China" AND Confirmed > 100 AND Deaths < 50 AND Recovered > 20

	//ToDo
	//parse seach parameter name that has several words in it
	//solve parsing problem with SELECT FROM WHERE in parameter names

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
	requestBody, err := request.GetRequest(cfg.DefaultFileName)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}

	err = request.LogRequest(requestBody)
	if err != nil {
		log.Error(err)
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

	err = p.ReadFile(ctx, os.Stdout)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}

}
