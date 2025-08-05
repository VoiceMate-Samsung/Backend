package models

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Response string `json:"response,omitempty"`
	Screen   string `json:"screen,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Existing Gemini models
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate  `json:"candidates"`
	Error      *GeminiError `json:"error,omitempty"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type GeminiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
