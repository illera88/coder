package codersdk

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/xerrors"

	"github.com/coder/websocket"
)

// FIDO2SessionTokenProvider implements SessionTokenProvider with two tokens:
//   - An operate token (long-lived, used for normal API calls).
//   - A connect token (retrieved via FIDO2 device interaction, used only
//     for workspace connections via DialAgent).
//
// When the connect token has expired, the WebSocket handshake will fail
// with a 401. SetDialOption detects this on subsequent calls and
// triggers RefreshConnectToken to create a new one (requiring a YubiKey
// touch).
//
// @typescript-ignore FIDO2SessionTokenProvider
type FIDO2SessionTokenProvider struct {
	OperateToken string

	// GetConnectToken triggers FIDO2 device interaction and returns the
	// decrypted connect token. Called on first connection.
	GetConnectToken func() (string, error)

	// RefreshConnectToken creates a new connect token (using the operate
	// token), encrypts it with the FIDO2-derived key, and returns the
	// new plaintext token. Called when the previous token is expired.
	RefreshConnectToken func(ctx context.Context) (string, error)

	mu           sync.Mutex
	connectToken string
	connectErr   error
	// tokenFailed is set when a previous connection attempt failed,
	// signaling that the token should be refreshed.
	tokenFailed bool
}

// AsRequestOption returns a request option that uses the operate token
// for normal API calls (list, start, stop workspaces, etc.).
func (f *FIDO2SessionTokenProvider) AsRequestOption() RequestOption {
	return func(req *http.Request) {
		req.Header.Set(SessionTokenHeader, f.OperateToken)
	}
}

// GetSessionToken returns the operate token.
func (f *FIDO2SessionTokenProvider) GetSessionToken() string {
	return f.OperateToken
}

// SetDialOption obtains the connect token (triggering a FIDO2 assertion
// if not cached) and sets it on the websocket dial headers. This is
// called by DialAgent when establishing a workspace connection.
//
// If FIDO2 key derivation fails, no token is set on the headers. This
// causes the WebSocket handshake to fail with an auth error rather than
// silently falling back to the operate token.
func (f *FIDO2SessionTokenProvider) SetDialOption(opts *websocket.DialOptions) {
	if opts.HTTPHeader == nil {
		opts.HTTPHeader = http.Header{}
	}

	if opts.HTTPHeader.Get(SessionTokenHeader) != "" {
		return
	}

	token, err := f.ensureConnectToken()
	if err != nil {
		// Do not set any token. The connection will fail with a clear
		// auth error instead of silently using the wrong token scope.
		return
	}
	opts.HTTPHeader.Set(SessionTokenHeader, token)
}

// ConnectError returns the error from the last FIDO2 key derivation
// attempt, if any.
func (f *FIDO2SessionTokenProvider) ConnectError() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.connectErr
}

// OnDialFailure implements codersdk.DialFailureHandler. When the
// WebSocket handshake fails with 401 (unauthorized), it marks the
// cached connect token as failed so the next ensureConnectToken call
// triggers a refresh.
func (f *FIDO2SessionTokenProvider) OnDialFailure(statusCode int) {
	if statusCode == 401 {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.tokenFailed = true
	}
}

func (f *FIDO2SessionTokenProvider) ensureConnectToken() (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// If the token was previously rejected, try to refresh.
	if f.tokenFailed && f.connectToken != "" {
		f.connectToken = ""
		f.tokenFailed = false
		return f.refreshLocked()
	}

	if f.connectToken != "" {
		return f.connectToken, nil
	}

	if f.GetConnectToken == nil {
		f.connectErr = xerrors.New("no FIDO2 token function configured")
		return "", f.connectErr
	}

	token, err := f.GetConnectToken()
	if err != nil {
		f.connectErr = err
		return "", err
	}

	f.connectToken = token
	f.connectErr = nil
	return token, nil
}

// refreshLocked creates a new connect token. Caller must hold f.mu.
func (f *FIDO2SessionTokenProvider) refreshLocked() (string, error) {
	if f.RefreshConnectToken == nil {
		f.connectErr = xerrors.New("connect token expired and no refresh function configured")
		return "", f.connectErr
	}

	token, err := f.RefreshConnectToken(context.Background())
	if err != nil {
		f.connectErr = xerrors.Errorf("refresh connect token: %w", err)
		return "", f.connectErr
	}

	f.connectToken = token
	f.connectErr = nil
	return token, nil
}
