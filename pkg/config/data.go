package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type TechStackItem struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ImagePath   string `yaml:"imagePath"`
	Url         string `yaml:"url"`
	Section     string `yaml:"section"`
}

type SocialItem struct {
	Name      string `yaml:"name"`
	Url       string `yaml:"url"`
	ImagePath string `yaml:"imagePath"`
}

type MetaData struct {
	Name    string `yaml:"name"`
	Role    string `yaml:"role"`
	Image   string `yaml:"image"`
	AboutMe string `yaml:"aboutMe"`
}

type Data struct {
	MetaData  MetaData
	TechStack []TechStackItem `yaml:"techStack"`
	Socials   []SocialItem    `yaml:"socials"`
}

func LoadData(filename string) (Data, error) {
	var config Data
	b, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
