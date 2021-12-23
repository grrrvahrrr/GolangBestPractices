package main

import (
	"CourseWork/process"
)

func main() {
	//SELECT SNo FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Country/Region = Hong Kong AND Confirmed > 10000 AND Deaths < 500 AND Recovered > 5000

	//Print file path (make config yaml) and commit version
	//Process SIGINT with graceful shutdown
	//Set up timeout (from config)
	//All requests write to access.log, all errors to error.log
	//Write tests with mocks
	//make files
	//parse seach parameter name that has several words in it

	var r process.UserRequest

	r.GetRequest()
	r.ParseRequest()
	r.ReadFile()

}
