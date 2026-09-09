package auth

import (
	"net/http"
	"sync"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var oauthConfig = &oauth2.Config{
	Endpoint:    discord.Endpoint,
	Scopes:      []string{discord.ScopeIdentify},
	RedirectURL: "http://localhost:3225/callback",
	ClientID:    "1220960048704913448",
}

var (
	pendingMu sync.Mutex
	pending   = map[string]string{}
)

var httpSrvChan = make(chan bool)
var httpServer = &http.Server{
	Addr: ":3225",
}

// Disclaimer: Portions of this code have been adapted from LLM outputs.
