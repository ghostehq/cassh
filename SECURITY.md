# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in `cassh`, please report it by emailing me directly at **<SecOps@shawnschwartz.com>**.

Please include:

- A description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)

### What to Expect

- **Acknowledgment**: I will do my best to acknowledge receipt within 48-72 hours
- **Assessment**: I will assess the vulnerability and determine its severity
- **Updates**: I will keep you informed of critical progress updates during issue resolution
- **Resolution**: I will do my best to resolve critical vulnerabilities within 7 days
- **Credit**: I will credit you in the release notes (unless you prefer anonymity, let me know, otherwise you will be credited publicly!)

## Security Model

`cassh` is designed with security as a core principle:

### Certificate Authority

- CA private keys should never be exposed to clients
- Certificates are time-bound (default 12 hours) to limit exposure
- Certificate principals are controlled by server-side policy

### Authentication

- OIDC authentication via Microsoft Entra ID
- CSRF protection using state parameter
- Nonce verification to prevent replay attacks

### Configuration

- Split configuration model separates IT policy from user preferences
- Policy configuration is bundled in signed app bundles (enterprise)
- Sensitive values loaded from environment variables in production

### Transport

- All production traffic should use HTTPS
- Loopback listener only accepts connections from localhost

## Best Practices for Deployment

1. **Protect the CA private key** - Use hardware security modules (HSM) or secure key management in production
2. **Use HTTPS** - Never deploy the server without TLS
3. **Rotate certificates regularly** - The 12-hour default provides good balance
4. **Monitor certificate issuance** - Log all certificate signing events
5. **Restrict access** - Use network policies to limit who can reach the server
