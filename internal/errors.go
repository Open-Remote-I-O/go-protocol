package errors

import (
	"fmt"
	"io"
)

var (
	ErrHeaderFormat    = "invalid protocol header format sent"
	ErrHeaderFormatEOF = fmt.Errorf("%s: %w", ErrHeaderFormat, io.EOF)
)
