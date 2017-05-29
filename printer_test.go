package mdqi

import (
	"bytes"
	"testing"
)

var printTestResults = []Result{
	Result{
		Database: "db1",
		Columns: []string{
			"id",
			"name",
		},
		Rows: []map[string]interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "foo",
			},
			map[string]interface{}{
				"id":   2,
				"name": "barbaz",
			},
		},
	},
	Result{
		Database: "db2",
		Columns: []string{
			"id",
			"name",
		},
		Rows: []map[string]interface{}{
			map[string]interface{}{
				"id":   3,
				"name": "qux",
			},
		},
	},
}

func TestHorizontalPrinter(t *testing.T) {
	var out bytes.Buffer
	if err := Fprint(&out, HorizontalPrinter{}, printTestResults); err != nil {
		t.Fatalf("error on print: %s", err.Error())
	}

	expect := `
+-----+----+--------+
| db  | id |  name  |
+-----+----+--------+
| db1 | 1  | foo    |
| db1 | 2  | barbaz |
| db2 | 3  | qux    |
+-----+----+--------+
`

	if s := out.String(); !compareAfterTrim(s, expect) {
		t.Fatalf("unexpected output:\n%s", s)
	}
}

func TestVerticalPrinter(t *testing.T) {
	var out bytes.Buffer
	if err := Fprint(&out, VerticalPrinter{}, printTestResults); err != nil {
		t.Fatalf("error on print: %s", err.Error())
	}

	expect := `
*************************** 1. row ***************************
  db   : db1
  id   : 1
  name : foo
*************************** 2. row ***************************
  db   : db1
  id   : 2
  name : barbaz
*************************** 3. row ***************************
  db   : db2
  id   : 3
  name : qux
`

	if s := out.String(); !compareAfterTrim(s, expect) {
		t.Fatalf("unexpected output:\n%s", s)
	}
}

func TestPrinterWithEmptyResults(t *testing.T) {
	// avoid panic

	if err := Print(HorizontalPrinter{}, []Result{}); err != nil {
		t.Error("unexpected error:", err)
	}

	if err := Print(VerticalPrinter{}, []Result{}); err != nil {
		t.Error("unexpected error:", err)
	}
}
