package cryptutils

import (
	"errors"
)

var (
	ErrPrivkeyNotFound = errors.New("private key not found")
)

// def get_decrypt_decoding(tenant_id):
//     filepath = "privkeys/{tenant_id}".format(tenant_id=tenant_id) + "/private.pem"

//     cache_key = "tenant_privkey:{hash}".format(hash=hashlib.sha3_256(filepath.encode()).hexdigest())
//     private_key = redis_client.get(cache_key)
//     if not private_key:
//         try:
//             private_key = storage.load(filepath)
//         except FileNotFoundError:
//             raise PrivkeyNotFoundError("Private key not found, tenant_id: {tenant_id}".format(tenant_id=tenant_id))

//         redis_client.setex(cache_key, 120, private_key)

//     rsa_key = RSA.import_key(private_key)
//     cipher_rsa = gmpy2_pkcs10aep_cipher.new(rsa_key)

//     return rsa_key, cipher_rsa
