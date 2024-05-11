package api

// TestResultResponse represents the result of the run
type TestResultResponse struct {
	Message string `json:"message"`
}

// ResponseProblems represents json response for GET /api/problems
type ResponseProblems struct {
	Problems []string `json:"problems"`
}

// ResponseMe represents json response for GET /api/me
type ResponseMe struct {
	Username string `json:"username"`
}

// StartTimeResponse represents startTime of the contest and duration of the contest
type StartTimeResponse struct {
	StartTime int64 `json:"startTime"`
	Duration  int64 `json:"duration"`
}
