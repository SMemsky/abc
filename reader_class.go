package abc

import (
	"io"
)

func (c *ClassInfo) readClass(r io.Reader) error {
	var err error
	if c.StaticInit, err = readU30(r); err != nil {
		return err
	}

	traitCount, err := readU30(r)
	if err != nil {
		return err
	}
	c.Traits = make([]TraitInfo, traitCount)
	for j := uint32(0); j < traitCount; j++ {
		if err := c.Traits[j].readTrait(r); err != nil {
			return err
		}
	}

	return nil
}
