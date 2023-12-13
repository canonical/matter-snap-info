package config

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/canonical/matter-snap-info/logger"
)

type Config struct {
	Snaps map[string]struct {
		GithubRepo string
	}
}

func Load(path string) (c *Config, err error) {
	var reader io.Reader

	if strings.HasPrefix(path, "http") {
		logger.Infoln("Fetching config file from:", path)

		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		reader = res.Body
	} else {
		logger.Infoln("Reading local config file from:", path)
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader = file
	}

	err = json.NewDecoder(reader).Decode(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
