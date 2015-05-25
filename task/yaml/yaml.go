package yaml

type Config struct {
	ProcConfig ProcConfig   `yaml:"procs"`
	Godir      string       `yaml:"godir"`
	Buildpack  string       `yaml:"buildpack"`
	Addons     AddonsConfig `yaml:"addons"`
	Env        EnvConfig    `yaml:"env"`
}

type ProcConfig map[string]string

type EnvConfig [string]string

type AddonsConfig map[string]DatabaseConfig

type DatabaseConfig struct {
	Version  string `yaml:"version"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	URL      string `yaml:"url"`
}
