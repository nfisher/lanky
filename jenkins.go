package main

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"sort"
	"time"
)

type JenkinsClient struct {
	*Config
}

func (j *JenkinsClient) TrayFeed(p *Projects) (err error) {
	trayFeedUrl := j.Config.TrayFeedUrl()

	client := http.Client{
		Timeout: j.Config.ClientTimeout(),
	}

	resp, err := client.Get(trayFeedUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = ReadTrayFeed(resp.Body, p)
	if err != nil {
		return errors.New(err.Error() + " from " + trayFeedUrl)
	}

	sort.Sort(ByStatus{p})

	return nil
}

func (j *JenkinsClient) TriggerBuild() (err error) {
	return nil
}

type Project struct {
	WebUrl          string    `xml:"webUrl,attr"`
	Name            string    `xml:"name,attr"`
	LastBuildLabel  string    `xml:"lastBuildLabel,attr"`
	LastBuildTime   time.Time `xml:"lastBuildTime,attr"`
	LastBuildStatus string    `xml:"lastBuildStatus,attr"`
	Activity        string    `xml:"activity,attr"`
}

func (p *Project) BuildTime() string {
	return p.LastBuildTime.Format("2006-01-02 15:04")
}

func (p *Project) ConsoleUrl() string {
	return p.WebUrl + p.LastBuildLabel + "/console"
}

type Projects struct {
	XMLName xml.Name `xml:"Projects"`
	Project []Project
}

func (p *Projects) Len() int      { return len(p.Project) }
func (p *Projects) Swap(i, j int) { p.Project[i], p.Project[j] = p.Project[j], p.Project[i] }

type ByStatus struct{ *Projects }

func (p ByStatus) Less(i, j int) bool {
	return p.Project[i].LastBuildStatus < p.Project[j].LastBuildStatus
}

func ReadTrayFeed(r io.Reader, p *Projects) error {
	dec := xml.NewDecoder(r)
	err := dec.Decode(p)

	return err
}
