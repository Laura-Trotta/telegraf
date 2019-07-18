package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Comment this line to disable pprof endpoint.
	"os"
	"strconv"
	"strings"
)

//Generic input file
type InputFromURL struct {
	//file name
	Name string `json:"name"`
	//file content
	Content []byte `json:"content"`
}

//Type of input file, can be plugin or configuration file
type FileType int

const (
	Plugin FileType = 0
	Config FileType = 1
)

type GeneralHttpGet func(url string) (*http.Response, error)

func HttpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	return resp, err
}

//function to get files
//Gets an array of type InputFromUrl from the given url,
// creates corresponding files in the given directory
func GetExternalFiles(url, dir string, filetype FileType, fget GeneralHttpGet) error {

	resp, err := fget(url)
	if err != nil {
		return err
	}
	if strings.HasPrefix(strconv.Itoa(resp.StatusCode), "5") {
		fmt.Printf("E! Connection error to external resources URL %s\n", url)
		return err
	} else if strings.HasPrefix(strconv.Itoa(resp.StatusCode), "4") {
		fmt.Printf("W! File from external resources URL %s not found\n", url)
	}

	defer resp.Body.Close()

	var inputs []InputFromURL
	json.NewDecoder(resp.Body).Decode(&inputs)

	lastChar := dir[len(dir)-1:]
	if lastChar != string(os.PathSeparator) {
		dir = dir + string(os.PathSeparator)
	}

	for _, input := range inputs {

		hasPluginPrefix := strings.HasSuffix(input.Name, ".so") || strings.HasSuffix(input.Name, ".dll")
		hasConfigPrefix := strings.HasSuffix(input.Name, ".conf")
		if !(filetype == Plugin && hasPluginPrefix || filetype == Config && hasConfigPrefix) {
			fmt.Printf("W! File %v has an uncorrect name\n", input.Name)
		}

		file, err := os.OpenFile(dir+input.Name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf(err.Error())
		}
		defer file.Close()

		if _, err := file.Write(input.Content); err != nil {
			fmt.Printf(err.Error())
		}
	}
	return nil
}
