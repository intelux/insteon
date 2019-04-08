package plm

import "io"

// GetNextAllLinkRecordRequest is sent when the next all-link record is requested.
type GetNextAllLinkRecordRequest struct{}

func (GetNextAllLinkRecordRequest) commandCode() CommandCode { return GetNextAllLinkRecord }
func (GetNextAllLinkRecordRequest) marshal(io.Writer) error  { return nil }

// GetNextAllLinkRecordResponse is returned the next all-link record is about to be sent.
type GetNextAllLinkRecordResponse struct{}

func (*GetNextAllLinkRecordResponse) commandCode() CommandCode { return GetNextAllLinkRecord }
func (res *GetNextAllLinkRecordResponse) unmarshal(r io.Reader) error {
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
