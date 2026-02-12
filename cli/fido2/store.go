// Package fido2 provides FIDO2 HMAC-secret-based encryption for protecting
// session tokens with a hardware security key (e.g. YubiKey).
//
// The HMAC-secret extension allows a FIDO2 device to derive a deterministic
// symmetric key from a credential and a salt. This key is used to
// AES-256-GCM encrypt/decrypt a session token stored on disk. The token
// cannot be decrypted without physical interaction with the security key.
//
// Device interaction (registration, key derivation) requires the
// coder-fido2 helper binary, which is built separately with CGo and
// libfido2. The main coder binary stays pure Go.
package fido2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Files are stored in the same config directory as url/session, which
// is a single global directory (not namespaced per server). This is
// consistent with how the CLI stores session state — only one active
// deployment at a time. Logging into a different server with --fido2
// overwrites the previous credential and token.
const (
	// nolint:gosec // G101: not a credential, just a filename.
	credentialFileName     = "fido2-credential.json"
	encryptedTokenFileName = "fido2-connect-token.enc"
)

// ZeroBytes overwrites a byte slice with zeros. Use to clear derived keys
// from memory after use.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// StoredCredential holds the FIDO2 credential metadata persisted to disk.
type StoredCredential struct {
	CredentialID []byte `json:"credential_id"`
	Salt         []byte `json:"salt"`
	RPId         string `json:"rp_id"`
}

// Encrypt encrypts a token with a 32-byte key and writes to configDir.
// The file format is: nonce (12 bytes) || AES-256-GCM ciphertext.
func Encrypt(key []byte, token string, configDir string) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return xerrors.Errorf("create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return xerrors.Errorf("create GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return xerrors.Errorf("generate nonce: %w", err)
	}
	// Seal appends ciphertext to nonce, producing nonce||ciphertext.
	out := gcm.Seal(nonce, nonce, []byte(token), nil)

	path := filepath.Join(configDir, encryptedTokenFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return xerrors.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(path, out, 0o600)
}

// Decrypt reads the encrypted token from configDir and decrypts it.
func Decrypt(key []byte, configDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(configDir, encryptedTokenFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", xerrors.Errorf("read encrypted token: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", xerrors.Errorf("create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", xerrors.Errorf("create GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", xerrors.New("encrypted token file is too short")
	}
	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", xerrors.Errorf("decrypt token (wrong key or corrupted data): %w", err)
	}
	return string(plaintext), nil
}

// LoadCredential reads stored FIDO2 credential metadata from configDir.
func LoadCredential(configDir string) (*StoredCredential, error) {
	data, err := os.ReadFile(filepath.Join(configDir, credentialFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, os.ErrNotExist
		}
		return nil, xerrors.Errorf("read credential: %w", err)
	}
	var cred StoredCredential
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, xerrors.Errorf("unmarshal credential: %w", err)
	}
	return &cred, nil
}

// SaveCredential writes FIDO2 credential metadata to configDir.
func SaveCredential(configDir string, cred *StoredCredential) error {
	data, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return xerrors.Errorf("marshal credential: %w", err)
	}
	path := filepath.Join(configDir, credentialFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return xerrors.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// HasCredential returns true if a FIDO2 credential exists in configDir.
func HasCredential(configDir string) bool {
	_, err := os.Stat(filepath.Join(configDir, credentialFileName))
	return err == nil
}

// HasEncryptedToken returns true if an encrypted connect token exists.
func HasEncryptedToken(configDir string) bool {
	_, err := os.Stat(filepath.Join(configDir, encryptedTokenFileName))
	return err == nil
}

// RemoveCredential deletes the credential and encrypted token files.
func RemoveCredential(configDir string) error {
	for _, name := range []string{credentialFileName, encryptedTokenFileName} {
		if err := os.Remove(filepath.Join(configDir, name)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
