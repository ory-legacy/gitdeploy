package config

type Config struct {
	Version    string            `yaml:"version"`
	ProcConfig map[string]string `yaml:"procs,flow"`
	Godir      string            `yaml:"godir"`
	Buildpack  string            `yaml:"buildpack"`
	Addons     []DatabaseConfig  `yaml:"addons,flow"`
	Env        map[string]string `yaml:"env,flow"`
}

type DatabaseConfig struct {
	Version  string `yaml:"version"`
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	URL      string `yaml:"url"`
	Database string `yaml:"database"`
}
