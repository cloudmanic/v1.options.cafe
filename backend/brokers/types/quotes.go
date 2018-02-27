package types

type Quote struct {
	Type             string  `json:"type"`
	Symbol           string  `json:"symbol"`
	Size             int     `json:"size"`
	Last             float64 `json:"last"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Bid              float64 `json:"bid"`
	Ask              float64 `json:"ask"`
	Close            float64 `json:"close"`
	PrevClose        float64 `json:"prev_close"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"change_percentage"`
	Volume           int     `json:"volume"`
	AverageVolume    int     `json:"average_volume"`
	LastVolume       int     `json:"last_volume"`
	Description      string  `json:"description"`
}
