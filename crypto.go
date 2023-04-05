package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
)

var ecc_session_key *ecdsa.PrivateKey

func crypto_init() (err error) {
	err = crypto_ecdh_generate_session_key()
	if err != nil {
		return
	}
	return
}

func crypto_ecdh_generate_session_key() (err error) {

	ecc_session_key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return
}

func crypto_ecc_get_session_key() *ecdsa.PrivateKey {
	return ecc_session_key
}

func crypto_ecdh_perform_ecdh(remote_public_key_x []byte, remote_public_key_y []byte) (secret [32]byte, err error) {

	/**/
	remote_public_key := new(ecdsa.PublicKey)

	XX := new(big.Int)
	XX.SetBytes(remote_public_key_x)
	YY := new(big.Int)
	YY.SetBytes(remote_public_key_y)

	remote_public_key.X = XX
	remote_public_key.Y = YY
	remote_public_key.Curve = elliptic.P256()
	/**/

	secret_int, _ := remote_public_key.Curve.ScalarMult(remote_public_key.X, remote_public_key.Y, ecc_session_key.D.Bytes())

	secret = sha256.Sum256(secret_int.Bytes())

	return
}

/* aes */

func crypto_aes_cbc_encrypt(key, plaintext []byte, iv []byte) (ciphertext []byte, err error) {
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	/* padding */
	bytes_to_add := aes.BlockSize - (len(plaintext) % 16)

	for i := 0; i < bytes_to_add; i++ {
		plaintext = append(plaintext, byte(bytes_to_add))
	}
	/* */

	ciphertext = make([]byte, len(plaintext))

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[:], plaintext)

	return
}
func crypto_aes_cbc_decrypt(key, ciphertext, iv []byte) (plaintext []byte, err error) {
	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		fmt.Printf("ciphertext too short")
		return
	}

	ciphertext = ciphertext[:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	/* padding size */
	padding_size := ciphertext[len(ciphertext)-1]
	/**/
	plaintext = ciphertext[:(len(ciphertext) - int(padding_size))]

	return
}

func crypto_aes_encrypt(data []byte, key []byte) (e_data []byte, iv []byte, err error) {
	iv = make([]byte, 16)
	_, err = rand.Read(iv)
	if err != nil {
		return
	}
	e_data, err = crypto_aes_cbc_encrypt(key, data, iv)
	return
}

func crypto_aes_decrypt(data []byte, key []byte, iv []byte) (p_data []byte, err error) {

	p_data, err = crypto_aes_cbc_decrypt(key, data, iv)

	return
}
