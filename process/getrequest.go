package process

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var GitCommit string

func GetRequest(defaultFile string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "Commit version is: %s\n", GitCommit)
	fmt.Fprintln(os.Stdout, "Please, Enter an SQL request: SELECT *column_name* FROM *file_name* WHERE *search_parameter* AND *search_parameter*.\nSearch parameters are optional.\nTo use default file from config type `default` instead of file name.\nPress `Enter` to confirm.")
	fmt.Fprintf(os.Stdout, "Default file is: %s\n", defaultFile)

	request, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if !strings.Contains(request, sel) && !strings.Contains(request, from) {
		err = fmt.Errorf("wrong syntax of user request")
		return "", err
	}

	f, err := os.OpenFile("logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + request)
	if err != nil {
		return "", err
	}

	err = w.Flush()
	if err != nil {
		return "", err
	}

	return request, nil
}
