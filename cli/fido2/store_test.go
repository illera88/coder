package fido2_test

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/cli/fido2"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	dir := t.TempDir()
	token := "test-session-token-abc123"

	err = fido2.Encrypt(key, token, dir)
	require.NoError(t, err)

	got, err := fido2.Decrypt(key, dir)
	require.NoError(t, err)
	require.Equal(t, token, got)
}

func TestDecryptWrongKey(t *testing.T) {
	t.Parallel()

	key1 := make([]byte, 32)
	_, _ = rand.Read(key1)
	key2 := make([]byte, 32)
	_, _ = rand.Read(key2)

	dir := t.TempDir()
	err := fido2.Encrypt(key1, "secret", dir)
	require.NoError(t, err)

	_, err = fido2.Decrypt(key2, dir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong key or corrupted data")
}

func TestDecryptCorruptedFile(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, _ = rand.Read(key)

	dir := t.TempDir()
	// Write garbage that's too short to contain a nonce.
	err := os.WriteFile(filepath.Join(dir, "fido2-connect-token.enc"), []byte("short"), 0o600)
	require.NoError(t, err)

	_, err = fido2.Decrypt(key, dir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "too short")
}

func TestCredentialSaveLoad(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cred := &fido2.StoredCredential{
		CredentialID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Salt:         make([]byte, 32),
		RPId:         "coder-cli",
	}
	_, _ = rand.Read(cred.Salt)

	err := fido2.SaveCredential(dir, cred)
	require.NoError(t, err)

	loaded, err := fido2.LoadCredential(dir)
	require.NoError(t, err)
	require.Equal(t, cred.CredentialID, loaded.CredentialID)
	require.Equal(t, cred.Salt, loaded.Salt)
	require.Equal(t, cred.RPId, loaded.RPId)
}

func TestHasCredential(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.False(t, fido2.HasCredential(dir))

	err := fido2.SaveCredential(dir, &fido2.StoredCredential{
		CredentialID: []byte{1},
		Salt:         make([]byte, 32),
		RPId:         "test",
	})
	require.NoError(t, err)
	require.True(t, fido2.HasCredential(dir))
}

func TestHasEncryptedToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, _ = rand.Read(key)
	dir := t.TempDir()

	require.False(t, fido2.HasEncryptedToken(dir))

	err := fido2.Encrypt(key, "token", dir)
	require.NoError(t, err)
	require.True(t, fido2.HasEncryptedToken(dir))
}

func TestRemoveCredential(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, _ = rand.Read(key)
	dir := t.TempDir()

	err := fido2.SaveCredential(dir, &fido2.StoredCredential{
		CredentialID: []byte{1},
		Salt:         make([]byte, 32),
		RPId:         "test",
	})
	require.NoError(t, err)
	err = fido2.Encrypt(key, "token", dir)
	require.NoError(t, err)

	require.True(t, fido2.HasCredential(dir))
	require.True(t, fido2.HasEncryptedToken(dir))

	err = fido2.RemoveCredential(dir)
	require.NoError(t, err)

	require.False(t, fido2.HasCredential(dir))
	require.False(t, fido2.HasEncryptedToken(dir))
}
