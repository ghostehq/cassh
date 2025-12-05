# Deployment

`cassh` server can be deployed to various platforms. Choose based on your needs:

| Platform | Cost | Complexity | Best For |
|----------|------|------------|----------|
| [Render](#render) | Free tier | Low | Quick setup, auto-deploy |
| [Fly.io](#flyio) | <$5/mo | Low | Edge deployment, scaling |
| [Railway](#railway) | <$5/mo | Low | Simple deploys |
| [Self-Hosted](#self-hosted-vps) | ~$5/mo | Medium | Full control |

## Render

Render offers a free tier with auto-deploy from GitHub.

### Step 1: Connect Repository

1. Push your repo to GitHub
2. Go to [render.com](https://render.com)
3. Click **New** → **Web Service**
4. Connect your GitHub repo

### Step 2: Configure Build

- **Build Command:** `go build -o cassh-server ./cmd/cassh-server`
- **Start Command:** `./cassh-server`

### Step 3: Set Environment Variables

Add these in the Render dashboard:

```
CASSH_SERVER_URL=https://your-app.onrender.com
CASSH_OIDC_CLIENT_ID=your-client-id
CASSH_OIDC_CLIENT_SECRET=your-client-secret
CASSH_OIDC_TENANT=your-tenant-id
CASSH_CA_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----...
```

### Step 4: Deploy

Click **Create Web Service**. Render will build and deploy automatically.

!!! tip "Using render.yaml"
    The repo includes a `render.yaml` for infrastructure-as-code deployment.

---

## Fly.io

### Step 1: Install Fly CLI

=== "macOS"
    ```bash
    brew install flyctl
    ```

=== "Linux/WSL"
    ```bash
    curl -L https://fly.io/install.sh | sh
    ```

### Step 2: Login

```bash
fly auth signup  # or fly auth login
```

### Step 3: Launch App

```bash
fly launch --name cassh-yourcompany --no-deploy
```

### Step 4: Set Secrets

```bash
fly secrets set CASSH_SERVER_URL="https://cassh-yourcompany.fly.dev"
fly secrets set CASSH_OIDC_CLIENT_ID="your-client-id"
fly secrets set CASSH_OIDC_CLIENT_SECRET="your-client-secret"
fly secrets set CASSH_OIDC_TENANT="your-tenant-id"
fly secrets set CASSH_CA_PRIVATE_KEY="$(cat ca_key)"
```

### Step 5: Configure fly.toml

```toml
app = "cassh-yourcompany"
primary_region = "sjc"

[build]

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0

[env]
  PORT = "8080"
```

### Step 6: Deploy

```bash
fly deploy
```

Your server is live at `https://cassh-yourcompany.fly.dev`

---

## Railway

### Step 1: Install CLI

```bash
npm install -g @railway/cli
```

### Step 2: Login and Initialize

```bash
railway login
railway init
```

### Step 3: Set Variables

```bash
railway variables set CASSH_SERVER_URL="https://your-app.up.railway.app"
railway variables set CASSH_OIDC_CLIENT_ID="your-client-id"
# ... set other variables
```

### Step 4: Deploy

```bash
railway up
```

---

## Self-Hosted (VPS)

Any $5/mo VPS works (DigitalOcean, Linode, Vultr, Hetzner).

### Step 1: Build on Server

```bash
git clone https://github.com/shawntz/cassh.git
cd cassh
go build -o cassh-server ./cmd/cassh-server
```

### Step 2: Create Systemd Service

```ini
# /etc/systemd/system/cassh.service
[Unit]
Description=cassh SSH Certificate Server
After=network.target

[Service]
Type=simple
User=cassh
WorkingDirectory=/opt/cassh
EnvironmentFile=/etc/cassh/cassh.env
ExecStart=/opt/cassh/cassh-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Step 3: Create Environment File

```bash
# /etc/cassh/cassh.env
CASSH_SERVER_URL=https://cassh.yourcompany.com
CASSH_OIDC_CLIENT_ID=your-client-id
CASSH_OIDC_CLIENT_SECRET=your-client-secret
CASSH_OIDC_TENANT=your-tenant-id
CASSH_CA_PRIVATE_KEY_PATH=/etc/cassh/ca_key
```

### Step 4: Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable cassh
sudo systemctl start cassh
```

### Step 5: Configure Reverse Proxy

Use nginx or Caddy for HTTPS termination.

**Caddy (recommended):**

```
cassh.yourcompany.com {
    reverse_proxy localhost:8080
}
```

**nginx:**

```nginx
server {
    listen 443 ssl http2;
    server_name cassh.yourcompany.com;
    
    ssl_certificate /etc/letsencrypt/live/cassh.yourcompany.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/cassh.yourcompany.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Docker

A Dockerfile is included for containerized deployments:

```bash
docker build -t cassh-server .
docker run -p 8080:8080 \
  -e CASSH_SERVER_URL="https://cassh.yourcompany.com" \
  -e CASSH_OIDC_CLIENT_ID="your-client-id" \
  -e CASSH_OIDC_CLIENT_SECRET="your-client-secret" \
  -e CASSH_OIDC_TENANT="your-tenant-id" \
  -e CASSH_CA_PRIVATE_KEY="$(cat ca_key)" \
  cassh-server
```

## Update Entra Redirect URI

After deployment, update your Entra app's redirect URI to match your production URL:

1. Go to Azure Portal → Entra ID → App registrations
2. Select your cassh app
3. Go to **Authentication**
4. Update the redirect URI to `https://your-production-url/auth/callback`
