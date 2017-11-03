package security

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/ricottadeploy/common/x509cert"
)

type Fingerprint string

func (f *Fingerprint) Valid() bool {
	matched, err := regexp.MatchString("^[a-f0-9]{40}$", string(*f))
	if err != nil {
		log.Fatal("Error while matching regex: %s", err)
	}
	return matched
}

func GetPeerFingerprint(conn *tls.Conn) string {
	conn.Handshake()
	state := conn.ConnectionState()
	peerCert := state.PeerCertificates[0]
	publicKey := peerCert.PublicKey.(*rsa.PublicKey)
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	fingerPrint := x509cert.GetFingerprintSHA1(pubKeyBytes)
	return fingerPrint
}

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
