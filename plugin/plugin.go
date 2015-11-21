package plugin

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter"
	"github.com/the-obsidian/mc/util"
)

type Plugin struct {
	Name       string
	URI        string
	Sha        string
	Processors []Processor `yaml:"-"`

	originalProcessors []map[string]interface{} `yaml:"processors"`
}

func (p *Plugin) Init() error {
	for _, pr := range p.originalProcessors {
		processor, err := NewProcessorFromConfig(p, pr)
		if err != nil {
			return err
		}
		p.Processors = append(p.Processors, processor)
	}

	return nil
}

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

	err = checksum(tmpDownload, hash, checksumValue)
	if err != nil {
		return fmt.Errorf("failed to checksum download: %v", err)
	}

	err = os.Rename(tmpDownload, dest)
	if err != nil {
		return err
	}

	return nil
}

type Processor interface {
	Process()
}

func NewProcessorFromConfig(p *Plugin, config map[string]interface{}) (Processor, error) {
	switch config["type"] {
	case "unzip":
		pr := &UnzipProcessor{
			Type:  "unzip",
			Files: config["files"].([]string),
		}
		return pr, nil
	default:
		return nil, fmt.Errorf("invalid processor type: %v", config["type"])
	}

	return nil, nil
}

type UnzipProcessor struct {
	Type  string
	Files []string
}

func (pr *UnzipProcessor) Process() {}
