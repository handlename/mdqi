package mdqi

import (
	"reflect"
	"sort"
	"testing"
)

func TestManageTags(t *testing.T) {
	sortEqual := func(a, b []string) bool {
		sort.Strings(a)
		sort.Strings(b)

		return reflect.DeepEqual(a, b)
	}

	app, _ := NewApp(Conf{})

	if tags := app.GetTags(); len(tags) != 0 {
		t.Errorf("unexpected tags: %+v", tags)
	}

	app.AddTag("foo")
	app.AddTag("bar")

	if tags := app.GetTags(); !sortEqual(tags, []string{"foo", "bar"}) {
		t.Errorf("unexpected tags: %+v", tags)
	}

	app.RemoveTag("foo")

	if tags := app.GetTags(); !sortEqual(tags, []string{"bar"}) {
		t.Errorf("unexpected tags: %+v", tags)
	}
}
