package planner

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGather(t *testing.T) {

	wd, err := os.Getwd()
	file := File{
		Directory:     filepath.Join(wd, "test/testfiles"),
		Plandirectory: filepath.Join(wd, "test/testplan"),
		TagsList:      "tag=newtag",
	}
	err = file.checkDirNames()
	assert.NoError(t, err)

	p, _ := parsers.NewParser(&parsers.Config{
		MetricName: "example_timeseries",
		DataFormat: "influx",
	})

	file.SetParser(p)

	acc := testutil.Accumulator{}

	todayDate := time.Now()
	testdate := time.Date(todayDate.Year(), todayDate.Month(), todayDate.Day(), 0, 0, 0, 0, time.UTC)

	err = file.Gather(&acc)

	assert.NoError(t, err)

	//check if plan folder has been created
	_, err = os.Stat(file.Plandirectory)
	assert.NoError(t, err)

	//check if plan.json exists and is correct
	plans := []Plan{}

	planfile, err := ioutil.ReadFile(file.Plandirectory + "plan.json")
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(planfile), &plans)
	assert.NoError(t, err)

	//file 1, the only one done
	assert.Equal(t, "0", plans[0].Filename)
	assert.Equal(t, true, plans[0].Done)

	planneddate := time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 0)
	assert.Equal(t, planneddate, plans[0].Day)

	//file 2
	assert.Equal(t, "2", plans[1].Filename)
	assert.Equal(t, false, plans[1].Done)

	planneddate = time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 2)
	assert.Equal(t, planneddate, plans[1].Day)

	//file 3
	assert.Equal(t, "4", plans[2].Filename)
	assert.Equal(t, false, plans[2].Done)

	planneddate = time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 4)
	assert.Equal(t, planneddate, plans[2].Day)

	//check if accumulator has metrics
	metrics := acc.Metrics

	resultTime := time.Date(todayDate.Year(), todayDate.Month(), todayDate.Day(), 13, 0, 0, 0, time.UTC)

	assert.Equal(t, len(metrics), 1)
	assert.Equal(t, metrics[0].Measurement, "example_timeseries")
	assert.Equal(t, metrics[0].Time.UTC(), resultTime)
	assert.Equal(t, metrics[0].Tags["tag"], "newtag")
	assert.Equal(t, metrics[0].Fields["field"], "value")

	os.RemoveAll(file.Plandirectory)

}
