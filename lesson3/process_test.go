package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestScanDir(t *testing.T) {
	var df dirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	input := []byte(path + "/testDir")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	err = df.scanDir()
	if err != nil {
		t.Errorf(`The directory "%s" doesn't exist.`, df.dir)
	}
}

func TestWalkDir(t *testing.T) {
	var df dirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	df.dir = path + "/testDir"
	err = df.walkDir()
	if err != nil {
		t.Error("Error walking directory")
	}
	filesExpected := []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
		path + "/testDir/testDir1/file1",
		path + "/testDir/testDir1/file2",
		path + "/testDir/testDir1/testDir3/file1",
		path + "/testDir/testDir1/testDir3/file2",
		path + "/testDir/testDir2/file1",
		path + "/testDir/testDir2/file2",
		path + "/testDir/testDir2/testDir4/file1",
		path + "/testDir/testDir2/testDir4/file2",
	}
	for i, v := range df.files {
		if v != filesExpected[i] {
			t.Errorf("In files slice expected %s, but got %s", filesExpected[i], v)
		}
	}
}

func TestFindDuplicates(t *testing.T) {
	var df dirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	df.dir = path + "/testDir"
	df.files = []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
		path + "/testDir/testDir1/file1",
		path + "/testDir/testDir1/file2",
		path + "/testDir/testDir1/testDir3/file1",
		path + "/testDir/testDir1/testDir3/file2",
		path + "/testDir/testDir2/file1",
		path + "/testDir/testDir2/file2",
		path + "/testDir/testDir2/testDir4/file1",
		path + "/testDir/testDir2/testDir4/file2",
	}

	err = df.findDuplicates()
	if err != nil {
		t.Error("Error finding duplicates in the directory.")
	}
	duplicatesExpected := []string{
		path + `/testDir/testDir1/file1`,
		path + `/testDir/testDir1/testDir3/file1`,
		path + `/testDir/testDir2/testDir4/file1`,
		path + `/testDir/testDir2/file1`,
		path + `/testDir/testDir2/testDir4/file2`,
		path + `/testDir/testDir1/testDir3/file2`,
		path + `/testDir/testDir2/file2`,
		path + `/testDir/testDir1/file2`,
	}
	for i, v := range df.duplicates {
		if v != duplicatesExpected[i] {
			t.Errorf("In duplicates slices expected %s, but got %s", duplicatesExpected[i], v)
		}
	}

	filesExpected := []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
	}
	for i, v := range df.files {
		if v != filesExpected[i] {
			t.Errorf("In files slice expected %s, but got %s", filesExpected[i], v)
		}
	}
}

func TestDeleteDuplicates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var df dirFiles
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	df.dir = path + "/testDir"
	df.files = []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
	}
	df.duplicates = []string{
		path + `/testDir/testDir1/file1`,
		path + `/testDir/testDir1/testDir3/file1`,
		path + `/testDir/testDir2/testDir4/file1`,
		path + `/testDir/testDir2/file1`,
		path + `/testDir/testDir2/testDir4/file2`,
		path + `/testDir/testDir1/testDir3/file2`,
		path + `/testDir/testDir2/file2`,
		path + `/testDir/testDir1/file2`,
	}
	*delFlag = true

	err = df.deleteDuplicates(ctx)
	if err != nil {
		t.Error("Couldn't delete duplicates.")
	}
	for _, v := range df.duplicates {
		if _, err := os.Stat(v); err == nil {
			t.Errorf("%s shouldn't exist, but it does", v)
		}
	}
}

func TestCopyOriginals(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var df dirFiles
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	df.dir = path + "/testDir"
	df.files = []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
	}
	df.duplicates = []string{
		path + `/testDir/testDir1/file1`,
		path + `/testDir/testDir1/testDir3/file1`,
		path + `/testDir/testDir2/testDir4/file1`,
		path + `/testDir/testDir2/file1`,
		path + `/testDir/testDir2/testDir4/file2`,
		path + `/testDir/testDir1/testDir3/file2`,
		path + `/testDir/testDir2/file2`,
		path + `/testDir/testDir1/file2`,
	}
	*copyFlag = true

	err = df.copyOriginals(ctx)
	if err != nil {
		t.Error("Couldn't copy originals.")
	}

	if _, err := os.Stat(df.dir + "/originals"); os.IsNotExist(err) {
		t.Error(`"originals" directory wasn't created`)
	}

	file1expected := path + "/testDir/originals/file1"
	file2expected := path + "/testDir/originals/file2"
	if _, err := os.Stat(file1expected); err != nil {
		t.Errorf("%s wasn't created", file1expected)
	}
	if _, err := os.Stat(file2expected); err != nil {
		t.Errorf("%s wasn't created", file2expected)
	}

	for _, v := range df.files {
		if _, err := os.Stat(v); err == nil {
			t.Errorf("%s shouldn't exist, but it does", v)
		}
	}
}
