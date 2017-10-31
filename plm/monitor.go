package plm

import (
	"fmt"
	"io"
)

// Monitor represents a type that can listen on device changes.
type Monitor interface {
	ResponseReceived(*PowerLineModem, Response)
}

type printMonitor struct {
	Writer io.Writer
}

// NewPrintMonitor instanciates a new monitor that prints out what it reads.
func NewPrintMonitor(w io.Writer) Monitor {
	return printMonitor{
		Writer: w,
	}
}

func (m printMonitor) ResponseReceived(plm *PowerLineModem, res Response) {
	fmt.Fprintf(m.Writer, "%s\n", res)
}
