package srv_sec

import (
	"srv/srv_conf"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"

	"log"
	"time"
	"os"
)

var CertFile string
var KeyFile string
var JwtSecret string

// GetSecret returns the JWT secret
func Env_SetSecret() {

	if srv_conf.IsGinModDebug() {
		JwtSecret = srv_conf.GetString("jwt_secret")
	} else {

		secret, err := generateSecret()
		if err != nil {
			log.Fatal("Error generating JWT secret")
		}
		JwtSecret = secret
	}
}

func generateSecret() (string, error) {
	// Since we want a 64-character secret and each character is 8 bits,
	// we need to generate 32 bytes and then encode it using base64
	const byteLength = 32

	secretBytes := make([]byte, byteLength)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}

	// Encoding the random bytes to base64
	secretBase64 := base64.RawURLEncoding.EncodeToString(secretBytes)
	return secretBase64, nil
}

func GenerateTLS(keySize int) error {

	log.Printf("Generating %d bits TLS keys", keySize)

	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	keyOut, err := os.OpenFile(KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return err
	}

	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}

	if err := pem.Encode(keyOut, keyBlock); err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"raadig"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365), // one year

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return err
	}

	certOut, err := os.OpenFile(CertFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer certOut.Close()

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	if err := pem.Encode(certOut, certBlock); err != nil {
		return err
	}

	return nil
}
