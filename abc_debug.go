package abc

import (
	"fmt"
)

var namespaceKind2String = map[NamespaceKind]string{
	Namespace:                "Namespace",
	PackageNamespace:         "PackageNamespace",
	PackageInternalNamespace: "PackageInternalNamespace",
	ProtectedNamespace:       "ProtectedNamespace",
	ExplicitNamespace:        "ExplicitNamespace",
	StaticProtectedNamespace: "StaticProtectedNamespace",
	PrivateNamespace:         "PrivateNamespace",
}

var multinameKind2String = map[MultinameKind]string{
	QName:       "QName",
	QNameA:      "QNameA",
	RTQName:     "RTQName",
	RTQNameA:    "RTQNameA",
	RTQNameL:    "RTQNameL",
	RTQNameLA:   "RTQNameLA",
	Multiname:   "Multiname",
	MultinameA:  "MultinameA",
	MultinameL:  "MultinameL",
	MultinameLA: "MultinameLA",
	GenericName: "GenericName",
}

func (a *ABCFile) Dump() {
	fmt.Printf("SWF Version %d.%d\n", a.Major, a.Minor)
}

func (a *ABCFile) DumpMultiname(index uint32) {
	multiname := &a.Constants.Multinames[index]
	fmt.Printf("%s(", multinameKind2String[multiname.Kind])
	if multiname.Kind == QName {
		a.dumpNamespace(multiname.Namespace)
		fmt.Printf(", \"%s\")\n", a.Constants.Strings[multiname.Name])
	}
}

func (a *ABCFile) dumpNamespace(index uint32) {
	namespace := &a.Constants.Namespaces[index]
	fmt.Printf("%s(\"%s\")", namespaceKind2String[namespace.Kind], a.Constants.Strings[namespace.Name])
}
