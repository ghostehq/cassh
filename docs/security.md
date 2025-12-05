# Security

cassh is a privileged authentication system. This page covers security considerations and best practices.

!!! danger "Critical Security Note"
    The CA private key can sign certificates that GitHub Enterprise will trust. Protect it accordingly.

## Security Model

### Certificate Authority

- CA private keys should **never** be exposed to clients
- Certificates are time-bound (default 12 hours) to limit exposure
- Certificate principals are controlled by server-side policy
- Key ID includes user email and timestamp for auditing

### Authentication

- OIDC authentication via Microsoft Entra ID
- CSRF protection using cryptographic state parameter
- Nonce verification prevents replay attacks
- State tokens expire after 10 minutes

### Configuration

- Split configuration model separates IT policy from user preferences
- Policy configuration is bundled in signed app bundles (enterprise)
- Sensitive values loaded from environment variables in production
- User config cannot override security-critical settings

### Transport

- All production traffic should use HTTPS
- OAuth tokens are transmitted over TLS only
- Loopback listener only accepts connections from localhost

---

## Best Practices

### CA Key Management

| Practice | Recommendation |
|----------|----------------|
| Storage | Use HSM or secure key management service |
| Access | Limit to cassh server process only |
| Rotation | Rotate annually or after suspected compromise |
| Backup | Secure offline backup with encryption |

### Server Deployment

- [ ] **Use HTTPS** - Never deploy without TLS
- [ ] **Restrict network access** - Use firewall rules
- [ ] **Monitor logs** - Alert on unusual certificate issuance
- [ ] **Keep updated** - Apply security patches promptly

### Entra App Configuration

- [ ] **Restrict users** - Limit who can authenticate
- [ ] **Use conditional access** - Require MFA, compliant devices
- [ ] **Monitor sign-ins** - Review authentication logs
- [ ] **Rotate secrets** - Change client secret periodically

### Client Distribution

- [ ] **Code sign** - Sign app bundles with Developer ID
- [ ] **Notarize** - Submit to Apple for notarization
- [ ] **MDM deployment** - Use managed distribution
- [ ] **Policy bundling** - Embed policy in signed app

---

## Threat Model

### Threats Addressed

| Threat | Mitigation |
|--------|------------|
| Stolen SSH key | Certificates expire automatically |
| Lost laptop | No action required - cert expires |
| Employee offboarding | Revoke Entra access, certs expire |
| Key compromise | Limited blast radius (12 hours) |
| CSRF attacks | State parameter validation |
| Replay attacks | Nonce verification |

### Threats NOT Addressed

| Threat | Notes |
|--------|-------|
| Compromised CA key | Attacker can sign arbitrary certs |
| Compromised Entra tenant | Attacker can authenticate as any user |
| Local privilege escalation | SSH key accessible to local user |
| Active session hijacking | Browser-based auth flow |

---

## Incident Response

### Suspected CA Key Compromise

1. **Immediately** revoke the CA in GitHub Enterprise
2. Generate a new CA key pair
3. Add new CA to GitHub Enterprise
4. Redistribute updated client policy
5. Investigate breach scope

### Suspicious Certificate Issuance

1. Review server logs for issuance events
2. Check Entra sign-in logs for anomalies
3. Revoke suspicious user access in Entra
4. Consider reducing cert validity temporarily

---

## Auditing

### Server Logs

cassh logs all certificate issuance events:

```
2024-12-04 10:30:45 User authenticated: user@company.com (username)
2024-12-04 10:30:46 Certificate issued: cassh:user@company.com:1701689446
```

### Recommended Monitoring

- Certificate issuance rate (alert on spikes)
- Failed authentication attempts
- Unusual access patterns (time, location)
- Certificate validity modifications

---

## Reporting Vulnerabilities

**Do not report security vulnerabilities through public GitHub issues.**

Email **SecOps@shawnschwartz.com** with:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

I will do my best to respond within 48-72 hours.
