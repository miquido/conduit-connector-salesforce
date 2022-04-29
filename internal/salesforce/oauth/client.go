package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/miquido/conduit-connector-salesforce/internal/salesforce/oauth/response"
	"github.com/miquido/conduit-connector-salesforce/internal/utils"
)

type Environment = string

const (
	EnvironmentSandbox Environment = "sandbox"

	grantType    = "password"
	loginURI     = "https://login.salesforce.com/services/oauth2/token"
	testLoginURI = "https://test.salesforce.com/services/oauth2/token"
)

func NewClient(
	environment Environment,
	clientID string,
	clientSecret string,
	username string,
	password string,
	securityToken string,
) *Client {
	return &Client{
		environment:   environment,
		clientID:      clientID,
		clientSecret:  clientSecret,
		username:      username,
		password:      password,
		securityToken: securityToken,
	}
}

type Client struct {
	environment   Environment
	clientID      string
	clientSecret  string
	username      string
	password      string
	securityToken string
}

func (a *Client) Authenticate(ctx context.Context) (response.TokenResponse, error) {
	payload := url.Values{
		"grant_type":    {grantType},
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
		"username":      {a.username},
		"password":      {fmt.Sprintf("%v%v", a.password, a.securityToken)},
	}

	// Build Uri
	uri := loginURI
	if EnvironmentSandbox == a.environment {
		uri = testLoginURI
	}

	// Build Body
	body := strings.NewReader(payload.Encode())

	// Build Request
	req, err := http.NewRequestWithContext(ctx, "POST", uri, body)
	if err != nil {
		return response.TokenResponse{}, fmt.Errorf("failed to prepare authentication request: %w", err)
	}

	// Add Headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip;q=1.0, *;q=0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "ConduitIO/Salesforce-v0.1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response.TokenResponse{}, fmt.Errorf("failed to send authentication request: %w", err)
	}

	respBytes, err := utils.DecodeHTTPResponse(resp)
	if err != nil {
		return response.TokenResponse{}, fmt.Errorf("could not read response data: %w", err)
	}

	resp.Body.Close()

	fmt.Println("OAuth response", string(respBytes))

	// Attempt to parse successful response
	var token response.TokenResponse
	if err := json.Unmarshal(respBytes, &token); err == nil {
		return token, nil
	}

	// Attempt to parse response as a force.com api error
	authFailureResponse := response.FailureResponseError{}

	if err := json.Unmarshal(respBytes, &authFailureResponse); err != nil {
		return response.TokenResponse{}, fmt.Errorf("unable to process authentication response: %w", err)
	}

	return response.TokenResponse{}, authFailureResponse
}