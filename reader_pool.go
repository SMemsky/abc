package abc

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func Read(r io.Reader) (*ABCFile, error) {
	abc := &ABCFile{}
	if err := binary.Read(r, binary.LittleEndian, &abc.Minor); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &abc.Major); err != nil {
		return nil, err
	}
	if abc.Major != 46 {
		return nil, fmt.Errorf("Unsupported ABC version: %d.%d", abc.Major, abc.Minor)
	}

	if err := abc.Constants.readConstantPool(r); err != nil {
		return nil, err
	}

	methodCount, err := readU30(r)
	if err != nil {
		return nil, err
	}
	abc.Methods = make([]MethodInfo, methodCount)
	for i := uint32(0); i < methodCount; i++ {
		if err := abc.Methods[i].readMethods(r); err != nil {
			return nil, err
		}
	}

	metadataCount, err := readU30(r)
	if err != nil {
		return nil, err
	}
	abc.Metadata = make([]MetadataInfo, metadataCount)
	for i := uint32(0); i < metadataCount; i++ {
		if err := abc.Metadata[i].readMetadata(r); err != nil {
			return nil, err
		}
	}

	classCount, err := readU30(r)
	if err != nil {
		return nil, err
	}
	abc.Instances = make([]InstanceInfo, classCount)
	for i := uint32(0); i < classCount; i++ {
		if err := abc.Instances[i].readInstance(r); err != nil {
			return nil, err
		}
	}
	abc.Classes = make([]ClassInfo, classCount)
	for i := uint32(0); i < classCount; i++ {
		if err := abc.Classes[i].readClass(r); err != nil {
			return nil, err
		}
	}

	scriptCount, err := readU30(r)
	if err != nil {
		return nil, err
	}
	abc.Scripts = make([]ScriptInfo, scriptCount)
	for i := uint32(0); i < scriptCount; i++ {
		if err := abc.Scripts[i].readScript(r); err != nil {
			return nil, err
		}
	}

	bodyCount, err := readU30(r)
	if err != nil {
		return nil, err
	}
	abc.MethodBodies = make([]MethodBodyInfo, bodyCount)
	for i := uint32(0); i < bodyCount; i++ {
		if err := abc.MethodBodies[i].readBody(r); err != nil {
			return nil, err
		}
	}

	return abc, nil
}

func (p *ConstantPool) readConstantPool(r io.Reader) error {
	if err := p.readSignedIntegers(r); err != nil {
		return err
	}
	if err := p.readUnsignedIntegers(r); err != nil {
		return err
	}
	if err := p.readDoubles(r); err != nil {
		return err
	}
	if err := p.readStrings(r); err != nil {
		return err
	}
	if err := p.readNamespaces(r); err != nil {
		return err
	}
	if err := p.readNamespaceSets(r); err != nil {
		return err
	}
	if err := p.readMultinames(r); err != nil {
		return err
	}
	return nil
}

func (p *ConstantPool) readSignedIntegers(r io.Reader) error {
	integerCount, err := readU30(r)
	if err != nil {
		return err
	}
	if integerCount == 0 {
		integerCount++
	}
	p.SignedIntegers = make([]int32, integerCount)
	for i := uint32(1); i < integerCount; i++ {
		p.SignedIntegers[i], err = readS32(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ConstantPool) readUnsignedIntegers(r io.Reader) error {
	integerCount, err := readU30(r)
	if err != nil {
		return err
	}
	if integerCount == 0 {
		integerCount++
	}
	p.UnsignedIntegers = make([]uint32, integerCount)
	for i := uint32(1); i < integerCount; i++ {
		p.UnsignedIntegers[i], err = readU32(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ConstantPool) readDoubles(r io.Reader) error {
	doublesCount, err := readU30(r)
	if err != nil {
		return err
	}
	if doublesCount == 0 {
		doublesCount++
	}
	p.Doubles = make([]float64, doublesCount)
	p.Doubles[0] = math.NaN()
	for i := uint32(1); i < doublesCount; i++ {
		if err := binary.Read(r, binary.LittleEndian, &p.Doubles[i]); err != nil {
			return err
		}
	}
	return nil
}

func (p *ConstantPool) readStrings(r io.Reader) error {
	stringsCount, err := readU30(r)
	if err != nil {
		return err
	}
	if stringsCount == 0 {
		stringsCount++
	}
	p.Strings = make([]string, stringsCount)
	for i := uint32(1); i < stringsCount; i++ {
		stringSize, err := readU30(r)
		if err != nil {
			return err
		}
		buffer := make([]byte, stringSize)
		if err := binary.Read(r, binary.LittleEndian, buffer); err != nil {
			return err
		}
		p.Strings[i] = string(buffer)
	}
	return nil
}

func (p *ConstantPool) readNamespaces(r io.Reader) error {
	namespaceCount, err := readU30(r)
	if err != nil {
		return err
	}
	if namespaceCount == 0 {
		namespaceCount++
	}
	p.Namespaces = make([]NamespaceInfo, namespaceCount)
	for i := uint32(1); i < namespaceCount; i++ {
		if err := p.Namespaces[i].read(r); err != nil {
			return err
		}
	}
	return nil
}

func (p *ConstantPool) readNamespaceSets(r io.Reader) error {
	setCount, err := readU30(r)
	if err != nil {
		return err
	}
	if setCount == 0 {
		setCount++
	}
	p.NamespaceSets = make([][]uint32, setCount)
	for i := uint32(1); i < setCount; i++ {
		namespaceCount, err := readU30(r)
		if err != nil {
			return err
		}
		p.NamespaceSets[i] = make([]uint32, namespaceCount)
		for j := uint32(0); j < namespaceCount; j++ {
			p.NamespaceSets[i][j], err = readU30(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *ConstantPool) readMultinames(r io.Reader) error {
	multinameCount, err := readU30(r)
	if err != nil {
		return err
	}
	if multinameCount == 0 {
		multinameCount++
	}
	p.Multinames = make([]MultinameInfo, multinameCount)
	for i := uint32(1); i < multinameCount; i++ {
		fmt.Println(i)
		if err := p.Multinames[i].read(r); err != nil {
			return err
		}
	}
	return nil
}

func (n *NamespaceInfo) read(r io.Reader) error {
	var kind uint8
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return err
	}
	n.Kind = NamespaceKind(kind)
	switch n.Kind {
	case Namespace, PackageNamespace, PackageInternalNamespace,
		ProtectedNamespace, ExplicitNamespace, StaticProtectedNamespace,
		PrivateNamespace:
		// Fine.
	default:
		return ErrBadNamespace
	}
	var err error
	n.Name, err = readU30(r)
	if err != nil {
		return err
	}
	return nil
}

func (n *MultinameInfo) read(r io.Reader) error {
	var kind uint8
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return err
	}
	fmt.Println("KIND:", kind)
	n.Kind = MultinameKind(kind)
	var err error
	switch n.Kind {
	case QName, QNameA:
		n.Namespace, err = readU30(r)
		if err != nil {
			return err
		}
		n.Name, err = readU30(r)
		if err != nil {
			return err
		}
	case RTQName, RTQNameA:
		n.Name, err = readU30(r)
		if err != nil {
			return err
		}
	case RTQNameL, RTQNameLA:
	case Multiname, MultinameA:
		n.Name, err = readU30(r)
		if err != nil {
			return err
		}
		n.NamespaceSet, err = readU30(r)
		if err != nil {
			return err
		}
	case MultinameL, MultinameLA:
		n.NamespaceSet, err = readU30(r)
		if err != nil {
			return err
		}
	case GenericName:
		n.Type, err = readU30(r)
		if err != nil {
			return err
		}
		paramCount, err := readU30(r)
		if err != nil {
			return err
		}
		n.Params = make([]uint32, paramCount)
		for i := uint32(0); i < paramCount; i++ {
			n.Params[i], err = readU30(r)
			if err != nil {
				return err
			}
		}
	default:
		fmt.Println(n.Kind)
		return ErrBadMultiname
	}
	return nil
}
