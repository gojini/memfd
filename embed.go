package memfd

import (
	"embed"
	"fmt"
	"io"
)

// LoadEmbedFile loads the specified file from the given embed file system into
// memory and returns the name of the memfd, a file closer function and error
// on failure to load the file.
//
// The name of the file is canonical representation of the file in the
// /proc/self/fd/<fd-num> format. The closer function must be called to release
// the file descriptor and the memory associated with it. It is usually a good
// idea to use defer to close this fd once in the scope of the caller/usage.
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
