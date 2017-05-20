package plm

import "io"

// Response is the interface for all responses.
type Response interface{}

// ParseResponse parses a response from a reader.
func ParseResponse(r io.Reader) (Response, error) {
	mark, err := skipToMessage(r)

	if err != nil {
		return nil, err
	}

	if mark == MessageFailure {
		// TODO: Return failure message
		return nil, nil
	}

	buffer := make([]byte, 16)

	n, err := r.Read(buffer)

	if err != nil {
		return nil, err
	}
}

// skipToMessage skips bytes on the specified io.Reader until a message start
// or message failure is met.
//
// The message start or message failure is returned.
func skipToMessageStart(r io.Reader) (byte, error) {
	buffer := make([]byte, 1)

	for buffer[0] != Message && buffer[0] != MessageFailure {
		_, err := r.Read(buffer)

		if err != nil {
			return nil, err
		}
	}

	return buffer[0], nil
}
