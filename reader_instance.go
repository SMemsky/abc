package abc

import (
	"encoding/binary"
	"fmt"
	"io"
)

func (i *InstanceInfo) readInstance(r io.Reader) error {
	var err error
	i.Name, err = readU30(r)
	if err != nil {
		return err
	}
	i.BaseName, err = readU30(r)
	if err != nil {
		return err
	}
	var flags uint8
	if err := binary.Read(r, binary.LittleEndian, &flags); err != nil {
		return err
	}
	if (flags & 0xf0) != 0 {
		return ErrUnknownInstanceFlag
	}
	i.Flags = ClassFlags(flags)
	i.ProtectedNamespace, err = readU30(r)
	if err != nil {
		return err
	}

	interfaceCount, err := readU30(r)
	if err != nil {
		return err
	}
	i.Interfaces = make([]uint32, interfaceCount)
	for j := uint32(0); j < interfaceCount; j++ {
		i.Interfaces[j], err = readU30(r)
		if err != nil {
			return err
		}
	}

	i.InstanceInit, err = readU30(r)
	if err != nil {
		return err
	}

	traitCount, err := readU30(r)
	if err != nil {
		return err
	}
	i.Traits = make([]TraitInfo, traitCount)
	for j := uint32(0); j < traitCount; j++ {
		fmt.Println("reading trait", j+1, "/", traitCount)
		if err := i.Traits[j].readTrait(r); err != nil {
			return err
		}
	}

	return nil
}
