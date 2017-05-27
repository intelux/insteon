package plm

import "io"

// GetIMInfoRequest is sent when information about is PLM is requested.
type GetIMInfoRequest struct{}

func (GetIMInfoRequest) commandCode() CommandCode { return GetIMInfo }

func (GetIMInfoRequest) marshal(io.Writer) error { return nil }

// GetIMInfoResponse is returned when information about is PLM is requested.
type GetIMInfoResponse struct {
	IMInfo IMInfo
}

func (*GetIMInfoResponse) commandCode() CommandCode { return GetIMInfo }
func (res *GetIMInfoResponse) unmarshal(r io.Reader) error {
	buffer := make([]byte, 7)

	_, err := io.ReadFull(r, buffer)

	if err != nil {
		return err
	}

	if buffer[len(buffer)-1] != MessageAck {
		return ErrCommandFailure
	}

	copy(res.IMInfo.Identity[:], buffer[:3])
	res.IMInfo.Category = Category{
		mainCategory: MainCategory(buffer[3]),
		subCategory:  SubCategory(buffer[4]),
	}
	res.IMInfo.FirmwareVersion = buffer[5]

	return nil
}
