package client

const (
	CollectorPath    = "collectors"
	CollectorPathGet = CollectorPath + "/%s"
)

type Collector struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name"`
	TimeZone      string `json:"timeZone,omitempty"`
	CollectorType string `json:"collectorType,omitempty"`
	Ephemeral     bool   `json:"ephemeral"`
	UseExisting   bool   `json:"-"`
}

type CollectorResponse struct {
	Collector Collector `json:"collector"`
}

type CollectorsResponse struct {
	Collectors []Collector `json:"collectors"`
}
