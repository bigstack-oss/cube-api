package v1

type Integration struct {
	Name                    string `json:"name"`
	IsHeaderShortcutEnabled bool   `json:"isHeaderShortcutEnabled"`
	Description             string `json:"description"`
	IsBuiltIn               bool   `json:"isBuiltIn"`
	Url                     string `json:"url"`
}