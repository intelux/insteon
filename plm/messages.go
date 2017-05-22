package plm

import (
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
	write(io.Writer) error
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
	_, err := w.Write([]byte{
		byte(MessageStart),
		byte(request.commandCode()),
	})

	if err != nil {
		return err
	}

	return request.write(w)
}

// UnmarshalResponse parses a response from a reader.
func UnmarshalResponse(r io.Reader, response Response) error {
	for {
		mark, err := skipToMessage(r)

		if err != nil {
			return err
		}

		if mark == MessageNak {
			return ErrCommandFailure
		}

		buffer := make([]byte, 1)

		_, err = r.Read(buffer)

		if err != nil {
			return err
		}

		commandCode := CommandCode(buffer[0])

		if commandCode != response.commandCode() {
			buffer := make([]byte, responsesSizes[commandCode])
			_, err = io.ReadFull(r, buffer)

			if err != nil {
				return err
			}
		} else {
			err = response.unmarshal(r)

			if err != nil {
				return err
			}

			return nil
		}
	}
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
	GetIMInfo: 7,
}

// GetIMInfoRequest is sent when information about is PLM is requested.
type GetIMInfoRequest struct{}

func (GetIMInfoRequest) commandCode() CommandCode { return GetIMInfo }

func (GetIMInfoRequest) write(io.Writer) error { return nil }

// GetIMInfoResponse is returned when information about is PLM is requested.
type GetIMInfoResponse struct {
	IMInfo IMInfo
}

func (*GetIMInfoResponse) commandCode() CommandCode { return GetIMInfo }
func (res *GetIMInfoResponse) unmarshal(r io.Reader) error {
	buffer := make([]byte, 7)

	_, err := io.ReadFull(r, buffer)

	if err != nil {
		return err
	}

	if buffer[len(buffer)-1] != MessageAck {
		return ErrCommandFailure
	}

	copy(res.IMInfo.Identity[:], buffer[:3])
	res.IMInfo.Category = Category{
		mainCategory: MainCategory(buffer[3]),
		subCategory:  SubCategory(buffer[4]),
	}
	res.IMInfo.FirmwareVersion = buffer[5]

	return nil
}
