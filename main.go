package main

import (
	"CourseWork/process"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	//SELECT SNo, Province/State FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Country/Region = "Mainland China" AND Confirmed > 100 AND Deaths < 50 AND Recovered > 20
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000 AND Deaths < 500 AND Recovered > 5000
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000

	//Print file path (make config yaml) and commit version
	//Set up timeout (from config)
	//All requests write to access.log, all errors to error.log
	//Write tests with mocks

	//parse seach parameter name that has several words in it

	var r process.Request

	requestBody, err := process.GetRequest()
	if err != nil {
		log.Error(err)
	}
	err = r.ParseRequest(requestBody)
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chSig := make(chan os.Signal, 10)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			sig := <-chSig
			switch sig {
			case syscall.SIGTERM:
				log.Info("Signal SIGTERM caught\n")
				cancel()
			case syscall.SIGINT:
				log.Info("Signal SIGINT caught\n")
				cancel()
			}
		}
	}()

	err = r.ReadFile(ctx)
	if err != nil {
		log.Error(err)
	}

}
