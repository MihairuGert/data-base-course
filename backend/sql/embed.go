package sqlfiles

import (
	"embed"
	"fmt"
)

// FS stores raw SQL used by the backend.
//
//go:embed queries/*.sql reports/*.sql
var FS embed.FS

func MustRead(name string) string {
	b, err := FS.ReadFile(name)
	if err != nil {
		panic(fmt.Errorf("read sql %s: %w", name, err))
	}
	return string(b)
}
