package openstack

import (
	"bufio"
	"os"
	"strings"

	"github.com/bigstack-oss/cube-api/internal/config"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	log "go-micro.dev/v5/logger"
)

var (
	Opts *Options
)

type Option func(*Options)

type Helper struct {
}

func NewConf(file string) (*Options, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		log.Errorf("failed to load ops conf: %s (%s)", file, err.Error())
		return nil, err
	}
	defer openedFile.Close()
	s := bufio.NewScanner(openedFile)
	s.Split(bufio.ScanLines)

	opts := &Options{}
	for s.Scan() {
		switch {
		case strings.Contains(s.Text(), "OS_AUTH_URL"):
			words := strings.Split(s.Text(), "=")
			opts.IdentityEndpoint = words[1]
		case strings.Contains(s.Text(), "OS_AUTH_TYPE"):
			words := strings.Split(s.Text(), "=")
			opts.AuthType = words[1]
		case strings.Contains(s.Text(), "OS_USERNAME"):
			words := strings.Split(s.Text(), "=")
			opts.Username = words[1]
		case strings.Contains(s.Text(), "OS_USER_DOMAIN_NAME"):
			words := strings.Split(s.Text(), "=")
			opts.UserDomainName = words[1]
		case strings.Contains(s.Text(), "OS_PASSWORD"):
			words := strings.Split(s.Text(), "=")
			opts.Password = words[1]
		case strings.Contains(s.Text(), "OS_PROJECT_NAME"):
			words := strings.Split(s.Text(), "=")
			opts.TenantName = words[1]
		case strings.Contains(s.Text(), "OS_PROJECT_DOMAIN_NAME"):
			words := strings.Split(s.Text(), "=")
			opts.ProjectDomainName = words[1]
		}
	}

	opts.DomainName = "default"
	systemScope := os.Getenv("OS_SYSTEM_SCOPE")
	if systemScope == "all" {
		opts.Scope = &gophercloud.AuthScope{
			System: true,
		}
	}

	return opts, nil
}

func syncOptions(opts []Option) (*Options, error) {
	options, err := NewConf(config.Conf.Spec.Auth.Openstack)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		o(options)
	}

	return options, nil
}

func NewProvider(opts ...Option) (*gophercloud.ProviderClient, error) {
	syncedOpts, err := syncOptions(opts)
	if err != nil {
		return nil, err
	}

	return openstack.AuthenticatedClient(
		gophercloud.AuthOptions{
			IdentityEndpoint: syncedOpts.IdentityEndpoint,
			UserID:           syncedOpts.UserID,
			Username:         syncedOpts.Username,
			Password:         syncedOpts.Password,
			Passcode:         syncedOpts.Passcode,
			TenantID:         syncedOpts.TenantID,
			TenantName:       syncedOpts.TenantName,
			DomainID:         syncedOpts.DomainID,
			DomainName:       syncedOpts.DomainName,
			Scope:            syncedOpts.Scope,
		},
	)
}

func NewAcceleratorV1(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	opts.ApplyDefaults("accelerator")
	url, err := provider.EndpointLocator(opts)
	if err != nil {
		return nil, err
	}

	client := new(gophercloud.ServiceClient)
	client.Type = "accelerator"
	client.ProviderClient = provider
	client.Endpoint = url

	return client, nil
}
