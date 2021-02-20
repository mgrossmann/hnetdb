package nodes

type Node struct {
	Name            string `json:"name"`
	Alias           string `json:"alias,omitempty"`
	IsGateway       bool   `json:"gateway"`
	Platform        string `json:"platform"`
	OperatingSystem string `json:"os"`
	Location        string `json:"location"`
}
