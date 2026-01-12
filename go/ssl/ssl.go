package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	nethttp "net/http"
	"strings"

	"github.com/github/orchestrator/go/config"
	"github.com/go-martini/martini"
	"github.com/howeyc/gopass"
	"github.com/openark/golib/log"
)

var cipherSuites = []uint16{
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
}

// Determine if a string element is in a string array
func HasString(elem string, arr []string) bool {
	for _, s := range arr {
		if s == elem {
			return true
		}
	}
	return false
}

// NewTLSConfig returns an initialized TLS configuration suitable for client
// authentication. If caFile is non-empty, it will be loaded.
func NewTLSConfig(caFile string, verifyCert bool) (*tls.Config, error) {
	var c tls.Config

	// Set to TLS 1.2 as a minimum.  This is overridden for mysql communication
	c.MinVersion = tls.VersionTLS12
	// Remove insecure ciphers from the list
	c.CipherSuites = cipherSuites
	c.PreferServerCipherSuites = true

	if verifyCert {
		log.Info("verifyCert requested, client certificates will be verified")
		c.ClientAuth = tls.VerifyClientCertIfGiven
	}
	if caFile != "" {
		data, err := ioutil.ReadFile(caFile)
		if err != nil {
			return &c, err
		}
		c.ClientCAs = x509.NewCertPool()
		if !c.ClientCAs.AppendCertsFromPEM(data) {
			return &c, errors.New("No certificates parsed")
		}
		log.Info("Read in CA file:", caFile)
	}
	c.BuildNameToCertificate()
	return &c, nil
}

// Verify that the OU of the presented client certificate matches the list
// of Valid OUs
func Verify(r *nethttp.Request, validOUs []string) error {
	if strings.Contains(r.URL.String(), config.Config.StatusEndpoint) && !config.Config.StatusOUVerify {
		return nil
	}
	if r.TLS == nil {
		return errors.New("No TLS")
	}
	for _, chain := range r.TLS.VerifiedChains {
		s := chain[0].Subject.OrganizationalUnit
		log.Debug("All OUs:", strings.Join(s, " "))
		for _, ou := range s {
			log.Debug("Client presented OU:", ou)
			if HasString(ou, validOUs) {
				log.Debug("Found valid OU:", ou)
				return nil
			}
		}
	}
	log.Error("No valid OUs found")
	return errors.New("Invalid OU")
}

type verifyOUsFunc func(r *nethttp.Request, validOUs []string) error

func VerifyOUs(validOUs []string) martini.Handler {
	return verifyOUsHandler(validOUs, Verify)
}

func verifyOUsHandler(validOUs []string, verify verifyOUsFunc) martini.Handler {
	return func(res nethttp.ResponseWriter, req *nethttp.Request, c martini.Context) {
		log.Debug("Verifying client OU")
		if err := verify(req, validOUs); err != nil {
			nethttp.Error(res, err.Error(), nethttp.StatusUnauthorized)
		}
	}
}

// AppendKeyPair loads the given TLS key pair and appends it to
// tlsConfig.Certificates.
func AppendKeyPair(tlsConfig *tls.Config, certFile string, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	return nil
}

// Read in a keypair where the key is password protected
func AppendKeyPairWithPassword(tlsConfig *tls.Config, certFile string, keyFile string, pemPass []byte) error {

	// Certificates aren't usually password protected, but we're kicking the password
	// along just in case.  It won't be used if the file isn't encrypted
	certData, err := ReadPEMData(certFile, pemPass)
	if err != nil {
		return err
	}
	keyData, err := ReadPEMData(keyFile, pemPass)
	if err != nil {
		return err
	}
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return err
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	return nil
}

// Read a PEM file and ask for a password to decrypt it if needed
func ReadPEMData(pemFile string, pemPass []byte) ([]byte, error) {
	pemData, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return pemData, err
	}

	pemBlock, err := decodePEMBlock(pemData, pemFile)
	if err != nil {
		return pemData, err
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		// Decrypt and get the ASN.1 DER bytes here
		pemData, err = x509.DecryptPEMBlock(pemBlock, pemPass)
		if err != nil {
			return pemData, err
		}
		log.Info("Decrypted", pemFile, "successfully")
		// Shove the decrypted DER bytes into a new pem Block with blank headers
		newBlock := pem.Block{
			Type:  pemBlock.Type,
			Bytes: pemData,
		}
		// This is now like reading in an uncrypted key from a file and stuffing it
		// into a byte stream
		pemData = pem.EncodeToMemory(&newBlock)
	}
	return pemData, nil
}

// Print a password prompt on the terminal and collect a password
func GetPEMPassword(pemFile string) []byte {
	fmt.Printf("Password for %s: ", pemFile)
	pass, err := gopass.GetPasswd()
	if err != nil {
		// We'll error with an incorrect password at DecryptPEMBlock
		return []byte("")
	}
	return pass
}

// Determine if PEM file is encrypted
func IsEncryptedPEM(pemFile string) bool {
	pemData, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return false
	}
	pemBlock, err := decodePEMBlock(pemData, pemFile)
	if err != nil {
		return false
	}
	if len(pemBlock.Bytes) == 0 {
		return false
	}
	return x509.IsEncryptedPEMBlock(pemBlock)
}

func decodePEMBlock(pemData []byte, pemFile string) (*pem.Block, error) {
	// We should really just get the pem.Block back here, if there's other
	// junk on the end, warn about it.
	pemBlock, rest := pem.Decode(pemData)
	if pemBlock == nil {
		return nil, fmt.Errorf("no PEM block found in %s", pemFile)
	}
	if len(rest) > 0 {
		log.Warning("Didn't parse all of", pemFile)
	}
	return pemBlock, nil
}

type tlsListenFunc func(network, addr string, config *tls.Config) (net.Listener, error)
type httpServeFunc func(l net.Listener, handler nethttp.Handler) error

// ListenAndServeTLS acts identically to http.ListenAndServeTLS, except that it
// expects TLS configuration.
func ListenAndServeTLS(addr string, handler nethttp.Handler, tlsConfig *tls.Config) error {
	return listenAndServeTLS(addr, handler, tlsConfig, tls.Listen, nethttp.Serve)
}

func listenAndServeTLS(addr string, handler nethttp.Handler, tlsConfig *tls.Config, listen tlsListenFunc, serve httpServeFunc) error {
	if addr == "" {
		// On unix Listen calls getaddrinfo to parse the port, so named ports are fine as long
		// as they exist in /etc/services
		addr = ":https"
	}
	l, err := listen("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	return serve(l, handler)
}
