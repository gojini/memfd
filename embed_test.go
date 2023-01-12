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
		t.Fail()
	}

	if name != "" {
		t.Fail()
	}

	if closer != nil {
		t.Fail()
	}

	name, closer, err = memfd.LoadEmbedFile(testContent, "README.md")
	if err != nil {
		t.Fail()
	}

	if !strings.HasPrefix(name, "/proc/self/fd/") {
		t.Fail()
	}

	if closer == nil {
		t.Fail()
	}

	if e := closer(); e != nil {
		t.Error(e)
	}
}
