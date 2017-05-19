package plm

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNewDialError(t *testing.T) {
	_, err := New("tcp://0.0.0.0:0")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestNewOpenError(t *testing.T) {
	_, err := New("nonexistingdevice")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestNewUnsupportedScheme(t *testing.T) {
	_, err := New("foo://bar")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestNewDialSuccess(t *testing.T) {
	srv, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		t.Fatalf("failed to start tcp server: %s", err)
	}

	defer srv.Close()

	plm, err := New("tcp://" + srv.Addr().String())

	if err != nil {
		t.Fatalf("failed to create device: %s", err)
	}

	if plm == nil {
		t.Fatalf("plm is not supposed to be nil")
	}

	defer plm.Close()
}

func runEchoSerialPort(device string) (chan struct{}, error) {
	socat := exec.Command("socat", fmt.Sprintf("pty,raw,link=%s", device), "-")
	stdout, err := socat.StdoutPipe()

	if err != nil {
		return nil, err
	}

	stdin, err := socat.StdinPipe()

	if err != nil {
		return nil, err
	}

	go io.Copy(stdin, stdout)

	err = socat.Start()

	if err != nil {
		return nil, err
	}

	stop := make(chan struct{})

	go func() {
		<-stop
		socat.Process.Signal(os.Interrupt)
		socat.Wait()
	}()

	// Wait for the device to exist.
	for _, err = os.Stat(device); os.IsNotExist(err); _, err = os.Stat(device) {
		time.Sleep(time.Millisecond)
	}

	return stop, nil
}

func TestNewLocalDevice(t *testing.T) {
	device := "mydevice"
	stop, err := runEchoSerialPort(device)

	if err != nil {
		t.Fatalf("failed to create serial port: %s", err)
	}

	defer close(stop)

	plm, err := New(device)

	if err != nil {
		t.Fatalf("failed to create device: %s", err)
	}

	if plm == nil {
		t.Fatalf("plm is not supposed to be nil")
	}

	defer plm.Close()
}
