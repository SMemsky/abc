package abc

import (
	"fmt"
)

var (
	ErrBadEncoded   = fmt.Errorf("Malformed encoded value.")
	ErrBadNamespace = fmt.Errorf("Unknown namespace kind.")
	ErrBadMultiname = fmt.Errorf("Unknown multiname kind.")
    ErrRestArguments = fmt.Errorf("Both NEED_REST and NEED_ARGUMENTS specified.")
    ErrBadOptionKind = fmt.Errorf("Unknown option kind.")
)
