package crypt

import (
	"crypto/rsa"
	"crypto/sha1"
	"hash"
	"io"
)

type PKCS1OAepCipher struct {
	key      *rsa.PrivateKey
	hashAlgo func() hash.Hash
	mgf      func(seed []byte, length int) []byte
	label    []byte
	rand     io.Reader
}

func NewPKCS1OAepCipher(key *rsa.PrivateKey, hashAlgo func() hash.Hash, mgf func([]byte, int) []byte, label []byte, rand io.Reader) *PKCS1OAepCipher {
	if hashAlgo == nil {
		hashAlgo = sha1.New
	}
	if mgf == nil {
		mgf = defaultMGF1SHA1
	}
	return &PKCS1OAepCipher{key: key, hashAlgo: hashAlgo, mgf: mgf, label: label, rand: rand}
}

func (c *PKCS1OAepCipher) CanEncrypt() bool {
	// Simplified check for encryption capability
	return c.key.PublicKey.E != 0 && c.key.PublicKey.N != nil
}

func (c *PKCS1OAepCipher) CanDecrypt() bool {
	// Simplified check for decryption capability
	return c.key.D != nil
}

func (c *PKCS1OAepCipher) Encrypt(message []byte) ([]byte, error) {
	// Implement the encryption logic here based on RFC3447
	// This is a simplified example and does not include all steps from the Python code.
	return rsa.EncryptOAEP(c.hashAlgo(), c.rand, &c.key.PublicKey, message, c.label)
}

func (c *PKCS1OAepCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	// Implement the decryption logic here based on RFC3447
	return rsa.DecryptOAEP(c.hashAlgo(), c.rand, c.key, ciphertext, c.label)
}

// Default mask generation function for MGF1 with SHA-1
func defaultMGF1SHA1(seed []byte, length int) []byte {
	// Implementation of MGF1 using SHA-1
	// Note: Go's standard library does not provide a direct way to do this,
	// so you would need to implement it or use a third-party package.
	return nil
}

func ObfuscatedToken(token string) string {
	if token == "" {
		return token
	}
	if len(token) <= 8 {
		return "********************"
	}
	return token[:6] + "************" + token[len(token)-2:]
}

// func DecryptToken(tenantID string, token string) (string, error) {
// 	decodedToken, err := base64.StdEncoding.DecodeString(token)
// 	if err != nil {
// 		return "", err
// 	}

// 	return decrypt(decodedToken, tenantID)
// }
// func decrypt(encryptedText []byte, tenantID string) (string, error) {
// 	rsaKey, cipherRSA, err := getDecryptDecoding(tenantID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return decryptTokenWithDecoding(encryptedText, rsaKey, cipherRSA)
// }

// func getDecryptDecoding(tenantID string) (*rsa.PrivateKey, *gmpy2.PKCS1OAEP, error) {
// 	filepath := fmt.Sprintf("privkeys/%s/private.pem", tenantID)

// 	hash := sha3.New256()
// 	hash.Write([]byte(filepath))
// 	hashBytes := hash.Sum(nil)
// 	cacheKey := fmt.Sprintf("tenant_privkey:%x", hashBytes)

// 	privateKeyBytes := cache.Instance().GetString(cacheKey)
// 	if len(privateKeyBytes) == 0 {
// 		privateKeyBytes, err := os.ReadFile(filepath)
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				return nil, nil, fmt.Errorf("private key not found, tenant_id: %s", tenantID)
// 			}
// 			return nil, nil, err
// 		}

// 		cache.Instance().SetEx(cacheKey, privateKeyBytes, 120)
// 	}

// 	rsaKey, err := rsa.VerifyPSS().ImportKeyFromBytes(privateKeyBytes)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	cipherRSA := gmpy2.NewPKCS1OAEP(rsaKey)

// 	return rsaKey, cipherRSA, nil
// }

// func decryptTokenWithDecoding(encryptedText []byte, rsaKey *rsa.PrivateKey, cipherRSA *gmpy2.PKCS1OAEP) (string, error) {
// 	if len(encryptedText) == 0 {
// 		return "", errors.New("encrypted text is empty")
// 	}

// 	if string(encryptedText[:len(prefixHybrid)]) == prefixHybrid {
// 		encryptedText = encryptedText[len(prefixHybrid):]

// 		keySize := rsaKey.Size()
// 		encAESKey := encryptedText[:keySize]
// 		nonce := encryptedText[keySize : keySize+16]
// 		tag := encryptedText[keySize+16 : keySize+32]
// 		ciphertext := encryptedText[keySize+32:]

// 		aesKey, err := cipherRSA.Decrypt(encAESKey)
// 		if err != nil {
// 			return "", err
// 		}

// 		cipherAES, err := aes.NewCipher(aesKey)
// 		if err != nil {
// 			return "", err
// 		}

// 		cipherEAX, err := aes.NewEAX(cipherAES, nonce)
// 		if err != nil {
// 			return "", err
// 		}

// 		decryptedText, err := cipherEAX.Open(nil, ciphertext, tag)
// 		if err != nil {
// 			return "", err
// 		}

// 		return string(decryptedText), nil
// 	} else {
// 		decryptedBytes, err := cipherRSA.Decrypt(encryptedText)
// 		if err != nil {
// 			return "", err
// 		}
// 		return string(decryptedBytes), nil
// 	}
// }
