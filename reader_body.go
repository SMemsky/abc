package abc

import (
	"encoding/binary"
	"io"
)

func (b *MethodBodyInfo) readBody(r io.Reader) error {
	var err error
	if b.Method, err = readU30(r); err != nil {
		return err
	}
	if b.StackLimit, err = readU30(r); err != nil {
		return err
	}
	if b.LocalCount, err = readU30(r); err != nil {
		return err
	}
	if b.InitScopeDepth, err = readU30(r); err != nil {
		return err
	}
	if b.MaxScopeDepth, err = readU30(r); err != nil {
		return err
	}
	codeLength, err := readU30(r)
	if err != nil {
		return err
	}
	b.RawCode = make([]byte, codeLength)
	if err := binary.Read(r, binary.LittleEndian, b.RawCode); err != nil {
		return err
	}

	exceptionCount, err := readU30(r)
	if err != nil {
		return err
	}
	b.Exceptions = make([]ExceptionInfo, exceptionCount)
	for i := uint32(0); i < exceptionCount; i++ {
		if err := b.Exceptions[i].readException(r); err != nil {
			return err
		}
	}

	traitCount, err := readU30(r)
	if err != nil {
		return err
	}
	b.Traits = make([]TraitInfo, traitCount)
	for j := uint32(0); j < traitCount; j++ {
		if err := b.Traits[j].readTrait(r); err != nil {
			return err
		}
	}

	return nil
}

func (e *ExceptionInfo) readException(r io.Reader) error {
	var err error
	if e.From, err = readU30(r); err != nil {
		return err
	}
	if e.To, err = readU30(r); err != nil {
		return err
	}
	if e.Target, err = readU30(r); err != nil {
		return err
	}
	if e.ExceptionType, err = readU30(r); err != nil {
		return err
	}
	if e.VarName, err = readU30(r); err != nil {
		return err
	}
	return nil
}
