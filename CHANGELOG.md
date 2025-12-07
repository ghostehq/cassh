# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-12-06

### Initial Release! 🎉

`cassh`: SSH Key & Certificate Manager for GitHub.

> Whether you're a solo developer managing personal projects or an enterprise team with hundreds of engineers,
cassh brings security best practices to your SSH workflow — without the complexity.


### Summary of Core Features

- macOS menu bar app with status indicator (green/yellow/red)
- **Enterprise mode**: SSH certificates signed by internal CA (12-hour default validity)
- **Personal mode**: SSH keys via `gh` CLI with automatic rotation
- OIDC authentication with Microsoft Entra ID (enterprise)
- Web server with meme landing page (LSP & Flash Slothmore)
- Development mode for local testing with mock authentication

### Added

- **Personal GitHub.com support**: Manage SSH keys for personal GitHub accounts without a server
  - Automatic key generation (Ed25519)
  - Key upload via `gh` CLI integration
  - Configurable rotation policies (4 hours to 90 days)
  - Automatic SSH config management
- **Multi-account support**: Manage multiple GitHub accounts (enterprise + personal) from a single menu bar app
- **Setup wizard**: First-run configuration wizard for adding connections
- **Automatic SSH config setup**: Automatically configures `~/.ssh/config` for both enterprise and personal connections
- **System notifications**: macOS notifications for certificate/key activation, expiring soon, and expired states
- **GitHub Enterprise URL in policy**: Added `[github] enterprise_url` field to policy config for SSH config auto-setup
- **Build release script**: Added `scripts/build-release` one-liner script to build all packages after configuring policy
- **Web page footer**: Added footer to landing and success pages with GitHub, Docs, Sponsor links, and copyright
- **Setup CTA banner**: Added "Deploy `cassh` for your team" call-to-action on landing page linking to getting started guide
- **SSH Clone URL setup for Enterprise**: Enterprise setup wizard now accepts any SSH clone URL from your GitHub Enterprise (e.g., `user_123@github.company.com:org/repo.git`) and automatically extracts the hostname and SCIM-provisioned SSH username
- **Per-connection Git identity management**: Configure separate `user.name` and `user.email` for each GitHub connection
  - Uses Git's `includeIf` with `hasconfig:remote.*.url` to automatically apply the correct identity based on repo remote URL
  - Per-connection gitconfig files stored in `~/.config/cassh/gitconfig-{connection-id}`
  - Git identity fields (optional) added to both Enterprise and Personal setup wizards
- **GitHub Enterprise SSH certificate extensions**: Certificates now include the required `login@HOSTNAME=USERNAME` extension for GitHub Enterprise authentication
- **Custom URL scheme (`cassh://`)**: Certificate installation from HTTPS pages now uses the `cassh://install-cert` URL scheme to bypass mixed-content browser restrictions
- **SCIM username support**: SSH config for enterprise connections now uses the SCIM-provisioned username (from SSH clone URL) instead of hardcoded `git`

### Menu Bar App

- Terminal icon with macOS template icon support (auto dark/light mode)
- Status indicators in dropdown menu (green/yellow/red)
- One-click certificate generation/renewal
- Auto SSH key generation (Ed25519)
- ssh-agent integration for certificate loading
- Multi-connection dropdown with individual status per account

### Configuration

- Environment variable support for cloud deployment
- Split configuration model (policy vs user preferences)
- TOML-based configuration files
- Configurable certificate validity period
- User-configurable key rotation policies for personal accounts

### Deployment

- Dockerfile for containerized deployments
- `render.yaml` for Render.com infrastructure-as-code
- Makefile with comprehensive build targets
- Support for Fly.io, Railway, Render, and self-hosted VPS

### Distribution

- PKG installer for MDM deployment and Homebrew (Jamf, Kandji, etc.)
- macOS app bundle with embedded policy
- LaunchAgent for auto-start on login
- Homebrew Cask support (`brew install --cask cassh`)
- **Liquid glass icon**: macOS 15+ dynamic icon with translucency effects (via Icon Composer)
- `make icon` target to compile `.icon` bundle to `Assets.car` using `actool`
- Fallback `.icns` for older macOS versions

### Documentation

- MkDocs Material documentation site
- GitHub Pages deployment via GitHub Actions
- Comprehensive guides: getting started, server setup, deployment, client distribution
- Configuration reference with all options
- Security best practices and threat model
- Project roadmap with planned features

### CI/CD

- GitHub Actions workflow for releases (triggered on `v*` tags)
- GitHub Actions workflow for documentation deployment
- Automated changelog parsing for release notes
- macOS PKG signing and notarization

### Community

- Apache 2.0 license
- Code of Conduct (Contributor Covenant 2.0)
- Contributing guidelines
- Security policy with vulnerability reporting
- Issue templates (bug report, feature request)
- Pull request template
- GitHub Sponsors funding configuration

### Security

- Split configuration model (IT policy vs user preferences)
- CSRF protection using cryptographic state parameter
- Nonce verification to prevent replay attacks
- Policy bundled in signed app bundles (enterprise mode)
- Sensitive values loaded from environment variables in production
- Loopback listener restricted to localhost connections

### Fixed

- **golangci-lint v2 configuration**: Added `//go:build darwin` build tag to menubar app to fix Linux CI builds
- **CA private key parsing**: Handle escaped newlines (`\n`) in environment variables for cloud deployments (Render, etc.)
- **Mixed content blocking**: Browser security blocked HTTP localhost requests from HTTPS cassh server pages; now uses custom URL scheme
- **Certificate rejection on GHE**: GitHub Enterprise rejected certificates missing the `login@` extension; certificates now include proper extensions
- **SSH connection failures**: Enterprise SSH connections failed with "Permission denied" when using wrong SSH username; now correctly uses SCIM-provisioned username from clone URL
- **Double app launch on install**: PKG postinstall script was launching two instances; fixed by removing redundant `open` command (LaunchAgent handles startup)
