package plugin

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hashicorp/go-getter"
	"github.com/the-obsidian/mc/util"
)

// Plugin represents an individual Minecraft plugin dependency
type Plugin struct {
	Name               string
	URI                string
	Sha                string
	Processors         []Processor              `yaml:"-"`
	OriginalProcessors []map[string]interface{} `yaml:"processors,flow"`
}

// Init processes the config into the required structure
func (p *Plugin) Init() error {
	for _, pr := range p.OriginalProcessors {
		processor, err := NewProcessorFromConfig(p, pr)
		if err != nil {
			return err
		}
		p.Processors = append(p.Processors, processor)
	}

	return nil
}

// Install downlods the plugin if it doesn't exist, processes it, and
// installs it to the plugins dir
func (p *Plugin) Install() error {
	dest := path.Join(".", "plugins", p.Name+".jar")
	tmpDir := path.Join(".", ".tmp", "plugins", p.Name)
	tmpDownload := path.Join(tmpDir, "download")
	hash := sha1.New()

	checksumValue, err := decodeChecksum(p.Sha)
	if err != nil {
		return err
	}

	if exists, err := util.FileExists(dest); err != nil {
		return err
	} else if exists {
		err = checksum(dest, hash, checksumValue)
		if err == nil {
			return nil
		}
	}

	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return err
	}

	err = getter.GetFile(tmpDownload, p.URI)
	if err != nil {
		return err
	}

	for _, pr := range p.Processors {
		err = pr.Process(p)
		if err != nil {
			return fmt.Errorf("failed to run processor: %s", err)
		}
	}

	err = checksum(tmpDownload, hash, checksumValue)
	if err != nil {
		return fmt.Errorf("failed to checksum download: %v", err)
	}

	err = os.Rename(tmpDownload, dest)
	if err != nil {
		return err
	}

	return os.RemoveAll(tmpDir)
}

// Processor represents a plugin post-processor
type Processor interface {
	Process(p *Plugin) error
}

// NewProcessorFromConfig builds a Processor from a config
func NewProcessorFromConfig(p *Plugin, config map[string]interface{}) (Processor, error) {
	switch config["type"] {
	case "unzip":
		pr := &UnzipProcessor{
			Type: "unzip",
			File: config["file"].(string),
		}
		return pr, nil
	default:
		return nil, fmt.Errorf("invalid processor type: %v", config["type"])
	}
}

// UnzipProcessor unzips a file to extract a plugin
type UnzipProcessor struct {
	Type string
	File string
}

// Process performs the unzip
func (pr *UnzipProcessor) Process(p *Plugin) error {
	tmpDir := path.Join(".", ".tmp", "plugins", p.Name)
	tmpDownload := path.Join(tmpDir, "download")
	tmpDest := path.Join(tmpDir, "download-extract")

	r, err := zip.OpenReader(tmpDownload)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != pr.File {
			continue
		}

		fmt.Printf("Extracting file %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		f, err := os.OpenFile(tmpDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}

		return os.Rename(tmpDest, tmpDownload)
	}

	return fmt.Errorf("failed to extract file ")
}
