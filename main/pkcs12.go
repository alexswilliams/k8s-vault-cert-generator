package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	log "github.com/sirupsen/logrus"

	"software.sslmate.com/src/go-pkcs12"
)

// Encodes the certificates into a PKCS12 - note that a default password is used (as is convention).
// The underlying library makes the suggestion that you retain the default password and secure the whole PFX
func createPkcs12(keyPEM string, certPEM string, caCertPEMs []string) ([]byte, error) {
	var randomSource = rand.Reader

	var keyDER, _ = pem.Decode([]byte(keyPEM))
	panicOnBadCert(keyDER, "Failed to encode private key PEM into DER: %s", keyPEM)
	var key, err = x509.ParsePKCS1PrivateKey(keyDER.Bytes)
	panicOnError(err, "Failed to parse DER-encoded private key whilst building RSA private key object")

	var certDER, _ = pem.Decode([]byte(certPEM))
	panicOnBadCert(certDER, "Failed to encode certificate PEM into DER: %s", certPEM)
	var cert, err2 = x509.ParseCertificate(certDER.Bytes)
	panicOnError(err2, "Failed to parse DER-encoded certificate whilst building X509 certificate object")

	caCerts := make([]*x509.Certificate, len(caCertPEMs))
	for index, caCertPEM := range caCertPEMs {
		var caCertDER, _ = pem.Decode([]byte(caCertPEM))
		panicOnBadCert(caCertDER, "Failed to encode CA certificate PEM into DER: %s", caCertPEM)
		var caCert, err3 = x509.ParseCertificate(caCertDER.Bytes)
		panicOnError(err3, "Failed to parse DER-encoded CA certificate whilst building X509 certificate object")
		caCerts[index] = caCert
	}

	return pkcs12.Encode(randomSource, key, cert, caCerts, pkcs12.DefaultPassword)
}

func panicOnBadCert(cert interface{}, msg string, v ...interface{}) {
	if cert == nil {
		log.Errorf(msg, v...)
		panic("Could not decode PEM certificate into DER-encoded object.")
	}
}
