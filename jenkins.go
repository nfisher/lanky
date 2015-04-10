package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

const (
	orderByDate   = "date"
	orderByStatus = "status"
)

func NewJenkins(config *Config) (client *JenkinsClient) {
	if config.Jenkins == nil || config.Jenkins.BaseUrl == "" || config.Jenkins.TrayFeed == "" {
		return nil
	}

	wc := &http.Client{
		Timeout: config.ClientTimeout(),
	}

	return &JenkinsClient{
		config,
		wc,
	}
}

type WebClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
}

type JenkinsClient struct {
	*Config
	WebClient
}

func (j *JenkinsClient) TrayFeed(p *Projects, by string) (err error) {
	trayFeedUrl := j.Config.TrayFeedUrl()

	resp, err := j.WebClient.Get(trayFeedUrl)
	if err != nil {
		msg := fmt.Sprintf("%v with a timeout of %v", err.Error(), j.Config.ClientTimeout())
		return errors.New(msg)
	}

	defer resp.Body.Close()
	err = ReadTrayFeed(resp.Body, p)
	if err != nil {
		return errors.New(err.Error() + " from " + trayFeedUrl)
	}

	switch by {
	case orderByDate:
		sort.Sort(ByStatus{p})
		sort.Stable(sort.Reverse(ByDate{p}))
		p.Order = orderByDate
		break
	default:
		sort.Sort(sort.Reverse(ByDate{p}))
		sort.Stable(ByStatus{p})
		p.Order = orderByStatus
		break
	}

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
	Order   string
}

func (p *Projects) ByDate() bool  { return p.Order == orderByDate }
func (p *Projects) Len() int      { return len(p.Project) }
func (p *Projects) Swap(i, j int) { p.Project[i], p.Project[j] = p.Project[j], p.Project[i] }

type ByStatus struct{ *Projects }

func (p ByStatus) Less(i, j int) bool {
	return p.Project[i].LastBuildStatus < p.Project[j].LastBuildStatus
}

type ByDate struct{ *Projects }

func (p ByDate) Less(i, j int) bool {
	return p.Project[i].LastBuildTime.Unix() < p.Project[j].LastBuildTime.Unix()
}

func ReadTrayFeed(r io.Reader, p *Projects) error {
	dec := xml.NewDecoder(r)
	err := dec.Decode(p)

	return err
}
