package sentinel_one

type SentinelOneFetchThreatsResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type Threat struct {
	Id         string     `json:"id"`
	ThreatInfo threatInfo `json:"threatInfo"`
}

type threatInfo struct {
	Classification string `json:"classification"`
}
