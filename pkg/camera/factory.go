package camera

type CameraFactory interface {
	Create(c Config) (Camera, error)
}

type Config struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}
