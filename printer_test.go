package mdqi

import (
	"bytes"
	"testing"
)

func TestPrinter(t *testing.T) {
	results := []Result{
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

	var out bytes.Buffer
	if err := Fprint(&out, results); err != nil {
		t.Fatalf("error on print: %s", err.Error())
	}

	expect := `
+-----+----+--------+
| DB  | ID |  NAME  |
+-----+----+--------+
| db1 | 1  | foo    |
| db1 | 2  | barbaz |
| db2 | 3  | qux    |
+-----+----+--------+
`

	if s := out.String(); !compareAfterTrim(s, expect, " \n") {
		t.Fatalf("unexpected output:\n%s", s)
	}
}

func TestPrinterWithEmptyResults(t *testing.T) {
	// avoid panic
	if err := Print([]Result{}); err != nil {
		t.Error("unexpected error:", err)
	}
}
