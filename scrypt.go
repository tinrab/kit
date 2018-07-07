package util

import (
	"golang.org/x/crypto/scrypt"
	"encoding/binary"
	"bytes"
	"encoding/base64"
)

// Scrypt computes scrypt key using data and cost parameters.
func Scrypt(data []byte, n, r, p, keyLength int) ([]byte, error) {
	salt := []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}
	hash, err := scrypt.Key(data, salt, n, r, p, keyLength)
	if err != nil {
		return nil, err
	}

	head := make([]byte, binary.MaxVarintLen64*4)
	s := binary.PutUvarint(head, uint64(n))
	s += binary.PutUvarint(head[s:], uint64(r))
	s += binary.PutUvarint(head[s:], uint64(p))
	s += binary.PutUvarint(head[s:], uint64(keyLength))
	hash = append(head[0:s], hash...)

	return hash, nil
}

// ScryptToBase64 computes scrypt key and returns the Base64 encoded result.
func ScryptToBase64(data []byte, n, r, p, keyLength int) (string, error) {
	hash, err := Scrypt(data, n, r, p, keyLength)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hash), nil
}

// ScryptEquals compares data with scrypt hash and returns true if they are equal.
func ScryptEquals(data, hash []byte) bool {
	buf := bytes.NewReader(hash)
	n, err := binary.ReadUvarint(buf)
	if err != nil {
		return false
	}
	r, err := binary.ReadUvarint(buf)
	if err != nil {
		return false
	}
	p, err := binary.ReadUvarint(buf)
	if err != nil {
		return false
	}
	keyLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return false
	}

	dataHash, err := Scrypt(data, int(n), int(r), int(p), int(keyLength))
	if err != nil {
		return false
	}

	diff := len(hash) ^ len(dataHash)
	for i := 0; i < len(hash) && i < len(dataHash); i++ {
		diff |= int(hash[i]) ^ int(dataHash[i])
	}
	return diff == 0
}

// ScryptEqualsBase64 compares data with scrypt hash encoded as Base64 string
func ScryptEqualsBase64(data []byte, hash string) bool {
	hashBytes, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	return ScryptEquals(data, hashBytes)
}
