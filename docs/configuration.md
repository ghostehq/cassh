# Configuration Reference

`cassh` uses a split configuration model:

- **Policy Config** - IT-controlled settings (bundled in app or from env vars)
- **User Config** - Personal preferences (editable by user)

## Environment Variables

Environment variables take precedence over file configuration.

### Server Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `CASSH_SERVER_URL` | Public server URL | Yes | - |
| `CASSH_OIDC_CLIENT_ID` | Entra app client ID | Yes* | - |
| `CASSH_OIDC_CLIENT_SECRET` | Entra app client secret | Yes* | - |
| `CASSH_OIDC_TENANT` | Entra tenant ID | Yes* | - |
| `CASSH_OIDC_REDIRECT_URL` | OAuth callback URL | No | `{server_url}/auth/callback` |
| `CASSH_CA_PRIVATE_KEY` | CA private key content | Yes** | - |
| `CASSH_CA_PRIVATE_KEY_PATH` | Path to CA private key file | Yes** | - |
| `CASSH_CERT_VALIDITY_HOURS` | Certificate lifetime in hours | No | `12` |
| `CASSH_LISTEN_ADDR` | Server listen address | No | `:8080` |
| `CASSH_DEV_MODE` | Enable development mode | No | `false` |
| `CASSH_POLICY_PATH` | Path to policy TOML file | No | `cassh.policy.toml` |

*Required in production mode
**One of these is required in production mode

### Development Mode

Set `CASSH_DEV_MODE=true` to:

- Skip OIDC authentication (mock auth)
- Use local CA key for testing
- Enable verbose logging

---

## Policy Config (TOML)

The policy TOML file contains IT-controlled settings.

### Server Configuration

```toml
# Server settings
server_base_url = "https://cassh.yourcompany.com"
cert_validity_hours = 12
dev_mode = false

# OIDC / Microsoft Entra
[oidc]
client_id = "your-client-id"
client_secret = "your-client-secret"
tenant = "your-tenant-id"
redirect_url = "https://cassh.yourcompany.com/auth/callback"

# Certificate Authority
[ca]
private_key_path = "./ca_key"

# GitHub Enterprise
[github]
enterprise_url = "https://github.yourcompany.com"
allowed_orgs = ["your-org"]  # Optional: restrict to specific orgs
```

### Client Configuration

Clients only need the server URL:

```toml
server_base_url = "https://cassh.yourcompany.com"
```

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `server_base_url` | string | Public URL of cassh server |
| `cert_validity_hours` | int | Certificate lifetime (default: 12) |
| `dev_mode` | bool | Enable development mode |
| `oidc.client_id` | string | Entra application client ID |
| `oidc.client_secret` | string | Entra application secret |
| `oidc.tenant` | string | Entra tenant ID |
| `oidc.redirect_url` | string | OAuth callback URL |
| `ca.private_key_path` | string | Path to CA private key file |
| `github.enterprise_url` | string | GitHub Enterprise base URL |
| `github.allowed_orgs` | []string | Restrict access to these orgs |

---

## User Config

User preferences are stored in `~/Library/Application Support/cassh/config.toml` (macOS).

```toml
# Refresh interval for status checks (seconds)
refresh_interval_seconds = 30

# Play sound on cert expiration warning
notification_sound = true

# Preferred meme character: "lsp", "sloth", or "random"
preferred_meme = "random"

# SSH key paths (auto-managed)
ssh_key_path = "~/.ssh/cassh_id_ed25519"
ssh_cert_path = "~/.ssh/cassh_id_ed25519-cert.pub"
```

### Field Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `refresh_interval_seconds` | int | `30` | How often to check cert status |
| `notification_sound` | bool | `true` | Play sound on warnings |
| `preferred_meme` | string | `"random"` | Landing page character |
| `ssh_key_path` | string | `~/.ssh/cassh_id_ed25519` | SSH private key location |
| `ssh_cert_path` | string | `~/.ssh/cassh_id_ed25519-cert.pub` | Certificate location |

---

## Config Locations

### macOS

| Config | Location |
|--------|----------|
| Policy (bundled) | `cassh.app/Contents/Resources/cassh.policy.toml` |
| Policy (fallback) | `./cassh.policy.toml` |
| User config | `~/Library/Application Support/cassh/config.toml` |
| SSH key | `~/.ssh/cassh_id_ed25519` |
| SSH cert | `~/.ssh/cassh_id_ed25519-cert.pub` |

### Linux

| Config | Location |
|--------|----------|
| Policy | `./cassh.policy.toml` or `CASSH_POLICY_PATH` |
| User config | `~/.config/cassh/config.toml` |
| SSH key | `~/.ssh/cassh_id_ed25519` |
| SSH cert | `~/.ssh/cassh_id_ed25519-cert.pub` |

---

## Config Precedence

1. **Environment variables** (highest priority)
2. **Policy file** (TOML)
3. **Default values** (lowest priority)

For security-critical settings (CA key, OIDC secrets), policy always wins over user config.
