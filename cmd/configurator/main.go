package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Illia-33/gym-localserver/pkg/config"

	"gopkg.in/yaml.v3"
)

var (
	netInterface = flag.String("interface", "wlp2s0", "networking interface through which cameras are accessible (use ip l for list all your interfaces)")
	outputFile   = flag.String("output", "config.yml", "config output file")
)

func main() {
	flag.Parse()
	fmt.Printf("running on interface %s\n", *netInterface)
	config, err := config.Run(*netInterface)
	if err != nil {
		fmt.Printf("generation config failure: %v\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("cannot open output file: %v\n", err)
		os.Exit(1)
	}

	enc := yaml.NewEncoder(file)
	err = enc.Encode(*config)
	if err != nil {
		fmt.Printf("encoding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("success\n")
}
