// Handles Microsoft Entra ID (Azure AD) auth
package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// EntraConfig holds Microsoft Entra ID configuration
type EntraConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Authenticator handles OIDC authentication flow with Entra ID
type Authenticator struct {
	config   *EntraConfig
	provider *oidc.Provider
	oauth2   *oauth2.Config
	verifier *oidc.IDTokenVerifier

	// State management for CSRF protection
	states     map[string]*authState
	statesLock sync.RWMutex
}

type authState struct {
	state     string
	nonce     string
	pubKey    string // User's SSH public key
	createdAt time.Time
}

// UserInfo contains verified user information from Entra ID
type UserInfo struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Username      string `json:"preferred_username"`
}

// NewAuthenticator creates a new Entra ID authenticator
func NewAuthenticator(ctx context.Context, cfg *EntraConfig) (*Authenticator, error) {
	issuerURL := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", cfg.TenantID)

	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &Authenticator{
		config:   cfg,
		provider: provider,
		oauth2:   oauth2Config,
		verifier: verifier,
		states:   make(map[string]*authState),
	}, nil
}

// StartAuth initiates the authentication flow
// Returns the authorization URL to redirect the user to
func (a *Authenticator) StartAuth(pubKey string) (string, error) {
	state, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	nonce, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Store state for verification
	a.statesLock.Lock()
	a.states[state] = &authState{
		state:     state,
		nonce:     nonce,
		pubKey:    pubKey,
		createdAt: time.Now(),
	}
	a.statesLock.Unlock()

	// Clean up old states
	go a.cleanupStates()

	url := a.oauth2.AuthCodeURL(state, oidc.Nonce(nonce))
	return url, nil
}

// HandleCallback processes the OIDC callback and returns user info
func (a *Authenticator) HandleCallback(ctx context.Context, r *http.Request) (*UserInfo, string, error) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" || code == "" {
		return nil, "", fmt.Errorf("missing state or code")
	}

	// Verify state
	a.statesLock.RLock()
	authState, exists := a.states[state]
	a.statesLock.RUnlock()

	if !exists {
		return nil, "", fmt.Errorf("invalid state - possible CSRF attack")
	}

	// Remove used state
	a.statesLock.Lock()
	delete(a.states, state)
	a.statesLock.Unlock()

	// Exchange code for token
	token, err := a.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", fmt.Errorf("no id_token in response")
	}

	// Verify ID token
	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to verify ID token: %w", err)
	}

	// Verify nonce
	var claims struct {
		Nonce string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, "", fmt.Errorf("failed to parse claims: %w", err)
	}
	if claims.Nonce != authState.nonce {
		return nil, "", fmt.Errorf("invalid nonce - possible replay attack")
	}

	// Extract user info
	var userInfo UserInfo
	if err := idToken.Claims(&userInfo); err != nil {
		return nil, "", fmt.Errorf("failed to extract user info: %w", err)
	}

	return &userInfo, authState.pubKey, nil
}

// cleanupStates removes expired auth states (older than 10 minutes)
func (a *Authenticator) cleanupStates() {
	a.statesLock.Lock()
	defer a.statesLock.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	for key, state := range a.states {
		if state.createdAt.Before(cutoff) {
			delete(a.states, key)
		}
	}
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
