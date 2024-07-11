package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// RSA, return order: privateKey, publicKey, error
func GenerateKeyPair(keySize int) (string, string, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return "", "", err
	}

	// Get public key from private key
	publicKey := &privateKey.PublicKey

	// Encode private key to PEM
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyPEMString := string(pem.EncodeToMemory(privateKeyPEM))

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	publicKeyPEMString := string(pem.EncodeToMemory(publicKeyPEM))

	return privateKeyPEMString, publicKeyPEMString, nil
}

func EncryptData(data string, publicKeyPEM string) (string, error) {
	// Decode public key
	pemBlock, _ := pem.Decode([]byte(publicKeyPEM))
	pub, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if err != nil {
		return "", err
	}

	// Encrypt data
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		return "", err
	}

	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)

	return encryptedBase64, nil
}

func DecryptData(data string, privateKeyPEM string) (string, error) {
	// Decode private key
	pemBlock, _ := pem.Decode([]byte(privateKeyPEM))
	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}
