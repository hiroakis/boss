package main

import (
	"flag"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"io/ioutil"
)

var config Config

func loadConfig(confPath string, c *Config) error {
	f, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		host string
		port int
		conf string
	)
	flag.StringVar(&host, "h", "0.0.0.0", "The listen IP. Default: 0.0.0.0")
	flag.IntVar(&port, "p", 9000, "The listen port. Default: 9000")
	flag.StringVar(&conf, "c", "./config.yml", "The configration file. Default: ./config.yml")
	flag.Parse()

	err := loadConfig(conf, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	listen := fmt.Sprintf("%s:%d", host, port)
	runServer(listen)
}
