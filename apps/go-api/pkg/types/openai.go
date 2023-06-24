package types

type GPTPrompt struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type OpenAIRequestError struct {
	Error string `json:"error"`
}
