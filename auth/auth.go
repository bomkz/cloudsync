package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func AuthenticateUser() error {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:3225/login").Start(); err != nil {
		return err
	}

	<-httpSrvChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)

}

func stopServer() {
	close(httpSrvChan)
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Step 1: redirect user to the provider's consent page
func handleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := randomState()
	if err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	verifier := oauth2.GenerateVerifier()

	pendingMu.Lock()
	pending[state] = verifier
	pendingMu.Unlock()

	url := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Step 2: handle the callback, exchange code for a token
func handleCallback(w http.ResponseWriter, r *http.Request) {
	if providerError := r.URL.Query().Get("error"); providerError != "" {
		http.Error(w, "authorization was denied", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")

	pendingMu.Lock()
	verifier, ok := pending[state]
	if ok {
		delete(pending, state)
	}
	pendingMu.Unlock()

	if !ok {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifier),
	)
	if err != nil {
		log.Printf("discord OAuth token exchange failed: %v", err)
		http.Error(w, "exchange failed", http.StatusInternalServerError)
		return
	}

	err = storeAuthToken(token)
	if err != nil {
		log.Printf("failed to store auth token: %v", err)
		http.Error(w, "failed to store auth token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html>
<body>
<p>Authentication complete. You can close this tab.</p>
<script>
window.open('', '_self');
window.close();
</script>
</body>
</html>`))
	stopServer()
}

func storeAuthToken(token *oauth2.Token) (err error) {
	svc_name := "cloudsave"
	svc_user := "default"

	err = keyring.Set(svc_name, svc_user, token.AccessToken)

	return
}

func GetAuthToken() (string, error) {
	svc_name := "cloudsave"
	svc_user := "default"

	token, err := keyring.Get(svc_name, svc_user)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Disclaimer: Portions of this code have been adapted from LLM outputs.
