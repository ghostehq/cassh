# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-12-04

### Added

- macOS menu bar app with certificate status indicator
- OIDC authentication with Microsoft Entra ID
- 12-hour SSH certificate signing
- Web server with meme landing page (LSP & Flash Slothmore)
- Development mode for local testing
- DMG and PKG installers for distribution

### Security

- Split configuration model (IT policy vs user preferences)
- CSRF protection for OIDC flow
- Nonce verification to prevent replay attacks
