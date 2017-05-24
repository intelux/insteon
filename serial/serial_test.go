package serial

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

func hasSocat() bool {
	_, err := exec.LookPath("socat")

	return err == nil
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

func TestOpen(t *testing.T) {
	if !hasSocat() {
		t.Skipf("`socat` is necessary for this test but was not found")
	}

	device := "mydevice"
	stop, err := runEchoSerialPort(device)

	if err != nil {
		t.Skipf("failed to create the echo serial port: %s", err)
	}

	defer close(stop)

	port, err := Open(device)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	sent := []byte("!")
	_, err = port.Write(sent)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	received := make([]byte, 10)
	n, err := port.Read(received)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	received = received[:n]

	if string(sent) != string(received) {
		t.Errorf("expected received data to be %s, but was %s", sent, received)
	}
}

func TestOpenFailure(t *testing.T) {
	device := "nonexistingdevice"

	_, err := Open(device)

	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestOpenInvalidDeviceType(t *testing.T) {
	device := "fixtures/device"

	_, err := Open(device)

	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}
