package tlshandler

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler"
)

// source: https://shaneutt.com/blog/golang-ca-and-signed-cert-go/

type CertData struct {
	ftp_filehandler.FileBasic
	Organization  string         `json:"organisation"`
	Country       string         `json:"country"`
	Province      string         `json:"province"`
	Locality      string         `json:"locality"`
	StreetAddress string         `json:"street_address"`
	PostalCode    string         `json:"postal_code"`
	NotAfter      NotAfterStruct `json:"add_date"`
	IPAddrresses  []net.IP       `json:"ip_addresses"`
}

type NotAfterStruct struct {
	Years  int `json:"years"`
	Months int `json:"months"`
	Days   int `json:"days"`
}

type CertSetup struct {
	ftp_filehandler.FileBasic
	cert       *x509.Certificate
	CertData   *CertData       `json:"cert_data"`
	PrivKey    *rsa.PrivateKey `json:"private_key"`
	PEM        *bytes.Buffer   `json:"pem"`
	PrivKeyPEM *bytes.Buffer   `json:"private_key_pem"`
	TlsCert    *TlsCertJson    `json:"tls_cert"`
	err        error
}

type CAJson struct {
	cert       *x509.Certificate
	CertData   CertData        `json:"cert_data"`
	PrivKey    *rsa.PrivateKey `json:"private_key"`
	PEM        string          `json:"pem"`
	PrivKeyPEM string          `json:"private_key_pem"`

	err error
}

func (cs *CertSetup) NewCAJson() (caj *CAJson) {
	caj = &CAJson{
		CertData:   *cs.CertData,
		cert:       cs.cert,
		PrivKey:    cs.PrivKey,
		PEM:        cs.PEM.String(),
		PrivKeyPEM: cs.PrivKeyPEM.String(),
	}

	return
}

func (caj *CAJson) ToJSON() (s string, err error) {
	loc := "func (caj *CAJson)ToJSON()(s string, err error)"
	j, err := json.Marshal(caj)
	if err != nil {
		return "", ftp_context.NewLogItem(loc, true).SetAfter("json.Marshal(caj)").SetMessage(err.Error()).AppendParentError(err)
	}
	s = string(j)

	return
}

type tlc_cert tls.Certificate

type TlsCertJson struct {
	tls.Certificate
	Cert [][]byte `json:"certificate"`
	// PrivateKey contains the private key corresponding to the public key in
	// Leaf. This must implement crypto.Signer with an RSA, ECDSA or Ed25519 PublicKey.
	// For a server up to TLS 1.2, it can also implement crypto.Decrypter with
	// an RSA PublicKey.
	PrivKey crypto.PrivateKey `json:"private_key"`
	// SupportedSignatureAlgorithms is an optional list restricting what
	// signature algorithms the PrivateKey can be used for.
	SSA []tls.SignatureScheme `json:"supported_signature_algorithms"`
	// OCSPStaple contains an optional OCSP response which will be served
	// to clients that request it.
	OCSPS []byte `json:"OCSP_staple"`
	// SignedCertificateTimestamps contains an optional list of Signed
	// Certificate Timestamps which will be served to clients that request it.
	SCT [][]byte `json:"signed_certificate_timestamps"`
	// Leaf is the parsed form of the leaf certificate, which may be initialized
	// using x509.ParseCertificate to reduce per-handshake processing. If nil,
	// the leaf certificate will be parsed as needed.
	L *x509.Certificate `json:"leaf"`
}

func NewTlsCertJson(tlc tls.Certificate) *TlsCertJson {

	return &TlsCertJson{
		Certificate: (tlc),
		Cert:        tlc.Certificate,
		PrivKey:     tlc.PrivateKey,
		SSA:         tlc.SupportedSignatureAlgorithms,
		OCSPS:       tlc.OCSPStaple,
		SCT:         tlc.SignedCertificateTimestamps,
		L:           tlc.Leaf,
	}
}

func (cs CertSetup) HasErr() bool {
	return cs.err != nil
}

func (cs CertSetup) Error() string {
	return cs.err.Error()
}

func (cs CertSetup) UnderlyingError() error {
	return cs.err
}

// set up our CA certificate
func NewCA(org *CertData) (cs *CertSetup) {
	loc := "func NewCA(org *CertData) (cs *CertSetup)"
	cs = new(CertSetup)
	var err error
	cs.CertData = org
	cs.cert = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{org.Organization},
			Country:       []string{org.Country},
			Province:      []string{org.Province},
			Locality:      []string{org.Locality},
			StreetAddress: []string{org.StreetAddress},
			PostalCode:    []string{org.PostalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(org.NotAfter.Years, org.NotAfter.Months, org.NotAfter.Days),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	cs.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		cs.err = ftp_context.NewLogItem(loc, true).
			SetAfter("rsa.GenerateKey(rand.Reader, 4096)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, cs.cert, cs.cert, &cs.PrivKey.PublicKey, cs.PrivKey)
	if err != nil {
		log.Println(err)
		cs.err = ftp_context.NewLogItem(loc, true).
			Set("after", "x509.CreateCertificate(rand.Reader, cs.cert, cs.cert, cs.PrivKey.PublicKey, cs.PrivKey)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return

	}

	// pem encode
	cs.PEM = new(bytes.Buffer)
	pem.Encode(cs.PEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	cs.PrivKeyPEM = new(bytes.Buffer)
	pem.Encode(cs.PrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cs.PrivKey),
	})

	return

}

func (c_a CertSetup) NewServerCert(org *CertData) (cs *CertSetup) {
	loc := "func (cs *CertSetup) ServerKey()(cd *CertSetup) "
	var err error
	cs = new(CertSetup)

	// set up our server certificate
	cs.cert = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{org.Organization},
			Country:       []string{org.Country},
			Province:      []string{org.Province},
			Locality:      []string{org.Locality},
			StreetAddress: []string{org.StreetAddress},
			PostalCode:    []string{org.PostalCode},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(org.NotAfter.Years, org.NotAfter.Months, org.NotAfter.Days),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	cs.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		cs.err = ftp_context.NewLogItem(loc, true).Set("after", "rsa.GenerateKey(rand.Reader, 4096)").AppendParentError(err)
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cs.cert, c_a.cert, &cs.PrivKey.PublicKey, c_a.PrivKey)
	if err != nil {
		cs.err = ftp_context.NewLogItem(loc, true).Set("after", "rsa.GenerateKey(rand.Reader, 4096)").AppendParentError(err)
		return
	}

	cs.PEM = new(bytes.Buffer)
	pem.Encode(cs.PEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	cs.PrivKeyPEM = new(bytes.Buffer)
	pem.Encode(cs.PrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cs.PrivKey),
	})

	tlc_c, err := tls.X509KeyPair(cs.PEM.Bytes(), cs.PrivKeyPEM.Bytes())
	if err != nil {
		cs.err = ftp_context.NewLogItem(loc, true).Set("after", "tls.X509KeyPair(cs.PEM.Bytes(), cs.PrivKeyPEM.Bytes())").AppendParentError(err)
		return
	}

	cs.TlsCert = NewTlsCertJson(tlc_c)
	return cs
}

func (cs CertSetup) ServerTlsConfig() (tlc *tls.Config, err error) {
	loc := "func (cs CertSetup) ServerTlsConfig()(tlc *tls.Config,err error)"
	if cs.err != nil {
		err = ftp_context.NewLogItem(loc, true).SetMessage("invalid CertSetup State").AppendParentError(cs.err)
		return
	}

	tlc = &tls.Config{
		Certificates: []tls.Certificate{cs.TlsCert.Certificate},
	}

	return
}
