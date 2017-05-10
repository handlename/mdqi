package mdqi

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrinter(t *testing.T) {
	app, _ := NewApp(Conf{})

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
	if err := app.Fprint(&out, results); err != nil {
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

func compareAfterTrim(a, b string, cutset string) bool {
	return strings.Trim(a, cutset) == strings.Trim(b, cutset)
}
