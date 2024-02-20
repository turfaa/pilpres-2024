package main

type KawalPemiluResponse struct {
	Result KawalPemiluResponseData `json:"result"`
}

type KawalPemiluResponseData struct {
	Aggregated map[string][]CountingResult `json:"aggregated"`
}

type CountingResult struct {
	CandidateOne           float64 `json:"pas1"`
	CandidateTwo           float64 `json:"pas2"`
	CandidateThree         float64 `json:"pas3"`
	TotalCompletedStations int64   `json:"totalCompletedTps"`
	TotalStations          int64   `json:"totalTps"`
	Place                  string  `json:"name"`
	UpdatedAt              int64   `json:"updateTs"`
}

func (r CountingResult) TotalVotes() float64 {
	return r.CandidateOne + r.CandidateTwo + r.CandidateThree
}

func (r CountingResult) CandidateOnePercent() float64 {
	return r.CandidateOne / r.TotalVotes() * 100
}

func (r CountingResult) CandidateTwoPercent() float64 {
	return r.CandidateTwo / r.TotalVotes() * 100
}

func (r CountingResult) CandidateThreePercent() float64 {
	return float64(r.CandidateThree) / float64(r.TotalVotes()) * 100
}
