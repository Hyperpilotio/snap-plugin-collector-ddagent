package dogstatsd

type Metric struct {
	SourceTypeName string      `json:"source_type_name"`
	MetricName     string      `json:"metric"`
	Points         [][]float64 `json:"points"`
	Type           string      `json:"type"`
	Host           string      `json:"host"`
}

type Metrics struct {
	Series []Metric `json:"series"`
}
