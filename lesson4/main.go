package main

import (
	"GolangBP/lesson4/process"
	"context"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var df process.ProcessAll = &process.DirFiles{}
	flag.Parse()

	if *process.LogFlag {
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

	err := df.ScanDir()
	for err != nil {
		log.WithError(err).Warn(`The directory doesn't exist, please, try again.`)
		err = df.ScanDir()
	}

	err = df.WalkDir()
	if err != nil {
		log.WithError(err).Fatal("Couldn't walk the directory.")
	}

	err = df.FindDuplicates()
	if err != nil {
		log.WithError(err).Fatal("Error finding duplicates in the directory.")
	}

	err = df.DeleteDuplicates(ctx)
	for err != nil {
		err = df.DeleteDuplicates(ctx)
	}

	err = df.CopyOriginals(ctx)
	for err != nil {
		err = df.CopyOriginals(ctx)
	}
}
