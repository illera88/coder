package codersdk_test

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/coder/coder/v2/codersdk"
	"github.com/coder/websocket"
)

func TestFIDO2SessionTokenProvider_AsRequestOption(t *testing.T) {
	t.Parallel()

	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
	}

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	p.AsRequestOption()(req)

	require.Equal(t, "operate-token", req.Header.Get("Coder-Session-Token"))
}

func TestFIDO2SessionTokenProvider_SetDialOption(t *testing.T) {
	t.Parallel()

	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			return "connect-token", nil
		},
	}

	opts := &websocket.DialOptions{}
	p.SetDialOption(opts)

	require.Equal(t, "connect-token", opts.HTTPHeader.Get("Coder-Session-Token"))
}

func TestFIDO2SessionTokenProvider_SetDialOptionCachesToken(t *testing.T) {
	t.Parallel()

	calls := 0
	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			calls++
			return "connect-token", nil
		},
	}

	opts1 := &websocket.DialOptions{}
	p.SetDialOption(opts1)
	require.Equal(t, 1, calls)

	opts2 := &websocket.DialOptions{}
	p.SetDialOption(opts2)
	require.Equal(t, 1, calls)
	require.Equal(t, "connect-token", opts2.HTTPHeader.Get("Coder-Session-Token"))
}

func TestFIDO2SessionTokenProvider_SetDialOptionError(t *testing.T) {
	t.Parallel()

	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			return "", http.ErrAbortHandler
		},
	}

	opts := &websocket.DialOptions{}
	p.SetDialOption(opts)

	require.Empty(t, opts.HTTPHeader.Get("Coder-Session-Token"))
	require.ErrorIs(t, p.ConnectError(), http.ErrAbortHandler)
}

func TestFIDO2SessionTokenProvider_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	calls := 0
	var mu sync.Mutex
	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			mu.Lock()
			calls++
			mu.Unlock()
			return "connect-token", nil
		},
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := &websocket.DialOptions{}
			p.SetDialOption(opts)
			assert.Equal(t, "connect-token", opts.HTTPHeader.Get("Coder-Session-Token"))
		}()
	}
	wg.Wait()

	mu.Lock()
	require.Equal(t, 1, calls, "GetConnectToken should be called exactly once")
	mu.Unlock()
}

func TestFIDO2SessionTokenProvider_RefreshAfterDialFailure(t *testing.T) {
	t.Parallel()

	refreshCalls := 0
	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			return "original-token", nil
		},
		RefreshConnectToken: func(_ context.Context) (string, error) {
			refreshCalls++
			return "refreshed-token", nil
		},
	}

	// First dial: uses original token.
	opts1 := &websocket.DialOptions{}
	p.SetDialOption(opts1)
	require.Equal(t, "original-token", opts1.HTTPHeader.Get("Coder-Session-Token"))
	require.Equal(t, 0, refreshCalls)

	// Simulate server 401 rejection.
	p.OnDialFailure(401)

	// Second dial: should refresh.
	opts2 := &websocket.DialOptions{}
	p.SetDialOption(opts2)
	require.Equal(t, "refreshed-token", opts2.HTTPHeader.Get("Coder-Session-Token"))
	require.Equal(t, 1, refreshCalls)
}

func TestFIDO2SessionTokenProvider_RefreshFailure(t *testing.T) {
	t.Parallel()

	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			return "original-token", nil
		},
		RefreshConnectToken: func(_ context.Context) (string, error) {
			return "", xerrors.New("YubiKey not found")
		},
	}

	opts1 := &websocket.DialOptions{}
	p.SetDialOption(opts1)
	require.Equal(t, "original-token", opts1.HTTPHeader.Get("Coder-Session-Token"))

	p.OnDialFailure(401)

	opts2 := &websocket.DialOptions{}
	p.SetDialOption(opts2)
	require.Empty(t, opts2.HTTPHeader.Get("Coder-Session-Token"))
	require.Error(t, p.ConnectError())
	require.Contains(t, p.ConnectError().Error(), "YubiKey not found")
}

func TestFIDO2SessionTokenProvider_Non401DoesNotTriggerRefresh(t *testing.T) {
	t.Parallel()

	refreshCalls := 0
	p := &codersdk.FIDO2SessionTokenProvider{
		OperateToken: "operate-token",
		GetConnectToken: func() (string, error) {
			return "valid-token", nil
		},
		RefreshConnectToken: func(_ context.Context) (string, error) {
			refreshCalls++
			return "should-not-be-called", nil
		},
	}

	opts1 := &websocket.DialOptions{}
	p.SetDialOption(opts1)
	require.Equal(t, "valid-token", opts1.HTTPHeader.Get("Coder-Session-Token"))

	p.OnDialFailure(500)

	opts2 := &websocket.DialOptions{}
	p.SetDialOption(opts2)
	require.Equal(t, "valid-token", opts2.HTTPHeader.Get("Coder-Session-Token"))
	require.Equal(t, 0, refreshCalls)
}
