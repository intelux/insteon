package plm

import (
	"fmt"
	"io"
)

// Monitor represents a type that can listen on device changes.
type Monitor interface {
	Initialize(*PowerLineModem) error
	Finalize(*PowerLineModem) error
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

func (printMonitor) Initialize(*PowerLineModem) error { return nil }
func (m printMonitor) Finalize(*PowerLineModem) error {
	if closer, ok := m.Writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
func (m printMonitor) ResponseReceived(plm *PowerLineModem, res Response) {
	fmt.Fprintf(m.Writer, "%s\n", res)
}
