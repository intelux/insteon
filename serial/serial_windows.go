package serial

import (
	"errors"
	"io"
)

func open(device string) (io.ReadWriteCloser, error) {
	return nil, errors.New("not yet implemented on Windows")
}
