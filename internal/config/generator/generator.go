package generator

import (
	"fmt"
	"strconv"
	"strings"

	cfg "github.com/Illia-33/gym-localserver/pkg/config"
	"github.com/beevik/etree"
	discovery "github.com/use-go/onvif/ws-discovery"
)

type Host struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func Run(interfaceName string) (cfg.Config, error) {
	devices, err := discovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		return cfg.Config{}, err
	}

	var config cfg.Config
	foundXaddrs := map[string]bool{}
	config.Cameras = make([]cfg.Camera, 0, len(devices))

	for i, deviceInfo := range devices {
		var camera cfg.Camera

		doc := etree.NewDocument()
		if err := doc.ReadFromString(deviceInfo); err != nil {
			continue
		}

		endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
		alreadyFound := false
		for _, xaddr := range endpoints {
			xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
			if foundXaddrs[xaddr] {
				alreadyFound = true
				break
			}

			splitXaddr := strings.Split(xaddr, ":")
			port, err := strconv.ParseUint(splitXaddr[1], 10, 64)
			if err != nil {
				continue
			}

			foundXaddrs[xaddr] = true
			camera.Ip = splitXaddr[0]
			camera.Port = cfg.Port(port)
		}

		if alreadyFound {
			continue
		}

		camera.Label = fmt.Sprintf("camera_%d", i)
		camera.Type = cfg.TypeOnvif
		camera.Login = "<camera_login>"
		camera.Password = "<camera_password>"

		config.Cameras = append(config.Cameras, camera)
	}

	config.Settings.AuthKey = "<auth_key>"

	return config, nil
}
