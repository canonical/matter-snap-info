package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	config "github.com/canonical/edgex-snap-info/src/config"
	logger "github.com/canonical/edgex-snap-info/src/log"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	configURL = "https://raw.githubusercontent.com/canonical/edgex-snap-info/main/config.json"
)

func main() {
	confFile := flag.String("conf", configURL, "URL or local path to config file")
	snapName := flag.String("snap", "", "Get info for a single snap only")
	flag.Parse()

	conf, err := config.LoadConfig(*confFile)
	if err != nil {
		logger.Fatalf("Error loading config file: %s", err)
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
		// filter by snap name
		if *snapName != "" && k != *snapName {
			continue
		}

		logger.Printf(logger.Green, "⏬ %s", k)

		// snap store
		info, err := querySnapStore(k)
		if err != nil {
			logger.Fatalf("Error querying snap store: %s", err)
		}

		// launchpad
		builds, err := queryLaunchpad(k)
		if err != nil {
			logger.Fatalf("Error querying launchpad: %s", err)
		}
		revisionBuildStatus := make(map[uint]string)
		for _, v := range builds.Entries {
			// Setting a check mark only if we find the successful build result for a given revision.
			// Alternative scenarios include results that have no revision number because:
			// - build or artifact upload has failed (an actual failure)
			// - build is too old and not returned in the query
			// - build or artifact upload is pending
			if v.StoreUploadRevision != nil && v.BuildState == "Successfully built" {
				revisionBuildStatus[*v.StoreUploadRevision] = "✅"
			}
		}

		// github
		runs, err := queryGithub(v.GithubRepo)
		if err != nil {
			logger.Fatalf("Error querying launchpad: %s", err)
		}
		var totalSnapRuns, failedSnapRuns uint
		testIcon := "🔴"
		for _, run := range runs.WorkflowRuns {
			if run.Name == "Snap Testing" {
				totalSnapRuns++
			}
			if run.Conclusion == "failure" {
				failedSnapRuns++
				logger.Printf(logger.Red, "🔴 %s (%s)", run.DisplayTitle, run.HTMLURL)
			}
		}
		if totalSnapRuns == 0 { // something is not right
			testIcon = "🟠"
		} else if failedSnapRuns == 0 {
			testIcon = "🟢"
		}

		// fill the table
		for _, cm := range info.ChannelMap {
			t.AppendRow(table.Row{
				k,
				cm.Channel.Track + "/" + cm.Channel.Risk,
				cm.Version,
				cm.Channel.Architecture,
				cm.Revision,
				cm.Channel.ReleasedAt.Format(time.Stamp),
				revisionBuildStatus[cm.Revision],
			}, table.RowConfig{AutoMerge: true})
		}
		t.AppendRow(table.Row{
			fmt.Sprintf("%s failed %d/%d", testIcon, failedSnapRuns, totalSnapRuns),
			"", "", "", "", "", "",
		}, table.RowConfig{AutoMerge: true})
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
	logger.Println(logger.White, "Querying Snap Store info for:", snapName)
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

type builds struct {
	Entries []struct {
		StoreUploadRevision *uint `json:"store_upload_revision"`
		BuildState          string
	}
}

func queryLaunchpad(projectName string) (*builds, error) {
	logger.Println(logger.White, "Querying Launchpad for:", projectName)
	res, err := http.Get(fmt.Sprintf("https://api.launchpad.net/devel/~canonical-edgex/+snap/%s/builds?ws.size=10&direction=backwards&memo=0", projectName))
	if err != nil {
		return nil, err
	}

	var builds builds
	err = json.NewDecoder(res.Body).Decode(&builds)
	if err != nil {
		return nil, err
	}

	return &builds, nil
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
	logger.Println(logger.White, "Querying Github workflow runs for:", project)
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
		logger.Printf(logger.Yellow, "🟠 %s", r.Message)
	}

	return &r, err
}
