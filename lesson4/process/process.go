package process

import (
	"GolangBP/lesson4/myErrors"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	DelFlag  = flag.Bool("del", false, "use to delete duplicate files.")
	CopyFlag = flag.Bool("copy", false, "use to move original files to */originals directory.")
	LogFlag  = flag.Bool("log", false, "use to redirect all logs to file */program.log in JSON")
)

type DirFiles struct {
	dir        string
	files      []string
	duplicates []string
	mu         sync.Mutex
	wg         sync.WaitGroup
}

type ProcessDir interface {
	ScanDir() error
	WalkDir() error
}

type ProcessFiles interface {
	FindDuplicates() error
	DeleteDuplicates(context.Context) error
	CopyOriginals(context.Context) error
}

type ProcessAll interface {
	ProcessDir
	ProcessFiles
}

func (df *DirFiles) ScanDir() error {
	fmt.Println("Please, Enter directory to scan for duplicate files.")

	_, err := fmt.Scan(&df.dir)
	if err != nil {
		return err
	}
	err = os.Chdir(df.dir)
	if err != nil {
		return err
	}

	log.Info(`Your directory is`, df.dir)
	return nil
}

func (df *DirFiles) WalkDir() error {
	err := filepath.Walk(df.dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			df.files = append(df.files, path)
		} else {
			log.Info("Scanning directory:", info.Name())
		}
		return nil
	})
	return err
}

func (df *DirFiles) FindDuplicates() error {
	for i := 0; i < len(df.files); i++ {
		iInfo, err := os.Stat(df.files[i])
		if err != nil {
			log.Errorf("Couldn't get %s stats", df.files[i])
			return err
		}
		for j := 1; j < len(df.files); j++ {
			jInfo, err := os.Stat(df.files[j])
			if err != nil {
				log.Errorf("Couldn't get %s stats", df.files[j])
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
		log.Info("No duplicate files found. Exiting..")
		os.Exit(1)
	} else {
		for _, f := range df.files {
			log.Infof(`File "%s" has duplicates`, f)
		}
		log.Infof("Duplicate files: %v", df.duplicates)
	}
	return nil
}

func (df *DirFiles) DeleteDuplicates(ctx context.Context) error {
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
				err := os.Remove(ff)
				if err != nil {
					log.Error(myErrors.CheckError("Couldn't delete file:"), ff)
				}
				df.wg.Done()
			}(f)
		}
		df.wg.Wait()

		log.Info("All duplicate files successfully deleted.")
	case "n":
		log.Info("All duplicate files remain.")
	default:
		err := myErrors.CheckError("Wrong input by user.")
		log.WithError(err).Warn(`Please, answer "y" or "n".`)
		return err
	}
	return nil
}

func (df *DirFiles) CopyOriginals(ctx context.Context) error {
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
		err := os.Mkdir(df.dir+"/originals", os.ModePerm)
		if err != nil {
			err = myErrors.CheckError(`Couldn't create directory "originals"`)
			log.Error(err)
			return err
		}
		for _, f := range df.files {
			df.wg.Add(1)
			go func(ff string) {
				df.mu.Lock()
				defer df.mu.Unlock()
				fInfo, err := os.Stat(ff)
				if err != nil {
					log.Error(myErrors.CheckError("Couldn't get stats of file:"), ff)
				}
				err = os.Rename(ff, df.dir+"/originals/"+fInfo.Name())
				if err != nil {
					log.Error(myErrors.CheckError("Couldn't move original file:"), ff)
				}
				df.wg.Done()
			}(f)
		}
		df.wg.Wait()

		log.Info("Files successfully moved to the new location.")
	case "n":
		log.Info("Program exited at user request.")
	default:
		err := myErrors.CheckError("Wrong input by user.")
		log.WithError(err).Warn(`Please, answer "y" or "n".`)
		return err
	}
	return nil
}
