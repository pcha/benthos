package manager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/Jeffail/benthos/v3/lib/stream"
	yaml "gopkg.in/yaml.v3"
)

func TestFromDirectory(t *testing.T) {
	testDir := t.TempDir()

	barDir := filepath.Join(testDir, "bar")
	if err := os.Mkdir(barDir, 0o777); err != nil {
		t.Fatal(err)
	}

	fooPath := filepath.Join(testDir, "foo.json")
	barPath := filepath.Join(barDir, "test.yaml")

	fooConf := stream.NewConfig()
	fooConf.Input.Type = "bloblang"

	barConf := stream.NewConfig()
	barConf.Input.Type = "nanomsg"

	expConfs := map[string]stream.Config{
		"foo":      fooConf,
		"bar_test": barConf,
	}

	var fooBytes []byte
	var err error
	if fooBytes, err = json.Marshal(fooConf); err != nil {
		t.Fatal(err)
	}
	var barBytes []byte
	if barBytes, err = yaml.Marshal(barConf); err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile(fooPath, fooBytes, 0o666); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(barPath, barBytes, 0o666); err != nil {
		t.Fatal(err)
	}

	var actConfs map[string]stream.Config
	if actConfs, err = LoadStreamConfigsFromDirectory(true, testDir); err != nil {
		t.Fatal(err)
	}

	var actKeys, expKeys []string
	for id := range actConfs {
		actKeys = append(actKeys, id)
	}
	sort.Strings(actKeys)
	for id := range expConfs {
		expKeys = append(expKeys, id)
	}
	sort.Strings(expKeys)

	if !reflect.DeepEqual(actKeys, expKeys) {
		t.Errorf("Wrong keys in loaded set: %v != %v", actKeys, expKeys)
	}

	if exp, act := "bloblang", actConfs["foo"].Input.Type; exp != act {
		t.Errorf("Wrong value in loaded set: %v != %v", act, exp)
	}
	if exp, act := "nanomsg", actConfs["bar_test"].Input.Type; exp != act {
		t.Errorf("Wrong value in loaded set: %v != %v", act, exp)
	}
}
