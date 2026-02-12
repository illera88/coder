package fido2

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

const helperBinary = "coder-fido2"

// ErrTouchTimeout indicates the user did not touch the security key in
// time. Callers should prompt the user to try again.
var ErrTouchTimeout = xerrors.New("security key touch timed out")

// ErrPinRequired indicates the security key requires a PIN that was
// not provided or was incorrect. Callers should prompt for the PIN.
var ErrPinRequired = xerrors.New("security key requires a PIN")

// runHelper executes the coder-fido2 helper with the given subcommand,
// piping the optional PIN via stdin (never as a CLI argument).
func runHelper(subcmd string, configDir string, pin string) (stdout string, err error) {
	cmd := exec.Command(helperBinary, subcmd, "--config-dir", configDir)
	// Send PIN on stdin. Empty string = no PIN (touch only).
	cmd.Stdin = strings.NewReader(pin + "\n")
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			switch exitErr.ExitCode() {
			case 2:
				return "", ErrTouchTimeout
			case 3:
				return "", ErrPinRequired
			}
		}
		return "", xerrors.Errorf("coder-fido2 %s: %s: %w", subcmd, strings.TrimSpace(errBuf.String()), err)
	}
	return outBuf.String(), nil
}

// RunRegister shells out to the coder-fido2 helper to register a FIDO2
// credential. PIN is optional (empty string for touch-only).
func RunRegister(configDir string, pin string) error {
	_, err := runHelper("register", configDir, pin)
	return err
}

// RunDeriveKey shells out to the coder-fido2 helper to perform a FIDO2
// assertion and returns the 32-byte derived key. The user must touch
// their security key. PIN is optional (empty string for touch-only).
func RunDeriveKey(configDir string, pin string) ([]byte, error) {
	out, err := runHelper("derive-key", configDir, pin)
	if err != nil {
		return nil, err
	}
	return decodeHexKey(out)
}

// IsHelperInstalled returns true if the coder-fido2 helper binary is
// found on PATH.
func IsHelperInstalled() bool {
	_, err := exec.LookPath(helperBinary)
	return err == nil
}

func decodeHexKey(out string) ([]byte, error) {
	key, err := hex.DecodeString(strings.TrimSpace(out))
	if err != nil {
		return nil, xerrors.Errorf("decode derived key: %w", err)
	}
	if len(key) != 32 {
		return nil, xerrors.Errorf("derived key has unexpected length %d", len(key))
	}
	return key, nil
}
