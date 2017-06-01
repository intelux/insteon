package plm

import "io"

// GetFirstAllLinkRecordRequest is sent when the first all-link record is requested.
type GetFirstAllLinkRecordRequest struct{}

func (GetFirstAllLinkRecordRequest) commandCode() CommandCode { return GetFirstAllLinkRecord }
func (GetFirstAllLinkRecordRequest) marshal(io.Writer) error  { return nil }

// GetFirstAllLinkRecordResponse is returned when information about is PLM is requested.
type GetFirstAllLinkRecordResponse struct{}

func (*GetFirstAllLinkRecordResponse) commandCode() CommandCode { return GetFirstAllLinkRecord }
func (res *GetFirstAllLinkRecordResponse) unmarshal(r io.Reader) error {
	buffer := make([]byte, 1)

	_, err := io.ReadFull(r, buffer)

	if err != nil {
		return err
	}

	if buffer[len(buffer)-1] != MessageAck {
		return ErrCommandFailure
	}

	return nil
}
