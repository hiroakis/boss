package main

type Config struct {
	RepoConfigs map[string]RepoConfig `yaml:"repos"`
}

type RepoConfig struct {
	Token   string   `yaml:"token"`
	Members []string `yaml:"members"`
	Labels  []string `yaml:"labels"`
}
