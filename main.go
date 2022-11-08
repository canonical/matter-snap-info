package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	configFile = "config.json"
)

type config struct {
	Snaps map[string]snap
}

type snap struct {
}

func main() {
	var err error
	log.Println("Reading config file:", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Error opening config file: %s\n", err)
	}
	defer file.Close()

	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Name", "Channel", "Version", "Arch", "Revision", "Date"})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
	})

	for k, _ := range conf.Snaps {
		info, err := querySnapStore(k)
		if err != nil {
			log.Fatalf("Error querying snap store: %s", err)
		}
		for _, cm := range info.ChannelMap {
			t.AppendRow(table.Row{
				k,
				cm.Channel.Track + "/" + cm.Channel.Risk,
				cm.Version,
				cm.Channel.Architecture,
				cm.Revision,
				cm.Channel.ReleasedAt.Format(time.Stamp),
			}, table.RowConfig{AutoMerge: true})
		}
		t.AppendSeparator()
	}

	t.Render()
}

type snapInfo struct {
	ChannelMap []struct {
		Channel struct {
			Architecture string
			Track, Risk  string
			ReleasedAt   time.Time `json:"released-at"`
		}
		Revision uint
		Version  string
	} `json:"channel-map"`
}

func querySnapStore(snapName string) (*snapInfo, error) {
	log.Println("Querying snap info for:", snapName)
	req, err := http.NewRequest(http.MethodGet, "https://api.snapcraft.io/v2/snaps/info/"+snapName, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Snap-Device-Series": {"16"},
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var info snapInfo
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	// log.Println("Snap info:", info)

	return &info, nil
}
