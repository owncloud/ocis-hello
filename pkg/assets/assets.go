package assets

import (
	"net/http"
	"os"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

//go:generate go run github.com/UnnoTed/fileb0x embed.yml

// assets gets initialized by New and provides the handler.
type assets struct {
	logger log.Logger
	path   string
}

// Open just implements the HTTP filesystem interface.
func (a assets) Open(original string) (http.File, error) {
	if a.path != "" {
		if stat, err := os.Stat(a.path); err == nil && stat.IsDir() {
			custom := path.Join(
				a.path,
				original,
			)

			if _, err := os.Stat(custom); !os.IsNotExist(err) {
				f, err := os.Open(custom)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			level.Warn(a.logger).Log(
				"msg", "Assets directory doesn't exist",
				"path", a.path,
			)
		}
	}

	return FS.OpenFile(
		CTX,
		original,
		os.O_RDONLY,
		0644,
	)
}

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	a := new(assets)

	for _, opt := range opts {
		opt(a)
	}

	return a
}
