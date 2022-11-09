package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	configURL = "https://raw.githubusercontent.com/farshidtz/edgex-snap-info/main/config.json"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Name", "Channel", "Version", "Arch", "Rev", "Date", "Build"})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
	})

	for k, v := range conf.Snaps {
		// snap store
		info, err := querySnapStore(k)
		if err != nil {
			log.Fatalf("Error querying snap store: %s", err)
		}

		// launchpad
		builds, err := queryLaunchpad(k)
		if err != nil {
			log.Fatalf("Error querying launchpad: %s", err)
		}
		builtRevs := make(map[uint]string)
		for _, v := range builds.Entries {
			if v.StoreUploadRevision != nil {
				builtRevs[*v.StoreUploadRevision] = v.BuildState
			}
		}

		// github
		runs, err := queryGithub(v.GithubRepo)
		if err != nil {
			log.Fatalf("Error querying launchpad: %s", err)
		}
		var totalSnapRuns, failedSnapRuns uint
		testIcon := "🔴"
		for _, run := range runs.WorkflowRuns {
			if run.Name == "Snap Testing" {
				totalSnapRuns++
			}
			if run.Conclusion == "failure" {
				failedSnapRuns++
			}
		}
		if failedSnapRuns == 0 {
			testIcon = "🟢"
		}

		for _, cm := range info.ChannelMap {
			var build string
			buildState, found := builtRevs[cm.Revision]
			if found && buildState == "Successfully built" {
				build = "✅"
			}
			t.AppendRow(table.Row{
				k,
				cm.Channel.Track + "/" + cm.Channel.Risk,
				cm.Version,
				cm.Channel.Architecture,
				cm.Revision,
				cm.Channel.ReleasedAt.Format(time.Stamp),
				build,
			}, table.RowConfig{AutoMerge: true})
		}
		t.AppendRow(table.Row{
			fmt.Sprintf("%s failed %d/%d", testIcon, failedSnapRuns, totalSnapRuns),
			"",
			"",
			"",
			"",
			"",
			"",
		}, table.RowConfig{AutoMerge: true})
		t.AppendSeparator()
	}

	t.Render()
}

type config struct {
	Snaps map[string]struct {
		GithubRepo string
	}
}

func loadConfig() (c *config, err error) {
	conf := flag.String("conf", configURL, "URL or local path to config file")
	flag.Parse()

	if strings.HasPrefix(*conf, "http") {
		log.Println("Fetching config file from:", *conf)

		res, err := http.Get(*conf)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&c)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Reading local config file from:", *conf)
		file, err := os.Open(*conf)
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

type builds struct {
	Entries []struct {
		StoreUploadRevision *uint `json:"store_upload_revision"`
		BuildState          string
	}
}

func queryLaunchpad(projectName string) (*builds, error) {
	log.Println("Querying launchpad for:", projectName)
	res, err := http.Get(fmt.Sprintf("https://api.launchpad.net/devel/~canonical-edgex/+snap/%s/builds?ws.size=2&direction=backwards&memo=0", projectName))
	if err != nil {
		return nil, err
	}

	var builds builds
	err = json.NewDecoder(res.Body).Decode(&builds)
	if err != nil {
		return nil, err
	}

	// log.Println("Builds:", builds)

	return &builds, nil
}

type runs struct {
	WorkflowRuns []struct {
		Name       string
		Conclusion string
		// DisplayTitle string `json:"display_title"`
		// HTMLURL      string `json:"html_url"`
	} `json:"workflow_runs"`
}

func queryGithub(project string) (*runs, error) {
	// https://api.github.com/repos/edgexfoundry/edgex-go/actions/runs?per_page=5&status=failure

	log.Println("Querying Github workflows for:", project)
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?per_page=10&event=pull_request", project))
	if err != nil {
		return nil, err
	}

	var r runs
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	// log.Println("runs:", r)

	return &r, err
}
