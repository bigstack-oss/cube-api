package config

import (
	"github.com/bigstack-oss/cube-api/internal/helpers/log"
	yaml "github.com/go-micro/plugins/v5/config/encoder/yaml"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/reader"
	"go-micro.dev/v5/config/reader/json"
	"go-micro.dev/v5/config/source/file"
)

var (
	Conf Config
)

type Config struct {
	Kind     string `json:"kind"`
	Metadata `json:"metadata"`
	Spec     `json:"spec"`
}

type Metadata struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

type Spec struct {
	Runtime string `json:"runtime"`
	Auth    `json:"auth"`
	Access  `json:"access"`
	Log     log.Options `json:"log"`
}

type Auth struct {
	Openstack string `json:"openstack"`
	K3s       string `json:"k3s"`
}

type Access struct {
	Port    int `json:"port"`
	Address `json:"Address"`
}

type Address struct {
	Local     string `json:"local"`
	Advertise string `json:"advertise"`
}

func NewConfiger() (config.Config, error) {
	return config.NewConfig(
		config.WithReader(
			json.NewReader(
				reader.WithEncoder(yaml.NewEncoder()),
			),
		),
	)
}

func Load(filePath string) (config.Config, error) {
	configer, err := NewConfiger()
	if err != nil {
		return nil, err
	}

	confSrc := file.NewSource(file.WithPath(filePath))
	err = configer.Load(confSrc)
	if err != nil {
		return nil, err
	}

	return configer, nil
}
