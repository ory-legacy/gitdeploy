package config

type Config struct {
	Version    string       `yaml:"version"`
	ProcConfig ProcConfig   `yaml:"procs"`
	Godir      string       `yaml:"godir"`
	Buildpack  string       `yaml:"buildpack"`
	Addons     AddonsConfig `yaml:"addons"`
	Env        EnvConfig    `yaml:"env"`
}

type ProcConfig map[string]string

type EnvConfig map[string]string

type AddonsConfig map[string]DatabaseConfig

type DatabaseConfig struct {
	Version  string `yaml:"version"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	URL      string `yaml:"url"`
	Database string `yaml:"database"`
}
