package file

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

// def get_signed_file_url(upload_file_id: str) -> str:
//     url = f"{gofy_config.FILES_URL}/files/{upload_file_id}/file-preview"

//     timestamp = str(int(time.time()))
//     nonce = os.urandom(16).hex()
//     key = gofy_config.SECRET_KEY.encode()
//     msg = f"file-preview|{upload_file_id}|{timestamp}|{nonce}"
//     sign = hmac.new(key, msg.encode(), hashlib.sha256).digest()
//     encoded_sign = base64.urlsafe_b64encode(sign).decode()

//     return f"{url}?timestamp={timestamp}&nonce={nonce}&sign={encoded_sign}"

// getSignedFileURL generates a signed URL for a file.
func GetSignedFileURL(uploadFileID string) string {
	baseURL := fmt.Sprintf("%s/files/%s/file-preview", confy.Get[string]("FILES_URL"), uploadFileID)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := generateNonce()
	sign := generateSignature("file-preview", uploadFileID, timestamp, nonce)

	query := url.Values{}
	query.Add("timestamp", timestamp)
	query.Add("nonce", nonce)
	query.Add("sign", sign)

	return fmt.Sprintf("%s?%s", baseURL, query.Encode())
}

// def verify_image_signature(*, upload_file_id: str, timestamp: str, nonce: str, sign: str) -> bool:
//     data_to_sign = f"image-preview|{upload_file_id}|{timestamp}|{nonce}"
//     secret_key = gofy_config.SECRET_KEY.encode()
//     recalculated_sign = hmac.new(secret_key, data_to_sign.encode(), hashlib.sha256).digest()
//     recalculated_encoded_sign = base64.urlsafe_b64encode(recalculated_sign).decode()

//     # verify signature
//     if sign != recalculated_encoded_sign:
//         return False

//     current_time = int(time.time())
//     return current_time - int(timestamp) <= gofy_config.FILES_ACCESS_TIMEOUT

// verifyImageSignature verifies the signature for an image.
func VerifyImageSignature(uploadFileID, timestamp, nonce, sign string) bool {
	return verifySignature("image-preview", uploadFileID, timestamp, nonce, sign)
}

// def verify_file_signature(*, upload_file_id: str, timestamp: str, nonce: str, sign: str) -> bool:
//     data_to_sign = f"file-preview|{upload_file_id}|{timestamp}|{nonce}"
//     secret_key = gofy_config.SECRET_KEY.encode()
//     recalculated_sign = hmac.new(secret_key, data_to_sign.encode(), hashlib.sha256).digest()
//     recalculated_encoded_sign = base64.urlsafe_b64encode(recalculated_sign).decode()

//     # verify signature
//     if sign != recalculated_encoded_sign:
//         return False

//	current_time = int(time.time())
//	return current_time - int(timestamp) <= gofy_config.FILES_ACCESS_TIMEOUT
//
// verifyFileSignature verifies the signature for a file.
func VerifyFileSignature(uploadFileID, timestamp, nonce, sign string) bool {
	return verifySignature("file-preview", uploadFileID, timestamp, nonce, sign)
}

// verifySignature verifies the signature based on the provided parameters.
func verifySignature(action, uploadFileID, timestamp, nonce, sign string) bool {
	// dataToSign := fmt.Sprintf("%s|%s|%s|%s", action, uploadFileID, timestamp, nonce)
	recalculatedSign := generateSignature(action, uploadFileID, timestamp, nonce)

	if sign != recalculatedSign {
		return false
	}

	currentTime := time.Now().Unix()
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	return currentTime-timestampInt <= confy.Get[int64]("FILES_ACCESS_TIMEOUT")
}

// generateSignature generates a signature for the given parameters.
func generateSignature(action, uploadFileID, timestamp, nonce string) string {
	dataToSign := fmt.Sprintf("%s|%s|%s|%s", action, uploadFileID, timestamp, nonce)
	sign := hmac.New(sha256.New, []byte(confy.Get[string]("SECRET_KEY")))
	sign.Write([]byte(dataToSign))
	return base64.RawURLEncoding.EncodeToString(sign.Sum(nil))
}

// generateNonce generates a random nonce.
func generateNonce() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		mlog.Errorf("rand read failed:%v", err)
		return "fd631ee32b653634d92a76c6c07a117d"
	}
	return hex.EncodeToString(b)
}

// func main() {
// 	// Example usage
// 	uploadFileID := "example-file-id"
// 	signedURL := getSignedFileURL(uploadFileID)
// 	fmt.Println("Signed URL:", signedURL)

// 	// Example verification
// 	timestamp := "1680000000"
// 	nonce := "random-nonce"
// 	sign := "example-signature"
// 	isValid := verifyFileSignature(uploadFileID, timestamp, nonce, sign)
// 	fmt.Println("Signature valid:", isValid)
// }
