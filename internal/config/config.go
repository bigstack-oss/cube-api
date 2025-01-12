package config

import (
	"github.com/bigstack-oss/cube-cos-api/internal/helpers/log"
	yaml "github.com/go-micro/plugins/v5/config/encoder/yaml"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/reader"
	"go-micro.dev/v5/config/reader/json"
	"go-micro.dev/v5/config/source/file"
)

var (
	Data Payload
)

type Payload struct {
	Kind     string `json:"kind"`
	Metadata `json:"metadata"`
	Spec     `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Runtime    string `json:"runtime"`
	Dependency `json:"auth"`
	Listen     `json:"listen"`
	Log        log.Options `json:"log"`
}

type Dependency struct {
	CubeCos   string `json:"cubeCos"`
	Openstack string `json:"openstack"`
	K3s       string `json:"k3s"`
}

type Listen struct {
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
