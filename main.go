package main

import (
	"CourseWork/process"

	log "github.com/sirupsen/logrus"
)

func main() {
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Country/Region = Hong Kong AND Confirmed > 10000 AND Deaths < 500 AND Recovered > 5000
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000 AND Deaths < 500 AND Recovered > 5000
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 10000

	//Print file path (make config yaml) and commit version
	//Process SIGINT with graceful shutdown
	//Set up timeout (from config)
	//All requests write to access.log, all errors to error.log
	//Write tests with mocks

	//parse seach parameter name that has several words in it

	var r process.UserRequest

	err := r.GetRequest()
	if err != nil {
		log.Error(err)
	}
	err = r.ParseRequest()
	if err != nil {
		log.Error(err)
	}
	err = r.ReadFile()
	if err != nil {
		log.Error(err)
	}

}
