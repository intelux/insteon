package plm

import "io"

// Request is the interface for all requests.
type Request interface{}

// Response is the interface for all responses.
type Response interface {
}

// ResponseFailure indicates that a previous message failed.
type ResponseFailure struct{}

// ParseResponse parses a response from a reader.
func ParseResponse(r io.Reader) (Response, error) {
	mark, err := skipToMessage(r)

	if err != nil {
		return nil, err
	}

	if mark == MessageNak {
		return ResponseFailure{}, nil
	}

	buffer := make([]byte, 16)

	n, err := r.Read(buffer)

	if err != nil {
		return nil, err
	}

	buffer = buffer[:n]

	return nil, err
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
