package plm

import "io"

// AllLinkRecordResponse is returned when an all-link record is sent.
type AllLinkRecordResponse struct {
	Flags    AllLinkRecordFlags
	Group    Group
	Identity Identity
	LinkData LinkData
}

func (*AllLinkRecordResponse) commandCode() CommandCode { return AllLinkRecord }
func (res *AllLinkRecordResponse) unmarshal(r io.Reader) error {
	buffer := make([]byte, 2)

	_, err := io.ReadFull(r, buffer)

	if err != nil {
		return err
	}

	res.Flags = AllLinkRecordFlags(buffer[0])
	res.Group = Group(buffer[1])

	_, err = io.ReadFull(r, res.Identity[:])

	if err != nil {
		return err
	}

	_, err = io.ReadFull(r, res.LinkData[:])

	if err != nil {
		return err
	}

	return nil
}
