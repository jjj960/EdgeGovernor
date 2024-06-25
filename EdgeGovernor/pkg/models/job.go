package models

import (
	"encoding/xml"
)

type Adag struct {
	XMLName    xml.Name `xml:"adag"`
	Version    string   `xml:"version,attr"`
	ChildCount int      `xml:"childCount,attr"`
	Children   []Child  `xml:"child"`
	Jobs       []Job    `xml:"job"`
}

type Job struct {
	ID           string `xml:"id,attr"`
	Namespace    string `xml:"namespace,attr"`
	Name         string `xml:"name,attr"`
	DeployNode   string `xml:"deploynode,attr"`
	Version      string `xml:"version,attr"`
	Runtime      string `xml:"runtime,attr"`
	Image        string `xml:"image,attr"`
	Uses         []File `xml:"uses"`
	Status       bool   `json:"status"`
	WorkflowName string `json:"workflow_name"`
}

type File struct {
	File     string `xml:"file,attr"`
	Link     string `xml:"link,attr"`
	Register string `xml:"register,attr"`
	Transfer string `xml:"transfer,attr"`
	Optional string `xml:"optional,attr"`
	Type     string `xml:"type,attr"`
	Size     string `xml:"size,attr"`
}

type Child struct {
	XMLName xml.Name `xml:"child"`
	Ref     string   `xml:"ref,attr"`
	Parents []Parent `xml:"parent"`
}

type Parent struct {
	XMLName xml.Name `xml:"parent"`
	Ref     string   `xml:"ref,attr"`
}
