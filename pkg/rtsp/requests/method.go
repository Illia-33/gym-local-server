package requests

type Method string

const (
	MethodOptions  Method = "OPTIONS"
	MethodDescribe Method = "DESCRIBE"
	MethodSetup    Method = "SETUP"
	MethodPlay     Method = "PLAY"
	MethodPause    Method = "PAUSE"
	MethodTeardown Method = "TEARDOWN"
)
