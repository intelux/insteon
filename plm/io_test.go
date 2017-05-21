package plm

import (
	"errors"
	"io"
	"testing"
)

func TestConnectedPipeUnconnectedWrite(t *testing.T) {
	pipe := &ConnectedPipe{}

	n, err := pipe.Write([]byte("hello"))

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if n != 5 {
		t.Errorf("expected n to be 5 but got: %d", n)
	}
}

func TestConnectedPipeReadFailure(t *testing.T) {
	errTest := errors.New("expected error")
	pipe := &ConnectedPipe{}
	r, w := io.Pipe()

	connectDone := make(chan struct{})

	go func() {
		defer close(connectDone)
		err := pipe.Connect(w)

		if err != errTest {
			t.Fatalf("expected error to be %s but got: %s", errTest, err)
		}
	}()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, err := pipe.Write([]byte("hi"))

				if err != nil {
					t.Fatalf("expected no error but got: %s", err)
				}
			}
		}
	}()

	buf := make([]byte, 10)
	n, err := r.Read(buf)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if n != 2 {
		t.Errorf("expected n to be 2 but got: %d", n)
	}

	// Close the writer end of the pipe, causing the reading end to return errClosedConnectedPipe.
	r.CloseWithError(errTest)

	<-connectDone
}

func TestConnectedPipeClose(t *testing.T) {
	pipe := &ConnectedPipe{}
	r, w := io.Pipe()

	connectDone := make(chan struct{})

	go func() {
		defer close(connectDone)
		err := pipe.Connect(w)

		if err != io.ErrClosedPipe {
			t.Fatalf("expected error to be %s but got: %s", io.ErrClosedPipe, err)
		}
	}()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, err := pipe.Write([]byte("hi"))

				if err != nil {
					t.Fatalf("expected no error but got: %s", err)
				}
			}
		}
	}()

	// Just make sure the connect is pending by waiting a for a read to succeed.
	buf := make([]byte, 10)
	_, _ = r.Read(buf)

	pipe.Close()

	<-connectDone
}
