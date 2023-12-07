package config

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	logger "github.com/canonical/edgex-snap-info/src/log"
)

type Config struct {
	Snaps map[string]struct {
		GithubRepo string
	}
}

func LoadConfig(confFile string) (c *Config, err error) {
	var reader io.Reader

	if strings.HasPrefix(confFile, "http") {
		logger.Println("Fetching config file from:", logger.White, confFile)

		res, err := http.Get(confFile)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		reader = res.Body
	} else {
		logger.Println("Reading local config file from:", logger.White, confFile)
		file, err := os.Open(confFile)
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
