package abc

import (
	"encoding/binary"
	"fmt"
	"io"
)

func (t *TraitInfo) readTrait(r io.Reader) error {
	var err error
	if t.Name, err = readU30(r); err != nil {
		return err
	}
	var kind uint8
	if err = binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return err
	}
	if (kind & 0x80) != 0 {
		return ErrUnknownTraitAttribute
	}
	t.Attributes = TraitAttributes(kind >> 4)
	t.Kind = TraitKind(kind & 0xf)

	switch t.Kind {
	case SlotTrait, ConstTrait:
		if t.SlotId, err = readU30(r); err != nil {
			return err
		}
		if t.Type, err = readU30(r); err != nil {
			return err
		}
		if t.ValueIndex, err = readU30(r); err != nil {
			return err
		}
		if t.ValueIndex != 0 {
			if err = binary.Read(r, binary.LittleEndian, &t.ValueKind); err != nil {
				return err
			}
		}
	case ClassTrait:
		if t.SlotId, err = readU30(r); err != nil {
			return err
		}
		if t.Class, err = readU30(r); err != nil {
			return err
		}
	case FunctionTrait:
		if t.SlotId, err = readU30(r); err != nil {
			return err
		}
		if t.Function, err = readU30(r); err != nil {
			return err
		}
	case MethodTrait, GetterTrait, SetterTrait:
		if t.DisplayId, err = readU30(r); err != nil {
			return err
		}
		if t.Method, err = readU30(r); err != nil {
			return err
		}
	default:
		fmt.Println(t.Kind, kind)
		return ErrUnknownTraitKind
	}

	if (t.Attributes & MetadataAttribute) != 0 {
		metadataCount, err := readU30(r)
		if err != nil {
			return err
		}
		t.Metadata = make([]uint32, metadataCount)
		for i := uint32(0); i < metadataCount; i++ {
			if t.Metadata[i], err = readU30(r); err != nil {
				return err
			}
		}
	}

	return nil
}
