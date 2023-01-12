package memfd

import (
	"embed"
	"fmt"
	"io"
)

func LoadEmbedFile(fs embed.FS, srcFile string) (string, func() error, error) {
	f, err := fs.Open(srcFile)
	if err != nil {
		return "", nil, fmt.Errorf("unable to open source: %w", err)
	}

	defer f.Close()

	mfd, err := Create()
	if err != nil {
		return "", nil, fmt.Errorf("unable to create memfd: %w", err)
	}

	defer mfd.Unmap()

	_, err = io.Copy(mfd, f)
	if err != nil {
		mfd.Close()
		return "", nil, fmt.Errorf("unable to load to memfd: %w", err)
	}

	return fmt.Sprintf("/proc/self/fd/%d", mfd.Fd()), mfd.Close, nil
}
