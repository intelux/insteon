package insteon

import (
	"encoding"
	"fmt"
	"io"
)

// MessageEncoder writes messages to an io.Writer.
type MessageEncoder struct {
	Writer io.Writer
}

// Encode a value to a writer.
func (e *MessageEncoder) Encode(value interface{}) error {
	if value == nil {
		return nil
	}

	if value, ok := value.(encoding.BinaryMarshaler); ok {
		data, err := value.MarshalBinary()

		if err != nil {
			return err
		}

		_, err = e.Writer.Write(data)

		return err
	}

	panic(fmt.Errorf("unsupported value: %#v", value))
}

// NewMessageEncoder instantiates a new message encoder.
func NewMessageEncoder(w io.Writer) *MessageEncoder {
	return &MessageEncoder{Writer: w}
}
