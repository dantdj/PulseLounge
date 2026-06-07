package logging

import (
	"io"
	"os"
)

var defaultOutput io.Writer = os.Stdout
