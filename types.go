package abc

type ABCFile struct {
	Minor uint16
	Major uint16 // Major 46 supported

	Constants ConstantPool
	Methods      []MethodInfo
	// Metadata     []MetadataInfo
	// Instances    []InstanceInfo
	// Classes      []ClassInfo
	// Scripts      []ScriptInfo
	// MethodBodies []MethodBodyInfo
}

// Each slice holds an additional default value with index 0. This value is not
// present in the binary file.
type ConstantPool struct {
	SignedIntegers   []int32   // Default 0.
	UnsignedIntegers []uint32  // Default 0.
	Doubles          []float64 // Default NaN.
	Strings          []string  // Default "", interpreted as "any". UTF-8.

	Namespaces    []NamespaceInfo // Default is "any".
	NamespaceSets [][]uint32      // Default is not defined. Non-zero indexes into Namespace slice.
	Multinames    []MultinameInfo // Default is not defined.
}

type NamespaceKind uint

const (
	Namespace                NamespaceKind = 0x08
	PackageNamespace         NamespaceKind = 0x16
	PackageInternalNamespace NamespaceKind = 0x17
	ProtectedNamespace       NamespaceKind = 0x18
	ExplicitNamespace        NamespaceKind = 0x19
	StaticProtectedNamespace NamespaceKind = 0x1A
	PrivateNamespace         NamespaceKind = 0x05
)

type NamespaceInfo struct {
	Kind NamespaceKind
	Name uint32 // Strings. Zero is an empty string.
}

type MultinameKind uint

const (
	QName       MultinameKind = 0x07
	QNameA      MultinameKind = 0x0D
	RTQName     MultinameKind = 0x0F
	RTQNameA    MultinameKind = 0x10
	RTQNameL    MultinameKind = 0x11
	RTQNameLA   MultinameKind = 0x12
	Multiname   MultinameKind = 0x09
	MultinameA  MultinameKind = 0x0E
	MultinameL  MultinameKind = 0x1B
	MultinameLA MultinameKind = 0x1C
	// TODO: https://richardszalay.wordpress.com/2009/02/11/generics-vector-in-the-avm2/
	GenericName MultinameKind = 0x1D
)

type MultinameInfo struct {
	Kind MultinameKind

	// QName, QNameA.
	Namespace uint32 // Namespaces. Zero = any.

	// QName, QNameA, RTQName, RTQNameA, Multiname, MultinameA.
	Name uint32 // Strings. Zero = any.

	// Multiname, MultinameA, MultinameL, MultinameLA.
	NamespaceSet uint32 // NamespaceSets. Zero = any for Multiname. Non-zero for MultinameL.

	// GenericName.
	Type   uint32   // Multinames.
	Params []uint32 // Multinames.
}

type MethodFlags uint8

const (
    NeedsArguments MethodFlags = 0x01 // Needs an arguments object. Conflicts with NeedsRest.
    UsesActivation MethodFlags = 0x02 // Uses newactivation opcode.
    NeedsRest      MethodFlags = 0x04 // Needs rest arguments array. Conflicts with NeedsArguments.
    HasOptional    MethodFlags = 0x08 // Has optional parameters.
    UsesDXNS       MethodFlags = 0x40 // Uses dxns and dxnslate opcodes.
    HasParamNames  MethodFlags = 0x80 // Has argument names.
)

type ValueKind uint8

const (
    SignedInteger   ValueKind = 0x03
    UnsignedInteger ValueKind = 0x04
    Double          ValueKind = 0x06
    String          ValueKind = 0x01
    True            ValueKind = 0x0B
    False           ValueKind = 0x0A
    Null            ValueKind = 0x0C
    Undefined       ValueKind = 0x00
)

type MethodInfo struct {
    ReturnType    uint32   // Multinames. Zero = any.
    ArgumentTypes []uint32 // Multinames. Zero = any.
    Name          uint32   // Strings. Zero = no name.

    Flags MethodFlags

    // Only for HasOptional.
    Defaults []ParamDefaults
    // Only for HasParamNames.
    Names []uint32 // Strings. Argument names.
}

type ParamDefaults struct {
    Kind  uint8  // Either ValueKind, or NamespaceKind.
    Index uint32 // Index into corresponding Pool entry.
}
