package config_test

import (
	"github.com/ameamezhou/distributekv/config"
	"github.com/ameamezhou/distributekv/xlog"
	"os"
	"reflect"
	"testing"
)

func TestParseShards(t *testing.T) {
	contents := `
[[shard]]
name = "xiaoqizhou"
idx = 0
address = "localhost:8080"

[[shard]]
name = "xiawuyue"
idx = 1
address = "localhost:8081"

[[shard]]
name = "mengzhi"
idx = 2
address = "localhost:8082"
`
	f, err := os.CreateTemp(os.TempDir(), "config.toml")
	if err != nil {
		xlog.Fatal("could not create a temp file ", err)
	}
	defer f.Close()
	name := f.Name()
	defer os.Remove(name)

	_, err = f.WriteString(contents)
	if err != nil {
		t.Fatal("counld not write the config contents", err)
	}

	c, err := config.ParseFile(name)
	if err != nil {
		t.Fatalf("could not parse config: %v", err)
	}
	want := config.Config{
		Shard: []config.Shard {
			{
				Name: "xiaoqizhou",
				Index: 0,
				Address: "localhost:8080",
			},
			{
				Name: "xiawuyue",
				Index: 1,
				Address: "localhost:8081",
			},
			{
				Name: "mengzhi",
				Index: 2,
				Address: "localhost:8082",
			},
		},
	}

	if !reflect.DeepEqual(c, want) {
		t.Errorf("the config dose match : got: %#v, want: %#v", c, want)
	}

}


