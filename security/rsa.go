package security

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func GetSHA1Fingerprint(pubKey []byte) string {
	s := sha1.New()
	s.Write(pubKey)
	return fmt.Sprintf("%x", s.Sum(nil))
}

func GetX509CertSHA1Fingerprint(certFile string) string {
	b, _ := ioutil.ReadFile(certFile)
	block, _ := pem.Decode(b)
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(rsaPublicKey)
	return GetSHA1Fingerprint(pubKeyBytes)
}
