package config

type Type string
type Transport string
type Port uint32

const (
	TransportTcp Transport = "tcp"
	TransportUdp Transport = "udp"
)

const (
	TypeOnvif Type = "onvif"
)

type Camera struct {
	Label          string    `yaml:"label"`
	Type           Type      `yaml:type`
	Ip             string    `yaml:"ip"`
	Port           Port      `yaml:"port"`
	Login          string    `yaml:"login"`
	Password       string    `yaml:"password"`
	VideoTransport Transport `yaml:"transport"`
}

type Settings struct {
}

type Config struct {
	Settings Settings `yaml:"settings"`
	Cameras  []Camera `yaml:"cameras"`
}
