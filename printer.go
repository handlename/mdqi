package mdqi

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

var defaultOutput io.Writer
var defaultPrinter Printer = func() Printer {
	return HorizontalPrinter{}
}()

type Printer interface {
	Print(out io.Writer, results []Result) error
}

type HorizontalPrinter struct{}

func (p HorizontalPrinter) Print(out io.Writer, results []Result) error {
	table := tablewriter.NewWriter(out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)

	// set header
	headers := []string{"db"}
	for _, name := range results[0].Columns {
		headers = append(headers, name)
	}
	table.SetHeader(headers)

	// set rows
	for _, result := range results {
		for _, row := range result.Rows {
			line := []string{result.Database}

			for _, name := range result.Columns {
				line = append(line, fmt.Sprint(row[name]))
			}

			table.Append(line)
		}
	}

	table.Render()

	return nil
}

type VerticalPrinter struct{}

func (p VerticalPrinter) Print(out io.Writer, results []Result) error {
	num := 0

	for _, result := range results {
		for _, row := range result.Rows {
			num += 1

			fmt.Fprintf(out, "*************************** %d. row ***************************\n", num)

			table := tablewriter.NewWriter(out)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoFormatHeaders(false)
			table.SetAutoWrapText(false)
			table.SetColumnSeparator(":")
			table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})

			table.Append([]string{"db", result.Database})

			for _, name := range result.Columns {
				table.Append([]string{name, fmt.Sprint(row[name])})
			}

			table.Render()
		}
	}

	return nil
}

func Print(printer Printer, results []Result) error {
	Fprint(defaultOutput, printer, results)
	return nil
}

func Fprint(out io.Writer, printer Printer, results []Result) error {
	if len(results) == 0 {
		return nil
	}

	printer.Print(out, results)

	return nil
}
