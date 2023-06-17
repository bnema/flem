package types

import "github.com/gorilla/sessions"

type App struct {
	SessionStore     sessions.Store
	PBUrl            string
	PBTradeURL       string
	PBAuthMethodsURL string
	PBAuthRefreshURL string
	OAuthRedirectURL string
	PBUserURL        string
	OpenAI_API_Key   string
	OpenAI_URL       string
	OpenAI_Model     string
}
