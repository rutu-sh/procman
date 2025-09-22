package image

type Image struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ContextTempDir string `json:"context_temp_dir,omitempty"`
	ImgPath        string `json:"imgpath,omitempty"`
	Tag            string `json:"tag,omitempty"`
	Created        string `json:"created,omitempty"`
}

// image spec structs

type ImageBuildStep struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Source      string   `json:"source,omitempty"`
	Destination string   `json:"destination,omitempty"`
	Workdir     string   `json:"workdir,omitempty"`
	Command     []string `json:"command,omitempty"`
}

type ImageJob struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
}

type ImageSpec struct {
	Base  string           `json:"base"`
	Steps []ImageBuildStep `json:"steps"`
	Job   ImageJob         `json:"job"`
}
