package testcert

import (
	"crypto/tls"
	"crypto/x509"
)

// LocalhostCertPool returns CertPool for validating LocalhostCert.
func LocalhostCertPool() *x509.CertPool {
	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(LocalhostCert) {
		panic("failed to AppendCertsFromPEM(LocalhostCert)")
	}
	return ca
}

// LocalhostCertificate returns Certificate.
func LocalhostCertificate() tls.Certificate {
	cert, err := tls.X509KeyPair(LocalhostCert, LocalhostKey)
	if err != nil {
		panic(err)
	}
	return cert
}
