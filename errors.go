package abc

import (
	"fmt"
)

var (
	ErrBadEncoded   = fmt.Errorf("Malformed encoded value.")
	ErrBadNamespace = fmt.Errorf("Malformed namespace kind.")
	ErrBadMultiname = fmt.Errorf("Malformed multiname kind.")
)
