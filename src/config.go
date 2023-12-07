package config

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Snaps map[string]struct {
		GithubRepo string
	}
}

func LoadConfig(confFile string) (c *Config, err error) {
	var reader io.Reader

	if strings.HasPrefix(confFile, "http") {
		log.Println("Fetching config file from:", confFile)

		res, err := http.Get(confFile)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		reader = res.Body
	} else {
		log.Println("Reading local config file from:", confFile)
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
