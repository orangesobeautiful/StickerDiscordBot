package config

import (
	"encoding/base64"
	"errors"
)

type sessionKeyInfo struct {
	authentication    []byte
	AuthenticationB64 string
	encryption        []byte
	EncryptionB64     string
}

var (
	authenticationFormatError = errors.New("authentication key must be 32 or 64 bytes base64 encoded string")
	encryptionFormatError     = errors.New("encryption key must be 16, 24, or 32 bytes base64 encoded string")
)

func (s *sessionKeyInfo) Preprocessing() error {
	var err error
	s.authentication, err = base64.StdEncoding.DecodeString(s.AuthenticationB64)
	if err != nil {
		return authenticationFormatError
	}
	switch len(s.authentication) {
	case 32, 64:
	default:
		return authenticationFormatError
	}

	s.encryption, err = base64.StdEncoding.DecodeString(s.EncryptionB64)
	if err != nil {
		return encryptionFormatError
	}
	switch len(s.encryption) {
	case 0, 16, 24, 32:
	default:
		return encryptionFormatError
	}

	return nil
}

func (s *sessionKeyInfo) Authentication() []byte {
	return s.authentication
}

func (s *sessionKeyInfo) Encryption() []byte {
	return s.encryption
}

func (s *sessionKeyInfo) SessionKeyPair() [][]byte {
	keyPair := [][]byte{s.authentication}

	if len(s.encryption) != 0 {
		keyPair = append(keyPair, s.encryption)
	}
	return keyPair
}
