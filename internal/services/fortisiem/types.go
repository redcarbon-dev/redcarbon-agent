package fortisiem

type FortiSIEMFetchAlertsResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Pages   int                      `json:"pages"`
	QueryID string                   `json:"queryId"`
}

type Incident struct {
	IncidentID    int    `json:"incidentId"`
	IncidentTitle string `json:"incidentTitle"`
	EventSeverity int    `json:"eventSeverity"`
}
