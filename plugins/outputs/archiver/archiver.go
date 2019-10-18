package archiver

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

type File struct {
	Directory string `toml:"directory"`
	Tag       string `toml:"tag"`

	writer     io.Writer
	closers    []io.Closer
	serializer serializers.Serializer
}

var Plugin telegraf.Output

var sampleConfig = `
  ## Directory to write to
  directory = ""

  ## Tag 
  tag = ""

  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "influx"
`

func (f *File) SetSerializer(serializer serializers.Serializer) {
	f.serializer = serializer
}

func (f *File) Connect() error {
	if len(f.Directory) < 1 {
		return fmt.Errorf("Must provide a file path")
	}
	lastChar := f.Directory[len(f.Directory)-1:]
	if lastChar != string(os.PathSeparator) {
		f.Directory = f.Directory + string(os.PathSeparator)
	}
	return nil
}

func (f *File) Close() error {
	var err error
	return err
}

func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Description() string {
	return "Send telegraf metrics to file depending on content"
}

func (f *File) Write(metrics []telegraf.Metric) error {

	var writeErr error = nil

	for i, metric := range metrics {

		time := metrics[i].Time()
		year := time.Year()
		day := time.YearDay()

		//Getting fields to build filename
		var filename strings.Builder
		filename.WriteString(metrics[i].Name())
		filename.WriteString(string(os.PathSeparator))
		filename.WriteString(metrics[i].Tags()[f.Tag])
		filename.WriteString(string(os.PathSeparator))
		filename.WriteString(strconv.Itoa(year))
		dirname := filename.String()
		filename.WriteString(string(os.PathSeparator))
		filename.WriteString(strconv.Itoa(day))

		b, err := f.serializer.Serialize(metric)
		if err != nil {
			log.Printf("D! [outputs.archiver] Could not serialize metric: %v", err)
		}
		if _, err := os.Stat(f.Directory + dirname); os.IsNotExist(err) {
			os.MkdirAll(f.Directory+dirname, os.ModePerm)
		}
		file, errf := os.OpenFile(f.Directory+filename.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if errf != nil {
			log.Println(errf)
		}
		defer file.Close()

		if _, errf := file.Write(b); errf != nil {
			log.Println(errf)
		}
	}

	return writeErr
}

func init() {
	outputs.Add("archiver", func() telegraf.Output {
		return &File{}
	})
}
