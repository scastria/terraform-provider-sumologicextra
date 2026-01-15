package client

const (
	CollectorPath = "collectors/%s"
)

type Collector struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
	CollectorType string `json:"collectorType,omitempty"`
}
type CollectorResponse struct {
	Collector Collector `json:"collector"`
}
