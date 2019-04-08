package plm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// CommandCoder is the interface for objects that have a command code.
type CommandCoder interface {
	commandCode() CommandCode
}

// Request is the interface for all requests.
type Request interface {
	CommandCoder
	marshal(io.Writer) error
}

// Response is the interface for all responses.
type Response interface {
	CommandCoder
	unmarshal(io.Reader) error
}

var (
	// ErrCommandFailure indicates that a previous command failed.
	ErrCommandFailure = errors.New("command failed")
)

type errUnknownCommand struct {
	CommandCode
}

// Error returns the error string.
func (e errUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command code (%x)", e.CommandCode)
}

// MarshalRequest serializes a request to a writer.
func MarshalRequest(w io.Writer, request Request) error {
	// Make sure we send it all at once. Not that it is mandatory, but it's
	// easier to debug.
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.Write([]byte{
		byte(MessageStart),
		byte(request.commandCode()),
	})

	if err != nil {
		return err
	}

	return request.marshal(bw)
}

// UnmarshalResponses parses responses from a reader and returns the index of the response that was unmarshalled.
func UnmarshalResponses(r io.Reader, responses []Response) (int, error) {
	for {
		mark, err := skipToMessage(r)

		if err != nil {
			return -1, err
		}

		if mark == MessageNak {
			return -1, ErrCommandFailure
		}

		buffer := make([]byte, 1)

		_, err = r.Read(buffer)

		if err != nil {
			return -1, err
		}

		commandCode := CommandCode(buffer[0])

		for i, response := range responses {
			if commandCode == response.commandCode() {
				err = response.unmarshal(r)

				if err != nil {
					return i, err
				}

				return i, nil
			}
		}

		// No response matches the command code. Let's read it all and move-on to the next message.
		// If a command code is missing from the table, we read nothing and
		// wait for the next message start, effectively discarding any bytes
		// in-between.
		buffer = make([]byte, responsesSizes[commandCode])
		_, err = io.ReadFull(r, buffer)

		if err != nil {
			return -1, err
		}
	}
}

// UnmarshalResponse parses a response from a reader.
func UnmarshalResponse(r io.Reader, response Response) error {
	_, err := UnmarshalResponses(r, []Response{response})
	return err
}

// skipToMessage skips bytes on the specified io.Reader until a message start
// is met.
//
// The message start or message nak is returned.
func skipToMessage(r io.Reader) (byte, error) {
	buffer := make([]byte, 1)

	for buffer[0] != MessageStart && buffer[0] != MessageNak {
		_, err := r.Read(buffer)

		if err != nil {
			return 0, err
		}
	}

	return buffer[0], nil
}

var responsesSizes = map[CommandCode]int{
	StandardMessageReceived:       9,
	ExtendedMessageReceived:       23,
	X10Received:                   2,
	AllLinkingCompleted:           8,
	ButtonEventReport:             1,
	UserResetDetected:             0,
	AllLinkCleanupFailureReport:   5,
	AllLinkRecordMessage:          8,
	AllLinkCleanupStatusReport:    1,
	GetIMInfo:                     7,
	GetFirstAllLinkRecord:         1,
	GetNextAllLinkRecord:          1,
	StartAllLinking:               3,
	CancelAllLinking:              1,
	SendStandardOrExtendedMessage: 7,
}
