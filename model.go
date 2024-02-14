package main

type KawalPemiluRequest struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type KawalPemiluResponse struct {
	Result KawalPemiluResponseData `json:"result"`
}

type KawalPemiluResponseData struct {
	Aggregated map[string][]CountingResult `json:"aggregated"`
}

type CountingResult struct {
	CandidateOne           int64  `json:"pas1"`
	CandidateTwo           int64  `json:"pas2"`
	CandidateThree         int64  `json:"pas3"`
	TotalCompletedStations int64  `json:"totalCompletedTps"`
	TotalStations          int64  `json:"totalTps"`
	Place                  string `json:"name"`
	UpdatedAt              int64  `json:"updateTs"`
}

func (r CountingResult) TotalVotes() int64 {
	return r.CandidateOne + r.CandidateTwo + r.CandidateThree
}

func (r CountingResult) CandidateOnePercent() float64 {
	return float64(r.CandidateOne) / float64(r.TotalVotes()) * 100
}

func (r CountingResult) CandidateTwoPercent() float64 {
	return float64(r.CandidateTwo) / float64(r.TotalVotes()) * 100
}

func (r CountingResult) CandidateThreePercent() float64 {
	return float64(r.CandidateThree) / float64(r.TotalVotes()) * 100
}
