package config

import (
	"encoding/json"
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

	if strings.HasPrefix(confFile, "http") {
		log.Println("Fetching config file from:", confFile)

		res, err := http.Get(confFile)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&c)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Reading local config file from:", confFile)
		file, err := os.Open(confFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(&c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
