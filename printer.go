package mdqi

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Print(results []Result) error {
	Fprint(os.Stdout, results)
	return nil
}

func Fprint(out io.Writer, results []Result) error {
	printer := tablewriter.NewWriter(out)
	printer.SetAlignment(tablewriter.ALIGN_LEFT)

	// set header
	headers := []string{"DB"}
	for _, name := range results[0].Columns {
		headers = append(headers, name)
	}
	printer.SetHeader(headers)

	// set rows
	for _, result := range results {
		for _, row := range result.Rows {
			line := []string{result.Database}

			for _, name := range result.Columns {
				line = append(line, fmt.Sprint(row[name]))
			}

			printer.Append(line)
		}
	}

	// print
	printer.Render()

	return nil
}
