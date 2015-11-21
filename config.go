package main

import (
	"io/ioutil"
	"os"

	"github.com/the-obsidian/mc/plugin"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Plugins []*plugin.Plugin

	unknownKeys []string
}

func NewConfigFromString(data string) (*Config, error) {
	c := Config{}

	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return nil, err
	}

	for _, plugin := range c.Plugins {
		err = plugin.Init()
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func NewConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewConfigFromString(string(data))
}

func (c *Config) InstallPlugins() error {
	err := os.MkdirAll("plugins", 0755)
	if err != nil {
		return err
	}

	for _, plugin := range c.Plugins {
		err = plugin.Install()
		if err != nil {
			return err
		}
	}

	return nil
}
