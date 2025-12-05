# Server Setup

This guide covers generating your CA keys, creating the Microsoft Entra app, and configuring the cassh server.

## Generate CA Keys

The Certificate Authority (CA) key is used to sign all SSH certificates. **Keep the private key secure.**

```bash
ssh-keygen -t ed25519 -f ca_key -N "" -C "cassh-ca"
```

This creates two files:

| File | Purpose | Security |
|------|---------|----------|
| `ca_key` | Private key for signing | **SERVER ONLY** - never share |
| `ca_key.pub` | Public key | Upload to GitHub Enterprise |

!!! danger "Protect Your CA Private Key"
    Anyone with access to the CA private key can sign certificates that GitHub Enterprise will trust. Store it securely and limit access.

## Configure GitHub Enterprise

Add your CA public key to GitHub Enterprise so it trusts certificates signed by your CA:

1. Go to your **GitHub Enterprise admin panel**
2. Navigate to **Settings** → **SSH certificate authorities**
3. Click **Add CA**
4. Paste the contents of `ca_key.pub`
5. Save

!!! tip "Verify the CA"
    After adding, you can test by signing a certificate and attempting to clone a repo.

## Create Microsoft Entra App

cassh uses Microsoft Entra ID (Azure AD) for authentication.

### Step 1: Register the App

1. Go to [Azure Portal](https://portal.azure.com)
2. Navigate to **Microsoft Entra ID** → **App registrations**
3. Click **New registration**
4. Configure:
    - **Name:** `cassh`
    - **Supported account types:** Accounts in this organizational directory only
    - **Redirect URI:** `https://YOUR_SERVER_URL/auth/callback`
5. Click **Register**

### Step 2: Note the IDs

After registration, note these values from the **Overview** page:

- **Application (client) ID** - Used as `CASSH_OIDC_CLIENT_ID`
- **Directory (tenant) ID** - Used as `CASSH_OIDC_TENANT`

### Step 3: Create Client Secret

1. Go to **Certificates & secrets**
2. Click **New client secret**
3. Add a description and expiry
4. Click **Add**
5. **Copy the secret value immediately** - it's only shown once!

This becomes your `CASSH_OIDC_CLIENT_SECRET`.

### Step 4: Configure API Permissions (Optional)

By default, cassh only needs basic OpenID permissions. If you want to restrict access to specific groups:

1. Go to **API permissions**
2. Add **Microsoft Graph** → **User.Read** (usually already added)
3. For group-based access, add **GroupMember.Read.All**

## Server Configuration

cassh can be configured via TOML file, environment variables, or both. Environment variables take precedence.

### Option 1: TOML File

Create `cassh.policy.toml`:

```toml
server_base_url = "https://cassh.yourcompany.com"
cert_validity_hours = 12

[oidc]
client_id = "your-client-id"
client_secret = "your-client-secret"
tenant = "your-tenant-id"
redirect_url = "https://cassh.yourcompany.com/auth/callback"

[ca]
private_key_path = "./ca_key"

[github]
enterprise_url = "https://github.yourcompany.com"
```

### Option 2: Environment Variables

For cloud deployments, use environment variables:

```bash
export CASSH_SERVER_URL="https://cassh.yourcompany.com"
export CASSH_OIDC_CLIENT_ID="your-client-id"
export CASSH_OIDC_CLIENT_SECRET="your-client-secret"
export CASSH_OIDC_TENANT="your-tenant-id"
export CASSH_CA_PRIVATE_KEY="$(cat ca_key)"  # Full key content
export CASSH_CERT_VALIDITY_HOURS="12"
```

See [Configuration Reference](configuration.md) for all options.

## Run the Server

### Development Mode

```bash
# With dev CA (mock auth enabled)
make dev-server
```

### Production Mode

```bash
# With TOML config
./cassh-server

# Or with env vars
CASSH_SERVER_URL="https://..." ./cassh-server
```

The server listens on `:8080` by default. Use `CASSH_LISTEN_ADDR` to change.

## Verify Setup

1. Open `http://localhost:8080` in your browser
2. You should see the meme landing page
3. Click "Sign in with Microsoft"
4. Complete the authentication flow
5. You should receive a certificate

## Next Steps

- [Deployment](deployment.md) - Deploy to production
- [Client Distribution](client.md) - Distribute to users
