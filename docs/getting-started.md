# Getting Started

This guide will walk you through setting up `cassh` from scratch.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Microsoft Entra ID tenant** - For SSO authentication
- **GitHub Enterprise instance** - Where your repos live
- **A server** - To run cassh (see [Deployment](deployment.md))

## Quick Start (Development)

For local development and testing:

```bash
# Clone the repository
git clone https://github.com/shawntz/cassh.git
cd cassh

# Install dependencies
make deps

# Generate a development CA key
make dev-ca

# Run the server in dev mode (mock auth)
make dev-server
```

The server starts at `http://localhost:8080` with mock authentication enabled.

## Quick Start (Production)

For production deployment:

1. **Generate CA keys** (see [Server Setup](server-setup.md#generate-ca-keys))
2. **Create Entra app** (see [Server Setup](server-setup.md#create-microsoft-entra-app))
3. **Configure server** (see [Server Setup](server-setup.md#server-configuration))
4. **Deploy** (see [Deployment](deployment.md))
5. **Distribute client** (see [Client Distribution](client.md))

## Project Structure

```
cassh/
├── cmd/
│   ├── cassh-server/    # Web server (OIDC + cert signing)
│   ├── cassh-menubar/   # macOS menu bar app
│   └── cassh-cli/       # Headless CLI
├── internal/
│   ├── ca/              # Certificate authority logic
│   ├── config/          # Configuration handling
│   ├── memes/           # Meme content for landing page
│   └── oidc/            # Microsoft Entra ID integration
├── packaging/
│   └── macos/           # macOS distribution files
└── docs/                # Documentation (you are here)
```

## Build Commands

```bash
# Build all binaries
make build

# Build individual components
make server      # cassh-server
make menubar     # cassh-menubar (macOS)
make cli         # cassh CLI

# Run tests
make test

# Build macOS app bundle
make app-bundle

# Create DMG installer (requires sudo)
sudo make dmg

# Create PKG for MDM
make pkg
```

## Next Steps

1. [Server Setup](server-setup.md) - Configure CA and Entra
2. [Deployment](deployment.md) - Deploy to production
3. [Client Distribution](client.md) - Distribute to users
