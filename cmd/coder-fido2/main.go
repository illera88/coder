// Command coder-fido2 is a helper binary for FIDO2 device interaction.
// It is built separately with CGo + libfido2 so the main coder binary
// stays pure Go. This is a standalone Go module — it does not import
// the main coder module.
//
// Usage:
//
//	coder-fido2 register --config-dir <dir>
//	coder-fido2 derive-key --config-dir <dir>
//
// All subcommands read an optional PIN from stdin (one line). An empty
// line means no PIN (touch-only). This avoids exposing the PIN in
// process arguments or shell history.
//
// The "register" subcommand creates a FIDO2 credential with the
// HMAC-secret extension and writes the credential metadata to the
// config directory.
//
// The "derive-key" subcommand performs a FIDO2 assertion and prints
// the 32-byte HMAC-derived key as hex to stdout.
package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	libfido2 "github.com/keys-pub/go-libfido2"
)

// Exit code 2 signals a touch timeout to the caller, allowing retry.
const exitCodeTimeout = 2

// Exit code 3 signals that the security key requires a PIN.
const exitCodePinRequired = 3

const credentialFileName = "fido2-credential.json"

// storedCredential mirrors cli/fido2.StoredCredential. The JSON format
// must stay in sync between the helper and the main CLI.
type storedCredential struct {
	CredentialID []byte `json:"credential_id"`
	Salt         []byte `json:"salt"`
	RPId         string `json:"rp_id"`
}

func loadCredential(configDir string) (*storedCredential, error) {
	data, err := os.ReadFile(filepath.Join(configDir, credentialFileName))
	if err != nil {
		return nil, err
	}
	var cred storedCredential
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, fmt.Errorf("unmarshal credential: %w", err)
	}
	return &cred, nil
}

func saveCredential(configDir string, cred *storedCredential) error {
	data, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal credential: %w", err)
	}
	path := filepath.Join(configDir, credentialFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// isTouchTimeout returns true if the error is due to the user not
// touching the YubiKey in time.
func isTouchTimeout(err error) bool {
	return errors.Is(err, libfido2.ErrOperationDenied) ||
		errors.Is(err, libfido2.ErrActionTimeout) ||
		strings.Contains(err.Error(), "operation denied") ||
		strings.Contains(err.Error(), "timed out")
}

// isPinRequired returns true if the error indicates a PIN is needed.
func isPinRequired(err error) bool {
	return errors.Is(err, libfido2.ErrPinRequired) ||
		errors.Is(err, libfido2.ErrPinInvalid) ||
		strings.Contains(err.Error(), "pin required") ||
		strings.Contains(err.Error(), "pin invalid")
}

// exitForFIDOError prints the error and exits with the appropriate
// code: 2 for touch timeout, 3 for PIN required, 1 for everything else.
func exitForFIDOError(context string, err error) {
	if isTouchTimeout(err) {
		fmt.Fprintf(os.Stderr, "error: %s: touch timed out (try again)\n", context)
		os.Exit(exitCodeTimeout)
	}
	if isPinRequired(err) {
		fmt.Fprintf(os.Stderr, "error: %s: security key requires a PIN\n", context)
		os.Exit(exitCodePinRequired)
	}
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", context, err)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "register":
		cmdRegister()
	case "derive-key":
		cmdDeriveKey()
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: coder-fido2 <register|derive-key> --config-dir <dir>\n")
	os.Exit(1)
}

func parseConfigDir() string {
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--config-dir" {
			i++
			if i < len(os.Args) {
				return os.Args[i]
			}
		}
	}
	fmt.Fprintf(os.Stderr, "error: --config-dir is required\n")
	os.Exit(1)
	return ""
}

// readPIN reads an optional PIN from stdin. Empty line means no PIN.
func readPIN() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimRight(scanner.Text(), "\r\n")
	}
	return ""
}

func discoverDevice() (*libfido2.Device, *libfido2.DeviceLocation) {
	locs, err := libfido2.DeviceLocations()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: enumerate devices: %v\n", err)
		os.Exit(1)
	}
	if len(locs) == 0 {
		fmt.Fprintf(os.Stderr, "error: no FIDO2 device found; plug in your security key\n")
		os.Exit(1)
	}
	device, err := libfido2.NewDevice(locs[0].Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: open device: %v\n", err)
		os.Exit(1)
	}
	return device, locs[0]
}

func mustRandBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Fprintf(os.Stderr, "error: generate random bytes: %v\n", err)
		os.Exit(1)
	}
	return b
}

func cmdRegister() {
	configDir := parseConfigDir()
	pin := readPIN()
	device, loc := discoverDevice()

	fmt.Fprintf(os.Stderr, "Found: %s %s at %s\n", loc.Manufacturer, loc.Product, loc.Path)
	fmt.Fprintf(os.Stderr, "Touch your security key...\n")

	cdh := mustRandBytes(32)
	uid := mustRandBytes(32)
	salt := mustRandBytes(32)

	attest, err := device.MakeCredential(
		cdh,
		libfido2.RelyingParty{ID: "coder-cli", Name: "Coder CLI"},
		libfido2.User{ID: uid, Name: "coder"},
		libfido2.ES256,
		pin,
		&libfido2.MakeCredentialOpts{
			Extensions: []libfido2.Extension{libfido2.HMACSecretExtension},
			RK:         libfido2.Default,
		},
	)
	if err != nil {
		exitForFIDOError("make credential", err)
	}

	cred := &storedCredential{
		CredentialID: attest.CredentialID,
		Salt:         salt,
		RPId:         "coder-cli",
	}
	if err := saveCredential(configDir, cred); err != nil {
		fmt.Fprintf(os.Stderr, "error: save credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Credential registered (ID: %s...)\n",
		hex.EncodeToString(attest.CredentialID[:min(8, len(attest.CredentialID))]))
}

func cmdDeriveKey() {
	configDir := parseConfigDir()
	pin := readPIN()
	device, _ := discoverDevice()

	cred, err := loadCredential(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: load credential: %v\n", err)
		os.Exit(1)
	}

	cdh := mustRandBytes(32)

	assertion, err := device.Assertion(
		cred.RPId,
		cdh,
		[][]byte{cred.CredentialID},
		pin,
		&libfido2.AssertionOpts{
			Extensions: []libfido2.Extension{libfido2.HMACSecretExtension},
			UP:         libfido2.True,
			HMACSalt:   cred.Salt,
		},
	)
	if err != nil {
		exitForFIDOError("assertion", err)
	}

	if len(assertion.HMACSecret) != 32 {
		fmt.Fprintf(os.Stderr, "error: HMAC secret length %d, expected 32\n", len(assertion.HMACSecret))
		os.Exit(1)
	}

	// Print the derived key as hex to stdout.
	fmt.Print(hex.EncodeToString(assertion.HMACSecret))
}
