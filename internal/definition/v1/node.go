package v1

import (
	"crypto/sha256"
	"encoding/hex"

	log "go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

const (
	Nodes = "nodes"
)

var (
	HostID        string
	Hostname      string
	Controller    string
	ControllerVip string
	ListenAddr    string
	AdvertiseAddr string
	IsHaEnabled   bool
	IsGpuEnabled  bool
)

type Node struct {
	ID           string `json:"id" yaml:"id"`
	Hostname     string `json:"hostname" yaml:"hostname"`
	Role         string `json:"role" yaml:"role"`
	Protocol     string `json:"protocol,omitempty" yaml:"protocol,omitempty" bson:"protocol,omitempty"`
	Address      string `json:"address" yaml:"address"`
	ManagementIP string `json:"managementIP" yaml:"managementIP"`
	License      `json:"license,omitempty" yaml:"license,omitempty" bson:"license,omitempty"`
	Status       string            `json:"status" yaml:"status"`
	Vcpu         ComputeStatistic  `json:"vcpu" yaml:"vcpu" bson:"vcpu"`
	Memory       SpaceStatistic    `json:"memory" yaml:"memory" bson:"memory"`
	Storage      SpaceStatistic    `json:"storage" yaml:"storage" bson:"storage"`
	Uptime       string            `json:"uptime" yaml:"uptime"`
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" bson:"labels,omitempty"`
}

func GenerateNodeHashByMacAddr() (string, error) {
	macAddr, err := GetMacAddr(NetMajorInterface)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(macAddr))
	return hex.EncodeToString(hash[:])[:8], nil
}

func GetNodesByRole(role string) ([]*Node, error) {
	svcs, err := registry.GetService(role)
	if err != nil {
		log.Errorf("failed to get service %s (%s)", role, err.Error())
		return nil, err
	}
	if len(svcs) == 0 {
		return nil, nil
	}

	nodes := []*Node{}
	for _, svc := range svcs {
		nodes = append(nodes, getNodesByService(svc, svc.Name)...)
	}

	return nodes, nil
}
