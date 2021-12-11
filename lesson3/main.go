package main

import (
	"context"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var df processAll = &dirFiles{}
	flag.Parse()

	if *logFlag {
		log.SetFormatter(&log.JSONFormatter{})
		f, err := os.OpenFile("program.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Info("Program Started.")

	err := df.scanDir()
	for err != nil {
		log.WithError(err).Warn(`The directory "%s" doesn't exist, please, try again.`)
		//log.Printf(`The directory "%s" doesn't exist, please, try again.`, df.dir)
		err = df.scanDir()
	}

	err = df.walkDir()
	if err != nil {
		log.WithError(err).Fatal("Couldn't walk the directory.")
		//log.Fatal(checkError("Couldn't walk the directory."))
	}

	err = df.findDuplicates()
	if err != nil {
		log.WithError(err).Fatal("Error finding duplicates in the directory.")
		//log.Fatal(checkError("Error finding duplicates in the directory."))
	}

	err = df.deleteDuplicates(ctx)
	for err != nil {
		//log.Println(err)
		err = df.deleteDuplicates(ctx)
	}

	err = df.copyOriginals(ctx)
	for err != nil {
		//log.Println(err)
		err = df.copyOriginals(ctx)
	}
}
