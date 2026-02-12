# FIDO2 Hardware Key Authentication

FIDO2 protects workspace connections (SSH, port forwarding) with a hardware
security key such as a YubiKey. When enabled, connecting to a workspace requires
a physical touch on the key. Management operations like listing or starting
workspaces work without the key.

> [!NOTE]
> This feature requires the `coder-fido2` helper binary, which is built
> separately from the main Coder CLI. See [Build the helper](#build-the-helper)
> below.

## How it works

FIDO2 authentication uses a two-token model:

- **Operate token**: Used for normal API calls (list, start, stop workspaces).
  Stored as the CLI session token. Does not grant SSH access.
- **Connect token**: Used only for workspace connections (SSH, port forwarding).
  Encrypted on disk with a key derived from your hardware security key.
  Decrypting it requires physical interaction with the key.

The connect token is encrypted using AES-256-GCM with a 32-byte key derived
from the FIDO2 `hmac-secret` extension. The key derivation requires a FIDO2
assertion (touch), so the token cannot be decrypted without physical access to
the registered security key.

## Prerequisites

- A FIDO2-compatible security key with the `hmac-secret` extension (e.g.
  YubiKey 5 series)
- The `coder-fido2` helper binary on your `PATH`
- `libfido2` installed on your system

## Build the helper

The `coder-fido2` helper binary uses CGo and `libfido2` for USB HID
communication with the security key. Install the dependencies and build:

```sh
# macOS
brew install libfido2

# Build the helper
CGO_ENABLED=1 go build -o coder-fido2 ./cmd/coder-fido2/

# Install on PATH
sudo cp coder-fido2 /usr/local/bin/
```

On Linux, install `libfido2-dev` (Debian/Ubuntu) or `libfido2-devel`
(Fedora/RHEL) before building.

## Set up FIDO2

Run `coder login` with the `--fido2` flag. If no credential is registered, the
CLI registers one automatically:

```sh
coder login https://coder.example.com --fido2
```

On first setup, this requires two touches on the security key:

1. **Touch 1**: Registers a FIDO2 credential on the key.
2. **Touch 2**: Derives the encryption key to protect the connect token.

Subsequent logins with `--fido2` require only one touch (key derivation).

If a touch times out, the CLI retries automatically up to three times.

## Connect to a workspace

After FIDO2 setup, workspace connections require touching the security key:

```console
$ coder ssh my-workspace
Touch your security key to connect...
```

Management commands work without the key:

```sh
coder list           # No touch required
coder start my-ws    # No touch required
coder ssh my-ws      # Touch required
```

## Manage credentials

Check FIDO2 status:

```sh
coder fido2 status
```

Register a credential manually (without `login --fido2`):

```sh
coder fido2 register
```

Remove the credential and encrypted token:

```sh
coder fido2 remove
```

After removal, workspace connections use the normal session token without
requiring the security key.

## Token expiry and refresh

The connect token has the same lifetime as a standard API token (default 7
days). When it expires, the next `coder ssh` creates a new connect token
automatically using the operate token, then encrypts it with a fresh FIDO2 key
derivation (one touch).

## Security considerations

- The FIDO2 mechanism is **client-side only**. The server issues standard API
  tokens and has no awareness of the security key. The key protects the connect
  token at rest on the client machine.
- The operate token has broad API access but **cannot** establish SSH
  connections. The connect token is scoped to `coder:workspaces.access`
  (SSH and application access only).
- The credential metadata (credential ID and salt) is stored in the Coder
  config directory as `fido2-credential.json`. The encrypted connect token is
  stored as `fido2-connect-token.enc`. Neither file is useful without the
  registered security key.
- PIN entry is supported. If your security key has a PIN set, the CLI prompts
  for it during registration. The PIN is passed to the helper via stdin, never
  as a command-line argument.

## Limitations

- Requires the separate `coder-fido2` helper binary (CGo dependency).
- First-time setup requires two touches on the security key. Subsequent logins
  and connections require one touch each.
- The security boundary is the client, not the server. A future server-side
  WebAuthn challenge would provide stronger guarantees.
