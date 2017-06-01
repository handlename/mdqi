package mdqi

import (
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
