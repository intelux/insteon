package plm

import (
	"bytes"
	"errors"
	"testing"
)

func TestDebugStreamLogError(t *testing.T) {
	buf := bytes.Buffer{}
	stream := debugStream{
		Writer: &buf,
		Local:  "local",
		Remote: "remote",
	}
	stream.Log("->", nil, errors.New("error"))
	expected := "[local -> remote] error\n"

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, string(buf.Bytes()))
	}
}

func TestDebugStreamLog(t *testing.T) {
	buf := bytes.Buffer{}
	stream := debugStream{
		Writer: &buf,
		Local:  "local",
		Remote: "remote",
	}
	stream.Log("->", []byte{0x01, 0x02}, nil)
	expected := "[local -> remote] 0102\n"

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, string(buf.Bytes()))
	}
}

func TestDebugReader(t *testing.T) {
	buf := bytes.Buffer{}
	reader := DebugReader{
		Reader: bytes.NewBuffer([]byte{0x01, 0x02}),
		debugStream: debugStream{
			Writer: &buf,
		},
	}
	data := make([]byte, 5)
	n, err := reader.Read(data)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if n != 2 {
		t.Errorf("expected n to be %d, not %d", 2, n)
	}

	expected := "[ < ] 0102\n"

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, string(buf.Bytes()))
	}
}

func TestDebugWriter(t *testing.T) {
	buf := bytes.Buffer{}
	writer := DebugWriter{
		Writer: &bytes.Buffer{},
		debugStream: debugStream{
			Writer: &buf,
		},
	}
	n, err := writer.Write([]byte{0x01, 0x02})

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if n != 2 {
		t.Errorf("expected n to be %d, not %d", 2, n)
	}

	expected := "[ > ] 0102\n"

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, string(buf.Bytes()))
	}
}
