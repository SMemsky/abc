package abc

import (
	"io"
)

func (s *ScriptInfo) readScript(r io.Reader) error {
	var err error
	if s.Init, err = readU30(r); err != nil {
		return err
	}

	traitCount, err := readU30(r)
	if err != nil {
		return err
	}
	s.Traits = make([]TraitInfo, traitCount)
	for j := uint32(0); j < traitCount; j++ {
		if err := s.Traits[j].readTrait(r); err != nil {
			return err
		}
	}

	return nil
}
