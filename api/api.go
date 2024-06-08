package api

// TestResultResponse represents the result of the run
type TestResultResponse struct {
	Message string `json:"message"`
}

// UsernameResponse represents json response for GET /api/me
type UsernameResponse struct {
	Username string `json:"username"`
}

// StartTimeResponse represents startTime of the contest and duration of the contest
type StartTimeResponse struct {
	StartTime int64 `json:"startTime"`
	Duration  int64 `json:"duration"`
}

type RunInfoResponse struct {
	Username string `json:"username"`
	RunID    int    `json:"run_id"`
	Problem  string `json:"problem"`
	Result   string `json:"result"`
	Time     int64  `json:"time"`
	Language string `json:"language"`
}
