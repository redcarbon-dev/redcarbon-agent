package qradar

type Offense struct {
	Description                string  `json:"description"`
	StartTime                  float64 `json:"start_time"`
	LastUpdatedTime            float64 `json:"last_updated_time"`
	Magnitude                  int     `json:"magnitude"`
	ID                         int     `json:"id"`
	LocalDestinationAddressIds []int   `json:"local_destination_address_ids"`
	SourceAddressIds           []int   `json:"source_address_ids"`
	OffenseType                int     `json:"offense_type"`
}
