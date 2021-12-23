//go:build integration
// +build integration

package process

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func SetupTest() error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}

	testFiles := []string{
		path + "/testDir/",
		path + "/testDir/testDir1/",
		path + "/testDir/testDir1/testDir3/",
		path + "/testDir/testDir2/",
		path + "/testDir/testDir2/testDir4/",
	}
	for _, v := range testFiles {
		err := os.MkdirAll(v, os.ModePerm)
		if err != nil {
			log.Println(err)
			return err
		}

		f, err := os.Create(v + "/file2")
		if err != nil {
			log.Println(err)
			return err
		}
		data := []byte("this is test file2")
		_, err = f.Write(data)
		if err != nil {
			log.Println(err)
			return err
		}
		f.Close()

		f, err = os.Create(v + "/file1")
		if err != nil {
			log.Println(err)
			return err
		}
		data = []byte("this is test file1")
		_, err = f.Write(data)
		if err != nil {
			log.Println(err)
			return err
		}
		f.Close()

	}
	log.Println("Test setup complete.")
	return nil
}

func TestScanDir(t *testing.T) {
	var df DirFiles

	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(path + "/testDir"); err != nil {
		if os.IsNotExist(err) {
			if err = SetupTest(); err != nil {
				t.Fatal("Couldn't setup test.")
			}
		}
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

	err = df.ScanDir()
	if err != nil {
		t.Errorf(`The directory "%s" doesn't exist.`, df.dir)
	}
	log.Println("TestScanDir completed.")

	err = os.Chdir(path)
	if err != nil {
		t.Error("Couldn't restore working directory")
	}
}

func TestWalkDir(t *testing.T) {
	var df DirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path + "/testDir"); err != nil {
		if os.IsNotExist(err) {
			if err = SetupTest(); err != nil {
				t.Fatal("Couldn't setup test.")
			}
		}
	}

	df.dir = path + "/testDir"
	err = df.WalkDir()
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
	log.Println("TestWalkDir completed.")
}

func TestFindDuplicates(t *testing.T) {
	var df DirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path + "/testDir"); err != nil {
		if os.IsNotExist(err) {
			if err = SetupTest(); err != nil {
				t.Fatal("Couldn't setup test.")
			}
		}
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

	err = df.FindDuplicates()
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
	log.Println("TestFindDuplicates completed.")
}

func TestDeleteDuplicates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var df DirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path + "/testDir"); err != nil {
		if os.IsNotExist(err) {
			if err = SetupTest(); err != nil {
				t.Fatal("Couldn't setup test.")
			}
		}
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
	*DelFlag = true

	err = df.DeleteDuplicates(ctx)
	if err != nil {
		t.Error("Couldn't delete duplicates.")
	}
	for _, v := range df.duplicates {
		if _, err := os.Stat(v); err == nil {
			t.Errorf("%s shouldn't exist, but it does", v)
		}
	}
	log.Println("TestDeleteDuplicates completed.")
}

func TestCopyOriginals(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var df DirFiles

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path + "/testDir"); err != nil {
		if os.IsNotExist(err) {
			if err = SetupTest(); err != nil {
				t.Fatal("Couldn't setup test.")
			}
		}
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
	*CopyFlag = true

	err = df.CopyOriginals(ctx)
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
	log.Println("TestCopyOriginals completed.")
}
