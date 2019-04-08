package plm

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestParseDeviceDialError(t *testing.T) {
	_, err := ParseDevice("tcp://0.0.0.0:0")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestParseDeviceOpenError(t *testing.T) {
	_, err := ParseDevice("nonexistingdevice")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestParseDeviceUnsupportedScheme(t *testing.T) {
	_, err := ParseDevice("foo://bar")

	if err == nil {
		t.Error("expected a failure")
	}
}

func TestParseDeviceDialSuccess(t *testing.T) {
	srv, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		t.Fatalf("failed to start tcp server: %s", err)
	}

	defer srv.Close()

	device, err := ParseDevice("tcp://" + srv.Addr().String())

	if err != nil {
		t.Fatalf("failed to create device: %s", err)
	}

	if device == nil {
		t.Fatalf("device is not supposed to be nil")
	}

	defer device.Close()
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

func TestParseDeviceLocalDevice(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Not implemented yet on Windows")
	}

	name := "mydevice"
	stop, err := runEchoSerialPort(name)

	if err != nil {
		t.Fatalf("failed to create serial port: %s", err)
	}

	defer close(stop)

	device, err := ParseDevice(name)

	if err != nil {
		t.Fatalf("failed to create device: %s", err)
	}

	if device == nil {
		t.Fatalf("device is not supposed to be nil")
	}

	defer device.Close()
}
