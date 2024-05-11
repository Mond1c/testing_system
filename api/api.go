package api

// TestResultResponse represents the result of the run
type TestResultResponse struct {
	Message string `json:"message"`
}

// ResponseMe represents json response for GET /api/me
type UsernameResponse struct {
	Username string `json:"username"`
}

// StartTimeResponse represents startTime of the contest and duration of the contest
type StartTimeResponse struct {
	StartTime int64 `json:"startTime"`
	Duration  int64 `json:"duration"`
}
