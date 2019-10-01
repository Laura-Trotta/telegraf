package planner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
)

type File struct {
	Directory     string `toml:"directory"`
	Plandirectory string `toml:"plandirectory"`
	TagsList      string `toml:"tagslist"`
	parser        parsers.Parser
}

type Plan struct {
	Day      time.Time `json:"day"`
	OldDay   time.Time `json:"old_day"`
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

func getTags(f *File) map[string]string {

	tagsmap := make(map[string]string)

	couples := strings.Split(f.TagsList, ",")

	for _, couple := range couples {

		single := strings.Split(couple, "=")

		tagsmap[single[0]] = single[1]

	}

	return tagsmap
}

//Checks if the directory provided in configuration are existent and correct
func (f *File) checkDirNames() error {
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

	_, errdir := os.Stat(f.Plandirectory)
	if os.IsNotExist(errdir) {
		err := os.MkdirAll(f.Plandirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

//Changes the day and choosen tag of the metrics
//Returns the original date of the metrics
func (f *File) modifyMetrics(plan Plan, tagsmap map[string]string, acc telegraf.Accumulator) error {
	metrics, err := f.readMetric(filepath.Join(f.Directory, plan.Filename))
	if err != nil {
		return err
	}

	for _, m := range metrics {

		newtime := time.Date(plan.Day.Year(), plan.Day.Month(), plan.Day.Day(), m.Time().Hour(), m.Time().Minute(), m.Time().Second(), m.Time().Nanosecond(), m.Time().Location())

		m.SetTime(newtime)

		for key, value := range tagsmap {
			m.AddTag(key, value)
		}

		acc.AddFields(m.Name(), m.Fields(), m.Tags(), m.Time())

	}

	return nil
}

//Saves changes to the json file plan.json
func (f *File) savePlans(plans []Plan) error {

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

//Creates the initial plan
func (f *File) initialize(reference time.Time) error {

	var names []int

	err := filepath.Walk(f.Directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			intname, _ := strconv.Atoi(info.Name())
			names = append(names, intname)
		}
		return nil
	})

	if err != nil {
		return err
	}

	plans := make([]Plan, len(names))

	sort.Ints(names)

	date := time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

	for i, name := range names {

		stringName := strconv.Itoa(name)

		metrics, err := f.readMetric(filepath.Join(f.Directory, stringName))
		if err != nil {
			return err
		}

		oldDay := metrics[0].Time().UTC()

		plan := Plan{date.AddDate(0, 0, name), oldDay, stringName, false}

		plans[i] = plan

	}

	f.savePlans(plans)

	return nil
}

func (f *File) Gather(acc telegraf.Accumulator) error {

	err := f.checkDirNames()
	if err != nil {
		return err
	}

	_, errstat := os.Stat(f.Plandirectory + "plan.json")
	if os.IsNotExist(errstat) {
		err := f.initialize(time.Now().AddDate(0, 0, -1).UTC())
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

	tagsmap := getTags(f)

	workdone := true

	for i, plan := range plans {

		if plan.Done == false && time.Now().UTC().After(plan.Day) {

			err := f.modifyMetrics(plan, tagsmap, acc)

			if err == nil {
				plan.Done = true
				plans[i] = plan

				f.savePlans(plans)
			}

		}

		if plan.Done == false {
			workdone = false
		}
	}

	if workdone {
		lastday := plans[len(plans)-1].Day
		os.Remove(f.Plandirectory + "plan.json")
		err := f.initialize(lastday)
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
