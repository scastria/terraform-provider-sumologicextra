package client

const (
	CollectorPath        = "collectors"
	CollectorPathGet     = CollectorPath + "/%s"
	CollectorPathNameGet = CollectorPath + "/name/%s"
)

type Collector struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	TimeZone      string `json:"timeZone,omitempty"`
	CollectorType string `json:"collectorType,omitempty"`
	Ephemeral     bool   `json:"ephemeral"`
	UseExisting   bool   `json:"-"`
}

type CollectorResponse struct {
	Collector Collector `json:"collector"`
}
