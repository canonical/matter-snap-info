package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/canonical/edgex-snap-info/config"
	"github.com/canonical/edgex-snap-info/logger"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	configURL = "https://raw.githubusercontent.com/canonical/edgex-snap-info/main/config.json"
)

func main() {
	confFile := flag.String("conf", configURL, "URL or local path to config file")
	snapName := flag.String("snap", "", "Get info for a single snap only")
	flag.Parse()

	conf, err := config.Load(*confFile)
	if err != nil {
		logger.Fatalf("Error loading config file: %s", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Name", "Channel", "Version", "Arch", "Rev", "Date", "Build", "Test"})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
	})

	for k, v := range conf.Snaps {
		// filter by snap name
		if *snapName != "" && k != *snapName {
			continue
		}

		logger.Infof("⏬ %s", k)

		// snap store
		info, err := querySnapStore(k)
		if err != nil {
			logger.Fatalf("Error querying snap store: %s", err)
		}

		// github
		runs, err := queryGithub(v.GithubRepo)
		if err != nil {
			logger.Fatalf("Error querying launchpad: %s", err)
		}
		var build, test snapStatistcs
		build.runName = "Snap Builder"
		test.runName = "Snap Testing"
		handleRun(runs, &build)
		handleRun(runs, &test)

		// fill the table
		for _, cm := range info.ChannelMap {
			t.AppendRow(table.Row{
				k,
				cm.Channel.Track + "/" + cm.Channel.Risk,
				cm.Version,
				cm.Channel.Architecture,
				cm.Revision,
				cm.Channel.ReleasedAt.Format(time.Stamp),
				fmt.Sprintf("%d/%d", build.success, build.total),
				fmt.Sprintf("%d/%d", test.success, test.total),
			}, table.RowConfig{AutoMerge: true})
		}

		t.AppendSeparator()
	}

	t.Render()
}

type snapStatistcs struct {
	total, success uint
	runName        string
}

func handleRun(run *runs, s *snapStatistcs) {
	for _, r := range run.WorkflowRuns {
		if r.Name != s.runName {
			continue
		}
		s.total++
		if r.Conclusion == "success" {
			*&s.success++
			// logger.Successf("🟢 %s (%s)", r.DisplayTitle, r.HTMLURL)
		} else {
			logger.Errorf("🔴 %s (%s)", r.DisplayTitle, r.HTMLURL)
		}
	}
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
	logger.Infoln("Querying Snap Store info for:", snapName)
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

	return &info, nil
}

type runs struct {
	WorkflowRuns []struct {
		Name         string
		Conclusion   string
		DisplayTitle string `json:"display_title"`
		HTMLURL      string `json:"html_url"`
	} `json:"workflow_runs"`
	Message string
}

func queryGithub(project string) (*runs, error) {
	logger.Infoln("Querying Github workflow runs for:", project)
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?per_page=10&event=pull_request", project))
	if err != nil {
		return nil, err
	}

	var r runs
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Message != "" {
		logger.Warnf("🟠 %s", r.Message)
	}

	return &r, err
}
