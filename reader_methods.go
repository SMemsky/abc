package abc

import (
    "encoding/binary"
    // "fmt"
    "io"
)

func (m *MethodInfo) readMethods(r io.Reader) error {
    paramCount, err := readU30(r)
    if err != nil {
        return err
    }
    m.ArgumentTypes = make([]uint32, paramCount)
    m.ReturnType, err = readU30(r)
    if err != nil {
        return err
    }
    for i := uint32(0); i < paramCount; i++ {
        m.ArgumentTypes[i], err = readU30(r)
        if err != nil {
            return err
        }
    }
    m.Name, err = readU30(r)
    if err != nil {
        return err
    }

    var flags uint8
    if err := binary.Read(r, binary.LittleEndian, &flags); err != nil {
        return err
    }
    m.Flags = MethodFlags(flags)

    if (m.Flags & NeedsArguments) != 0 || (m.Flags & NeedsRest) != 0 {
        if (m.Flags & NeedsArguments) == (m.Flags & NeedsRest) {
            return ErrRestArguments
        }
    }

    if (m.Flags & HasOptional) != 0 {
        optionCount, err := readU30(r)
        if err != nil {
            return err
        }
        m.Defaults = make([]ParamDefaults, optionCount)
        for i := uint32(0); i < optionCount; i++ {
            if err := m.Defaults[i].read(r); err != nil {
                return err
            }
        }
    }

    if (m.Flags & HasParamNames) != 0 {
        m.Names = make([]uint32, paramCount)
        for i := uint32(0); i < paramCount; i++ {
            m.Names[i], err = readU30(r)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (d *ParamDefaults) read(r io.Reader) error {
    var err error
    d.Index, err = readU30(r)
    if err != nil {
        return err
    }
    if err := binary.Read(r, binary.LittleEndian, &d.Kind); err != nil {
        return err
    }

    switch ValueKind(d.Kind) {
    case SignedInteger, UnsignedInteger, Double, String, True, False, Null,
        Undefined:
        // Fine.
    default:
        // Try Namespace.
        switch NamespaceKind(d.Kind) {
        case Namespace, PackageNamespace, PackageInternalNamespace,
            ProtectedNamespace, ExplicitNamespace, StaticProtectedNamespace,
            PrivateNamespace:
            // Fine.
        default:
            return ErrBadOptionKind
        }
    }

    return nil
}
