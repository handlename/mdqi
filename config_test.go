package mdqi

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigFromFile(t *testing.T) {
	conf, err := ConfFromFile("./config.example.yaml")

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expected := Conf{
		Mdq: ConfMdq{
			Bin:    "~/bin/mdq",
			Config: "~/.config/mdq/config.yaml",
		},
		Mdqi: ConfMdqi{
			History:        "~/.mdqi_history",
			DefaultTag:     "db1",
			DefaultDisplay: "horizontal",
		},
	}

	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("unexpected conf: %+v", conf)
		t.Logf("  expected conf: %+v", expected)
	}
}

func TestConfigSetMdqBin(t *testing.T) {
	// invalid path
	{
		conf := Conf{Mdq: ConfMdq{Bin: "/path/to/mdq"}}

		_, err := NewApp(conf)
		if err == nil {
			t.Error("should be error")
		}
	}

	// valid path
	{
		wd, _ := os.Getwd()
		conf := Conf{Mdq: ConfMdq{Bin: filepath.Join(wd, "test", "bin", "mdq")}}

		app, err := NewApp(conf)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		if p := app.mdqPath; p != conf.Mdq.Bin {
			t.Error("failed to set mdq path:", p)
		}
	}
}

func TestConfigSetMdqConfig(t *testing.T) {
	conf := Conf{Mdq: ConfMdq{Config: "/path/to/conf.yaml"}}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if p := app.mdqConfigPath; p != conf.Mdq.Config {
		t.Error("unexpected mdq config path:", p)
	}
}

func TestConfigSetMdqiHistory(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "mdqi")
	conf := Conf{Mdqi: ConfMdqi{History: tmp.Name()}}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if p := app.historyPath; p != conf.Mdqi.History {
		t.Error("unexpected history path:", p)
	}
}

func TestConfigSetMdqiDefaultTag(t *testing.T) {
	conf := Conf{Mdqi: ConfMdqi{DefaultTag: "db1"}}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if tag := app.tag; tag != "db1" {
		t.Error("unexpected tag:", tag)
	}
}

func TestConfigSetMdqiDefaultDisplay(t *testing.T) {
	conf := Conf{Mdqi: ConfMdqi{DefaultDisplay: "vertical"}}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	switch app.printer.(type) {
	case VerticalPrinter:
	default:
		t.Errorf("unexpected printer %s:", app.printer)
	}
}
