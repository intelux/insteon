package plm

import "io"

// AllLinkRecordResponse is returned when an all-link record is sent.
type AllLinkRecordResponse struct {
	Record AllLinkRecord
}

func (*AllLinkRecordResponse) commandCode() CommandCode { return AllLinkRecordMessage }
func (res *AllLinkRecordResponse) unmarshal(r io.Reader) error {
	buffer := make([]byte, 2)

	_, err := io.ReadFull(r, buffer)

	if err != nil {
		return err
	}

	res.Record.Flags = AllLinkRecordFlags(buffer[0])
	res.Record.Group = Group(buffer[1])

	_, err = io.ReadFull(r, res.Record.Identity[:])

	if err != nil {
		return err
	}

	_, err = io.ReadFull(r, res.Record.LinkData[:])

	if err != nil {
		return err
	}

	return nil
}
