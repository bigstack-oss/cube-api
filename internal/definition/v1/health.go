package v1

type Health struct {
	Service string   `json:"service"`
	Status  string   `json:"status"`
	Modules []Module `json:"modules"`
}

type Module struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
