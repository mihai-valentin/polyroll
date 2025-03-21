package elk

type elkErrorReason struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type elkError struct {
	elkErrorReason
	RootCause []elkErrorReason `json:"root_cause"`
}

type errorResponse struct {
	Error  elkError `json:"error"`
	Status int      `json:"status"`
}
