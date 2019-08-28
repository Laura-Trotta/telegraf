package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
)

func TestOutputHierarchy(t *testing.T) {
	wd, err := os.Getwd()

	s, _ := serializers.NewInfluxSerializer()
	f := File{
		Directory:  filepath.Join(wd, "test"),
		serializer: s,
		Tag:        "tag1",
	}

	err = f.Connect()
	assert.NoError(t, err)

	err = f.Write(testutil.MockMetrics())
	assert.NoError(t, err)

	//in the test folder, check correct hierarchy
	folder1 := filepath.Join(f.Directory, "test1")
	_, err = os.Stat(folder1)
	assert.NoError(t, err)

	folder2 := filepath.Join(folder1, "value1")
	_, err = os.Stat(folder2)
	assert.NoError(t, err)

	folder3 := filepath.Join(folder2, "2009")
	_, err = os.Stat(folder3)
	assert.NoError(t, err)

	file := filepath.Join(folder3, "314")
	_, err = os.Stat(file)
	assert.NoError(t, err)

	os.RemoveAll(f.Directory)

}
