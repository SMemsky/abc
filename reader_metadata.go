package abc

import (
	// "encoding/binary"
	// "fmt"
	"io"
)

// TODO
func (m *MetadataInfo) readMetadata(r io.Reader) error {
	var err error
	if m.Name, err = readU30(r); err != nil {
		return err
	}

	itemCount, err := readU30(r)
	if err != nil {
		return err
	}

	m.Items = make([]MetadataItemInfo, itemCount)
	for i := uint32(0); i < itemCount; i++ {
		if err := m.Items[i].readMetadataItem(r); err != nil {
			return err
		}
	}
	return nil
}

func (m *MetadataItemInfo) readMetadataItem(r io.Reader) error {
	var err error
	if m.Key, err = readU30(r); err != nil {
		return err
	}
	if m.Value, err = readU30(r); err != nil {
		return err
	}
	return nil
}
