package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Title          string            `yaml:"title"`
	Subtitle       string            `yaml:"subtitle"`
	Email          string            `yaml:"email"`
	Description    string            `yaml:"description"`
	Baseurl        string            `yaml:"baseurl"`
	Url            string            `yaml:"url"`
	GithubUsername string            `yaml:"github_username"`
	Custom         map[string]string `yaml:"custom"`
}

func (c *Config) SetFromFile(p string) error {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s not found", p))
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("Error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return fmt.Errorf("Error setting config from YAML: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	return nil
}
