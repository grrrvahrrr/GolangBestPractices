//go:build !integration
// +build !integration

package process

import (
	"GolangBP/lesson4/myErrors"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var fs fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	Remove(name string) error
	Rename(oldPath string, newPath string) error
	Mkdir(name string) error
	Chdir(dir string) error
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type osFS struct{}

func (osFS) Open(name string) (file, error)              { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error)       { return os.Stat(name) }
func (osFS) Remove(name string) error                    { return nil }
func (osFS) Rename(oldPath string, newPath string) error { return nil }
func (osFS) Mkdir(name string) error                     { return nil }
func (osFS) Chdir(dir string) error                      { return nil }

type mockedFS struct {
	osFS

	reportErr   bool
	reportSize1 int64
	reportSize2 int64
	reportName  string
}

type mockedFileInfo struct {
	os.FileInfo
	size int64
	name string
}

func (m mockedFileInfo) Size() int64  { return m.size }
func (m mockedFileInfo) Name() string { return m.name }

func (m mockedFS) Stat(name string) (os.FileInfo, error) {
	if m.reportErr {
		return nil, os.ErrNotExist
	}
	if strings.Contains(name, "file1") {
		return mockedFileInfo{size: m.reportSize1, name: m.reportName}, nil
	} else {
		return mockedFileInfo{size: m.reportSize2, name: m.reportName}, nil
	}

}
func (m mockedFS) Chdir(dir string) error {
	if m.reportErr {
		return os.ErrNotExist
	}
	return nil
}

func (m mockedFS) Remove(name string) error {
	if m.reportErr {
		return os.ErrNotExist
	}
	return nil
}

func (m mockedFS) Mkdir(name string) error {
	if m.reportErr {
		return os.ErrNotExist
	}
	return nil
}

func (df *DirFiles) MockedScanDir() error {
	fmt.Println("Please, Enter directory to scan for duplicate files.")

	_, err := fmt.Scan(&df.dir)
	if err != nil {
		return err
	}
	err = fs.Chdir(df.dir)
	if err != nil {
		return err
	}

	log.Println(`Your directory is`, df.dir)
	return nil
}

func TestMockedScanDir(t *testing.T) {
	oldFs := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	var df DirFiles

	mfs.reportName = "/testDir"
	mfs.reportErr = false

	input := []byte(mfs.reportName)
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

	//No errors found
	if !assert.Nil(t, df.MockedScanDir()) {
		t.Errorf(`The directory "%s" doesn't exist.`, df.dir)
	}

	//There is an error
	mfs.reportErr = true
	if !assert.NotNil(t, df.MockedScanDir()) {
		t.Errorf(`The directory "%s" exists.`, df.dir)
	}
}

func (df *DirFiles) MockedFindDuplicates() error {
	for i := 0; i < len(df.files); i++ {
		iInfo, err := fs.Stat(df.files[i])
		if err != nil {
			log.Printf("Couldn't get %s stats", df.files[i])
			return err
		}
		for j := 1; j < len(df.files); j++ {
			jInfo, err := fs.Stat(df.files[j])
			if err != nil {
				log.Printf("Couldn't get %s stats", df.files[j])
				return err
			}
			if iInfo.Name() == jInfo.Name() && iInfo.Size() == jInfo.Size() && df.files[i] != df.files[j] {
				df.duplicates = append(df.duplicates, df.files[j])
				df.files[j] = df.files[len(df.files)-1]
				df.files[len(df.files)-1] = ""
				df.files = df.files[:len(df.files)-1]
				j = 1
			}
		}
	}
	if len(df.duplicates) == 0 {
		log.Println("No duplicate files found. Exiting..")
		os.Exit(1)
	} else {
		for _, f := range df.files {
			log.Printf(`File "%s" has duplicates`, f)
		}
		log.Printf("Duplicate files: %v", df.duplicates)
	}
	return nil
}

func TestMockedFindDuplicates(t *testing.T) {
	oldFs := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	var df DirFiles

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

	mfs.reportErr = false
	mfs.reportSize1 = 1
	mfs.reportSize2 = 2

	err = df.MockedFindDuplicates()
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

	if !assert.Equal(t, duplicatesExpected, df.duplicates) {
		t.Errorf("In files slice expected %s, but got %s", duplicatesExpected, df.duplicates)
	}

	filesExpected := []string{
		path + "/testDir/file1",
		path + "/testDir/file2",
	}

	//No errors found
	if !assert.Equal(t, filesExpected, df.files) {
		t.Errorf("In files slice expected %s, but got %s", filesExpected, df.files)
	}

	//There is an error getting stats
	mfs.reportErr = true
	if !assert.NotNil(t, df.MockedFindDuplicates()) {
		t.Errorf(`This should return an error`)
	}

}

func (df *DirFiles) MockedDeleteDuplicates(ctx context.Context) error {
	var answer string

	if *DelFlag {
		answer = "y"
	} else {
		fmt.Println("Do you wish to delete all duplicates? y/n")
		_, err := fmt.Scan(&answer)
		for err != nil {
			return err
		}
	}
	switch answer {
	case "y":
		for _, f := range df.duplicates {
			df.wg.Add(1)
			go func(ff string) {
				df.mu.Lock()
				defer df.mu.Unlock()
				err := fs.Remove(ff)
				if err != nil {
					log.Println(myErrors.CheckError("Couldn't delete file:"), ff)
				}
				df.wg.Done()
			}(f)
		}
		df.wg.Wait()

		log.Println("All duplicate files successfully deleted.")
	case "n":
		log.Println("All duplicate files remain.")
	default:
		err := myErrors.CheckError("Wrong input by user.")
		log.Println(`Please, answer "y" or "n".`)
		return err
	}
	return nil
}

func TestMockedDeleteDuplicates(t *testing.T) {
	oldFs := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var df DirFiles

	df.duplicates = []string{
		`/testDir/testDir1/file1`,
		`/testDir/testDir1/testDir3/file1`,
		`/testDir/testDir2/testDir4/file1`,
		`/testDir/testDir2/file1`,
		`/testDir/testDir2/testDir4/file2`,
		`/testDir/testDir1/testDir3/file2`,
		`/testDir/testDir2/file2`,
		`/testDir/testDir1/file2`,
	}
	*DelFlag = true

	mfs.reportErr = false

	//No errors found
	if !assert.Nil(t, df.MockedDeleteDuplicates(ctx)) {
		t.Error("Couldn't delete duplicates.")
	}

	//There is a goroutine error
	mfs.reportErr = true
	if assert.Nil(t, df.MockedDeleteDuplicates(ctx)) {
		log.Println("Test: Should recieve errors messages from goroutines.")
	}
}

func (df *DirFiles) MockedCopyOriginals(ctx context.Context) error {
	var answer string

	if *CopyFlag {
		answer = "y"
	} else {
		fmt.Println(`Do you wish to move all original files to a new directory "*your_dir*/originals"? y/n`)
		_, err := fmt.Scan(&answer)
		for err != nil {
			return err
		}
	}
	switch answer {
	case "y":

		err := fs.Mkdir(df.dir + "/originals")
		if err != nil {
			err = myErrors.CheckError(`Couldn't create directory "originals"`)
			log.Println(err)
			return err
		}
		for _, f := range df.files {
			df.wg.Add(1)
			go func(ff string) {
				df.mu.Lock()
				defer df.mu.Unlock()
				fInfo, err := fs.Stat(ff)
				if err != nil {
					log.Println(myErrors.CheckError("Couldn't get stats of file:"), ff)
				}
				err = fs.Rename(ff, df.dir+"/originals/"+fInfo.Name())
				if err != nil {
					log.Println(myErrors.CheckError("Couldn't move original file:"), ff)
				}
				df.wg.Done()
			}(f)
		}
		df.wg.Wait()

		log.Println("Files successfully moved to the new location.")
	case "n":
		log.Println("Program exited at user request.")
	default:
		err := myErrors.CheckError("Wrong input by user.")
		log.Println(`Please, answer "y" or "n".`)
		return err
	}
	return nil
}

func TestMockedCopyOriginals(t *testing.T) {
	oldFs := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var df DirFiles

	df.dir = "/testDir"
	df.files = []string{
		"/testDir/file1",
		"/testDir/file2",
	}
	*CopyFlag = true

	mfs.reportErr = false

	//No errors found
	if !assert.Nil(t, df.MockedCopyOriginals(ctx)) {
		t.Error("Couldn't copy originals.")
	}

	//There is an error
	mfs.reportErr = true
	if !assert.NotNil(t, df.MockedCopyOriginals(ctx)) {
		t.Errorf(`There should be an error.`)
	}
}
