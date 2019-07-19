package planner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
)

type File struct {
	Directory     string `toml:"directory"`
	Plandirectory string `toml:"plandirectory"`
	parser        parsers.Parser
}

type Plan struct {
	Day      time.Time `json:"day"`
	Filename string    `json:"filename"`
	Done     bool      `json:"done"`
}

var Plugin telegraf.Input

const sampleConfig = `
## Directory containing the files to be read
directory = ""

## Directory where the plan will be saved
plandirectory = ""

## The dataformat to be read from files
## Each data format has its own unique set of configuration options, read
## more about them here:
## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
data_format = "influx"
`

// SampleConfig returns the default configuration of the Input
func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Description() string {
	return "Modifies the timestamp of one file per day in the directory"
}

func checkDirNames(f *File) error {
	if len(f.Plandirectory) < 1 || len(f.Directory) < 1 {
		return fmt.Errorf("Must provide path for both directories")
	}
	lastChar := f.Plandirectory[len(f.Plandirectory)-1:]
	if lastChar != string(os.PathSeparator) {
		f.Plandirectory = f.Plandirectory + string(os.PathSeparator)
	}
	lastChar = f.Directory[len(f.Directory)-1:]
	if lastChar != string(os.PathSeparator) {
		f.Directory = f.Directory + string(os.PathSeparator)
	}

	return nil
}

func initialize(f *File, first bool) error {

	var names []string

	err := filepath.Walk(f.Directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			names = append(names, info.Name())
		}
		return nil
	})

	if err != nil {
		return err
	}

	plans := make([]Plan, len(names))

	for i, name := range names {

		plus := i
		if !first {
			plus++
		}
		plan := Plan{time.Now().UTC().AddDate(0, 0, plus).Truncate(time.Hour), name, false}

		plans[i] = plan

	}

	file, err := json.MarshalIndent(plans, "", "")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(f.Plandirectory+"plan.json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Gather(acc telegraf.Accumulator) error {

	err := checkDirNames(f)
	if err != nil {
		return err
	}

	_, errstat := os.Stat(f.Plandirectory + "plan.json")
	if os.IsNotExist(errstat) {
		err := initialize(f, true)
		if err != nil {
			return err
		}
	}

	plans := []Plan{}

	file, err := ioutil.ReadFile(f.Plandirectory + "plan.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(file), &plans)
	if err != nil {
		return err
	}

	workdone := true

	for i, plan := range plans {

		if plan.Done == false && time.Now().UTC().After(plan.Day) {

			metrics, err := f.readMetric(f.Directory + plan.Filename)
			if err != nil {
				return err
			}
			for _, m := range metrics {

				newtime := time.Date(plan.Day.Year(), plan.Day.Month(), plan.Day.Day(), m.Time().Hour(), m.Time().Minute(), m.Time().Second(), m.Time().Nanosecond(), time.UTC)
				m.SetTime(newtime)
				acc.AddFields(m.Name(), m.Fields(), m.Tags(), m.Time())

			}

			plan.Done = true

			plans[i] = plan

			file, err := json.MarshalIndent(plans, "", "")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(f.Plandirectory+"plan.json", file, 0644)
			if err != nil {
				return err
			}

		}

		if plan.Done == false {
			workdone = false
		}
	}

	if workdone {
		os.Remove(f.Plandirectory + "plan.json")
		err := initialize(f, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) SetParser(p parsers.Parser) {
	f.parser = p
}

func (f *File) readMetric(filename string) ([]telegraf.Metric, error) {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("E! Error file: %v could not be read, %s", filename, err)
	}
	return f.parser.Parse(fileContents)

}

func init() {
	inputs.Add("planner", func() telegraf.Input {
		return &File{}
	})
}
