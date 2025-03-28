package config

type Type string
type Transport string
type Port uint32

const (
	TypeOnvif Type = "onvif"
)

type Camera struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Type        Type   `yaml:"type"`
	Ip          string `yaml:"ip"`
	Port        Port   `yaml:"port"`
	Login       string `yaml:"login"`
	Password    string `yaml:"password"`
}

type Settings struct {
	AuthKey string `yaml:"auth_key"`
}

type Config struct {
	Settings Settings `yaml:"settings"`
	Cameras  []Camera `yaml:"cameras"`
}
