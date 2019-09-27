package planner

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {

	wd, err := os.Getwd()
	file := File{
		Directory:     filepath.Join(wd, "test/testfiles"),
		Plandirectory: filepath.Join(wd, "test/testplan"),
	}
	err = file.checkDirNames()
	assert.NoError(t, err)

	testdate := time.Now()

	err = file.initialize(testdate.AddDate(0, 0, -1))
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

	//file 1
	assert.Equal(t, "0", plans[0].Filename)
	assert.Equal(t, false, plans[0].DoneT)

	planneddate := time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 0)
	assert.Equal(t, planneddate, plans[0].Day)

	//file 2
	assert.Equal(t, "2", plans[1].Filename)
	assert.Equal(t, false, plans[1].DoneT)

	planneddate = time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 2)
	assert.Equal(t, planneddate, plans[1].Day)

	//file 3
	assert.Equal(t, "4", plans[2].Filename)
	assert.Equal(t, false, plans[2].DoneT)

	planneddate = time.Date(testdate.Year(), testdate.Month(), testdate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 4)
	assert.Equal(t, planneddate, plans[2].Day)

	os.RemoveAll(file.Plandirectory)

}
