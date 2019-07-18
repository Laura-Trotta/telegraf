package remote

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func MockHttpGet(url string) (*http.Response, error) {

	var err error
	array := []byte{}
	bod := InputFromURL{"testfile.conf", array}
	body := []InputFromURL{bod}

	jsn, err := json.Marshal(body)

	resp := http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(jsn)),
	}

	return &resp, err
}

func TestGetExternalFiles(t *testing.T) {

	currdir, err := os.Getwd()
	newdir := currdir + "/dev/"
	os.Mkdir(newdir, 0777)
	dir, err := os.OpenFile(newdir, os.O_CREATE|os.O_WRONLY, 0644)
	defer os.RemoveAll(newdir)

	//Mocked http call, must pass
	err = GetExternalFiles("", newdir, Config, MockHttpGet)

	if err != nil {
		t.Errorf("Function returned error = %s", err.Error())
	}

	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		t.Errorf("File was not created")
	}

}
