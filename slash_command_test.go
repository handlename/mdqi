package mdqi

import (
	"reflect"
	"testing"
)

func TestParseSlashCommand(t *testing.T) {
	type Expected struct {
		Result *SlashCommand
		Error  error
	}

	type Pattern struct {
		Title    string
		Query    string
		Expected Expected
	}

	patterns := []Pattern{
		Pattern{
			Title: "command only category",
			Query: "/quit",
			Expected: Expected{
				Result: &SlashCommand{
					Category: "quit",
					Name:     "",
					Args:     []string{},
				},
				Error: nil,
			},
		},
		Pattern{
			Title: "command with options",
			Query: "/tag set db1",
			Expected: Expected{
				Result: &SlashCommand{
					Category: "tag",
					Name:     "set",
					Args:     []string{"db1"},
				},
				Error: nil,
			},
		},
		Pattern{
			Title: "not a command",
			Query: "select * from item",
			Expected: Expected{
				Result: nil,
				Error:  ErrNotASlashCommand,
			},
		},
	}

	for _, pattern := range patterns {
		t.Log("test:", pattern.Title)

		r, err := ParseSlashCommand(pattern.Query)

		if !reflect.DeepEqual(r, pattern.Expected.Result) {
			t.Errorf("unexpected result: %#v", r)
		}

		if err != pattern.Expected.Error {
			t.Error("unexpected error:", err)
		}
	}
}
