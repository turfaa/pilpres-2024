package main

type SimplePredictor struct{}

func (p *SimplePredictor) PredictResult(results []CountingResult) CountingResult {
	var prediction CountingResult
	for _, r := range results {
		if r.TotalCompletedStations == 0 {
			continue
		}

		prediction.TotalStations += r.TotalStations
		prediction.TotalCompletedStations += r.TotalStations
		prediction.CandidateOne += p.predictFull(r.CandidateOne, r)
		prediction.CandidateTwo += p.predictFull(r.CandidateTwo, r)
		prediction.CandidateThree += p.predictFull(r.CandidateThree, r)

		if r.UpdatedAt > prediction.UpdatedAt {
			prediction.UpdatedAt = r.UpdatedAt
		}
	}

	prediction.Place = "Indonesia"

	return prediction
}

func (*SimplePredictor) predictFull(original float64, result CountingResult) float64 {
	prediction := original * float64(result.TotalStations)
	prediction = prediction / float64(result.TotalCompletedStations)
	return prediction
}
