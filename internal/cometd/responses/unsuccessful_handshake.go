package responses

// UnsuccessfulHandshakeResponseError represents the handshake failure response.
// See: https://docs.cometd.org/current7/reference/#_unsuccessful_handshake_response
type UnsuccessfulHandshakeResponseError struct {
	Channel                  string   `json:"channel"`
	Successful               bool     `json:"successful"`
	ErrorDetails             string   `json:"error"`
	SupportedConnectionTypes []string `json:"supportedConnectionTypes,omitempty"`
	Advice                   *advice  `json:"advice,omitempty"`
	Version                  string   `json:"version,omitempty"`
	MinimumVersion           string   `json:"minimumVersion,omitempty"`
	Ext                      *ext     `json:"ext,omitempty"`
	ID                       string   `json:"id,omitempty"`
}

func (e UnsuccessfulHandshakeResponseError) Error() string {
	return e.ErrorDetails
}