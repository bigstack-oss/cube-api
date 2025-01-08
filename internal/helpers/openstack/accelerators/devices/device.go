package devices

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type Device struct {
	UUID            string `json:"uuid"`
	Hostname        string `json:"hostname"`
	Type            string `json:"type"`
	Vendor          string `json:"vendor"`
	Model           string `json:"model"`
	StdBoardInfo    string `json:"std_board_info"`
	VendorBoardInfo string `json:"vendor_board_info"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ListOptsBuilder interface {
	ToDeviceListQuery() (string, error)
}

type ListOpts struct {
	Type     string            `q:"type"`
	Vendor   string            `q:"vendor"`
	Hostname string            `q:"hostname"`
	Filters  map[string]string `q:"-"`
}

func (opts ListOpts) ToDeviceListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func listDetailURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("devices")
}

type DevicePage struct {
	pagination.LinkedPageBase
}

func pageLinker(r pagination.PageResult) pagination.Page {
	return DevicePage{pagination.LinkedPageBase{PageResult: r}}
}

func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) ([]Device, error) {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToDeviceListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	devices := &[]Device{}
	pager := pagination.NewPager(client, url, pageLinker)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		err := ExtractDevicesInto(page, devices)
		return false, err
	})
	if err != nil {
		return nil, err
	}

	return *devices, nil
}

func (d DevicePage) IsEmpty() (bool, error) {
	if d.StatusCode == 204 {
		return true, nil
	}

	devices, err := ExtractDevices(d)
	return len(devices) == 0, err
}

func ExtractDevices(r pagination.Page) ([]Device, error) {
	d := []Device{}
	err := ExtractDevicesInto(r, &d)
	return d, err
}

func (r DevicePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"servers_links"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(s.Links)
}

func ExtractDevicesInto(r pagination.Page, v interface{}) error {
	return r.(DevicePage).Result.ExtractIntoSlicePtr(v, "devices")
}
