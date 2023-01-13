package memfd_test

import (
	"embed"
	"strings"
	"testing"

	"gojini.dev/memfd"
)

//go:embed README.md
var testContent embed.FS

func TestLoadFromEmbedFS(t *testing.T) {
	t.Parallel()

	name, closer, err := memfd.LoadEmbedFile(testContent, "randomFile")
	if err == nil {
		t.Error("missing embedded file was loaded successfully")
	}

	if name != "" {
		t.Errorf("name was not empty: %s", name)
	}

	if closer != nil {
		t.Error("closer was non-nil")
	}

	name, closer, err = memfd.LoadEmbedFile(testContent, "README.md")
	if err != nil {
		t.Errorf("failed to load embedded file: %s", err)
	}

	if !strings.HasPrefix(name, "/proc/self/fd/") {
		t.Errorf("incorrect prefix in the path: %s", name)
	}

	if closer == nil {
		t.Error("closer was nil")
	}

	if err := closer(); err != nil {
		t.Errorf("failed to close: %s", err)
	}
}
