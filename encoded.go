package abc

import (
	"encoding/binary"
	"io"
)

// Support for encoded U30, U32, S32 values.

func readU30(r io.Reader) (uint32, error) {
	value, err := readU32(r)
	if err != nil {
		return 0, err
	}
	if value >= 0x40000000 {
		return 0, ErrBadEncoded
	}
	return value, nil
}

func readU32(r io.Reader) (uint32, error) {
	var result uint32
	var buffer uint8
	for i := uint(0); i < 4; i++ {
		if err := binary.Read(r, binary.LittleEndian, &buffer); err != nil {
			return 0, err
		}
		result += uint32(buffer&0x7f) << (7 * i)
		if (buffer & 0x80) == 0 {
			return result, nil
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &buffer); err != nil {
		return 0, err
	}
	if (buffer & 0xf0) != 0 {
		return 0, ErrBadEncoded
	}
	return result + (uint32(buffer&0xf) << 28), nil
}

func readS32(r io.Reader) (int32, error) {
	var result int32
	var buffer uint8
	for i := uint(0); i < 4; i++ {
		if err := binary.Read(r, binary.LittleEndian, &buffer); err != nil {
			return 0, err
		}
		result += int32(buffer&0x7f) << (7 * i)
		if (buffer & 0x80) == 0 {
			if (buffer & 0x40) == 1 {
				result = -result
			}
			return result, nil
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &buffer); err != nil {
		return 0, err
	}
	if (buffer & 0xf0) != 0 {
		return 0, ErrBadEncoded
	}
	result += int32(buffer&0xf) << 28
	if (buffer & 0x40) == 1 {
		result = -result
	}
	return result, nil
}
