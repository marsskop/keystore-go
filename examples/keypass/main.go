package main

import (
	"crypto/x509"
	"log"
	"os"

	"github.com/pavel-v-chernykh/keystore-go/v4"
)

func readKeyStore(filename string, password []byte) keystore.KeyStore {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ks := keystore.New()
	if err := ks.Load(f, password); err != nil {
		log.Fatal(err) // nolint: gocritic
	}

	return ks
}

func zeroing(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

// nolint: godot
// keytool -genkeypair -alias alias -storepass password -keypass keypassword -keyalg RSA -keystore keystore.jks
func main() {
	password := []byte{'p', 'a', 's', 's', 'w', 'o', 'r', 'd'}
	defer zeroing(password)

	keyPassword := []byte{'k', 'e', 'y', 'p', 'a', 's', 's', 'w', 'o', 'r', 'd'}
	defer zeroing(keyPassword)

	ks := readKeyStore("keystore.jks", password)

	pke, err := ks.GetPrivateKeyEntry("alias", keyPassword)
	if err != nil {
		log.Fatal(err) // nolint: gocritic
	}

	key, err := x509.ParsePKCS8PrivateKey(pke.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", key)
}
