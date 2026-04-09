// This file is here as go embed only supports embedding files in a structure underneath it,
// and the frontend is in a different directory.

package uiassets

import (
	"embed"
	"fmt"
	"io/fs"
)

const stagedUIPath = "frontend/embed/generated"

// embeddedUI holds the tracked embed scaffold plus staged UI assets when present.
//
//go:embed frontend/embed
var embeddedUI embed.FS

func FS() (fs.FS, error) {
	if _, err := fs.Stat(embeddedUI, stagedUIPath+"/index.html"); err != nil {
		return nil, fmt.Errorf("staged UI assets not found at %s; run `make ui-build` before building the server: %w", stagedUIPath, err)
	}

	return fs.Sub(embeddedUI, stagedUIPath)
}
