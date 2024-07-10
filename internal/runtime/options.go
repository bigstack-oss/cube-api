package runtime

import "github.com/bigstack-oss/cube-api/internal/helpers/log"

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
