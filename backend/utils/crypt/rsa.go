package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha3"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/mlog"
)

func GetPrivateKey(tenant_id string) string {
	file_path := filepath.Join("privkeys", tenant_id, "private.pem")
	cache_key := "tenant_privkey:" + fmt.Sprintf("%x", sha3.Sum256([]byte(file_path)))

	private_key := cache.Instance().GetString(cache_key)
	if private_key == "" {
		bindata, err := os.ReadFile(file_path)
		if err != nil {
			mlog.Errorf("read file=%s failed:%v", file_path, err)
			panic(exceptions.NewPrivkeyNotFoundError("Private key not found, tenant_id: " + tenant_id))
		}
		private_key = string(bindata)
		cache.Instance().SetEx(cache_key, private_key, 120*time.Second)
	}
	return private_key
}

func RSADecrypt(encrypted_text []byte, tenant_id string) []byte {
	private_key := GetPrivateKey(tenant_id)
	//pem解码
	block, _ := pem.Decode([]byte(private_key))
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted_text)
	//返回明文
	return plainText
}

func RSAEncrypt(plain_text []byte, pub_key string) []byte {
	//pem解码
	block, _ := pem.Decode([]byte(pub_key))
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plain_text)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

func GenerateKeyPair(tenant_id string) string {
	// Generates private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		mlog.Errorf("rsa generate key failed:%v", err)
		return ""
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey := pem.EncodeToMemory(block)

	// Generates public key from private key.
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return ""
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey := pem.EncodeToMemory(block)
	os.WriteFile(filepath.Join("privkeys", tenant_id, "private.pem"), prvkey, 0644)
	return string(pubkey)
}
