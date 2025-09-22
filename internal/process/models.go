package process

import (
	"github.com/rutu-sh/procman/internal/image"
)

type ProcessEnv map[string]string

type Process struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Pid        uint           `json:"pid"`
	ContextDir string         `json:"context_dir"`
	Image      image.Image    `json:"image"`
	Job        image.ImageJob `json:"job"`
	Env        ProcessEnv     `json:"env"`
	Network    ProcessNetwork `json:"network,omitempty"`
}

type ProcessCreateImage struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type ProcessCreate struct {
	Name  string 				`json:"name"`
	Image ProcessCreateImage	`json:"image"`
	Env   ProcessEnv 			`json:"env"`
}

type PortMapping struct {
	HostPort uint `json:"host_port"`
	ProcPort uint `json:"proc_port"`
}

type ProcessNetwork struct {
	Ports []PortMapping `json:"ports,omitempty"`
}
