package cli

import (
	"io"
)

// ResultPrinter is an interface for result printer.
type ResultPrinter interface {
	// Print is a method to print the data in the specified format
	Print()
}

// ResultRenderer is an interface for result renderer that renders the result to a string.
// the reason for dividing ResultPrinter and ResultRenderer is because cli-common json printer is not returning string
type ResultRenderer interface {
	// Render is a method to render the data to a string in the specified format
	Render() string
}

// PrinterFactory is an interface for printer factory.
type PrinterFactory interface {
	// CreatePrinter is a method to create a new ResultPrinter that prints the data in the specified format
	CreatePrinter(writer io.Writer, data interface{}) (ResultPrinter, error)
}

// // ResultPrinterFactory is a struct to implement the PrinterFactory interface
// type ResultPrinterFactory struct {
// 	// currently we have json and table(text) printer format
// 	tablePrinterFactory PrinterFactory
// }

// // CreatePrinter is a method to create a new ResultPrinter based on the output format
// func (p *ResultPrinterFactory) CreatePrinter(format string, writer io.Writer, data interface{}) (ResultPrinter, error) {
// 	switch format {
// 	case "json":
// 		return NewJSONPrinter(writer, data), nil
// 	case "text":
// 		return p.tablePrinterFactory.CreatePrinter(writer, data)
// 	default:
// 		return nil, fmt.Errorf("unsupported output format: %s", format)
// 	}
// }
