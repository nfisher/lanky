package main

import (
	"encoding/xml"
	"io"
	"time"
)

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

type Projects struct {
	XMLName xml.Name `xml:"Projects"`
	Project []Project
}

func ReadTrayFeed(r io.Reader, p *Projects) error {
	dec := xml.NewDecoder(r)
	err := dec.Decode(p)

	return err
}
