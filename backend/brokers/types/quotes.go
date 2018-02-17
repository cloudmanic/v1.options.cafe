package types

type Quote struct {
	Type             string
	Symbol           string
	Size             int
	Last             float64
	Open             float64
	High             float64
	Low              float64
	Bid              float64
	Ask              float64
	Close            float64
	PrevClose        float64
	Change           float64
	ChangePercentage float64 `json:"change_percentage"`
	Volume           int
	AverageVolume    int `json:"average_volume"`
	LastVolume       int `json:"last_volume"`
	Description      string
}
