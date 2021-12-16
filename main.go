package main

import (
	"CourseWork/process"
)

func main() {
	//Print file path (make config yaml) and commit version
	//Process SIGINT with graceful shutdown
	//Set up timeout (from config)
	//All requests write to access.log, all errors to error.log
	//Write tests with mocks
	//golangci lint
	//make files

	var r process.UserRequest

	r.GetRequest()
	r.ParseRequest()
	r.ReadFile()

}
