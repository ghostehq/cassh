// Handles SSH cert generation and signing
// Certs issued with configurable validity (default 12 hours)
package ca

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// CertificateAuthority manages SSH cert signing
type CertificateAuthority struct {
	signer        ssh.Signer
	validityHours int
	principals    []string
}

// NewCA creates a new certificate authority from a private key
func NewCA(privateKeyPEM []byte, validityHours int, principals []string) (*CertificateAuthority, error) {
	signer, err := ssh.ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA private key: %w", err)
	}

	return &CertificateAuthority{
		signer:        signer,
		validityHours: validityHours,
		principals:    principals,
	}, nil
}

// SignPublicKey signs a user's public key, creating an SSH cert
func (ca *CertificateAuthority) SignPublicKey(userPubKey ssh.PublicKey, keyID string, username string) (*ssh.Certificate, error) {
	// Generate random serial
	serialBytes := make([]byte, 8)
	if _, err := rand.Read(serialBytes); err != nil {
		return nil, fmt.Errorf("failed to generate serial: %w", err)
	}
	serial := binary.BigEndian.Uint64(serialBytes)

	now := time.Now()
	validAfter := uint64(now.Unix())
	validBefore := uint64(now.Add(time.Duration(ca.validityHours) * time.Hour).Unix())

	// Determine principals
	principals := ca.principals
	if len(principals) == 0 {
		principals = []string{username}
	}

	cert := &ssh.Certificate{
		Key:             userPubKey,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           keyID,
		ValidPrincipals: principals,
		ValidAfter:      validAfter,
		ValidBefore:     validBefore,
		Permissions: ssh.Permissions{
			Extensions: map[string]string{
				"permit-agent-forwarding": "",
				"permit-port-forwarding":  "",
				"permit-pty":              "",
				"permit-user-rc":          "",
			},
		},
	}

	if err := cert.SignCert(rand.Reader, ca.signer); err != nil {
		return nil, fmt.Errorf("failed to sign certificate: %w", err)
	}

	return cert, nil
}

// GenerateKeyPair creates a new Ed25519 keypair for the user
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keypair: %w", err)
	}
	return pub, priv, nil
}

// MarshalCertificate converts a cert to the authorized_keys format
func MarshalCertificate(cert *ssh.Certificate) []byte {
	return ssh.MarshalAuthorizedKey(cert)
}

// MarshalPrivateKey converts an Ed25519 private key to PEM format
func MarshalPrivateKey(priv ed25519.PrivateKey) ([]byte, error) {
	// Use the ssh package to marshal to OpenSSH format
	block, err := ssh.MarshalPrivateKey(priv, "cassh generated key")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	return block.Bytes, nil
}

// ParsePublicKey parses an SSH public key from authorized_keys format
func ParsePublicKey(pubKeyBytes []byte) (ssh.PublicKey, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return pub, nil
}

// CertInfo contains human-readable cert information
type CertInfo struct {
	Serial      uint64
	KeyID       string
	Principals  []string
	ValidAfter  time.Time
	ValidBefore time.Time
	IsExpired   bool
	TimeLeft    time.Duration
}

// GetCertInfo extracts info from a cert for display
func GetCertInfo(cert *ssh.Certificate) *CertInfo {
	now := time.Now()
	validBefore := time.Unix(int64(cert.ValidBefore), 0)
	validAfter := time.Unix(int64(cert.ValidAfter), 0)

	return &CertInfo{
		Serial:      cert.Serial,
		KeyID:       cert.KeyId,
		Principals:  cert.ValidPrincipals,
		ValidAfter:  validAfter,
		ValidBefore: validBefore,
		IsExpired:   now.After(validBefore),
		TimeLeft:    validBefore.Sub(now),
	}
}

// ParseCertificate parses an SSH cert from file
func ParseCertificate(certBytes []byte) (*ssh.Certificate, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	cert, ok := pub.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("not a certificate")
	}

	return cert, nil
}
