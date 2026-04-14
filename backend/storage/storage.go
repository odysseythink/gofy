package storage

import (
	"io"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

// StorageProvider defines the interface for all storage backends.
type StorageProvider interface {
	Save(filename string, data []byte) error
	LoadOnce(filename string) ([]byte, error)
	LoadStream(filename string) (io.ReadCloser, error)
	Download(filename string, targetPath string) error
	Exists(filename string) (bool, error)
	Delete(filename string) error
}

var defaultStorage StorageProvider

// Init initializes the storage backend based on config.
func Init() {
	storageType := confy.GetWithDefault[string]("storage.type", "local")
	var err error

	switch storageType {
	case "local":
		basePath := confy.GetWithDefault[string]("storage.local.path", "./storage")
		defaultStorage, err = NewLocalStorage(basePath)
	case "s3":
		// TODO: Implement S3 storage initialization once aws-sdk-go-v2 is vendored.
		// defaultStorage, err = NewS3Storage(S3Config{
		//     Bucket:       confy.GetWithDefault[string]("storage.s3.bucket", ""),
		//     AccessKey:    confy.GetWithDefault[string]("storage.s3.access_key", ""),
		//     SecretKey:    confy.GetWithDefault[string]("storage.s3.secret_key", ""),
		//     Region:       confy.GetWithDefault[string]("storage.s3.region", "us-east-1"),
		//     Endpoint:     confy.GetWithDefault[string]("storage.s3.endpoint", ""),
		//     UsePathStyle: confy.GetWithDefault[bool]("storage.s3.use_path_style", false),
		// })
		mlog.Errorf("s3 storage backend is not yet implemented, falling back to local")
		basePath := confy.GetWithDefault[string]("storage.local.path", "./storage")
		defaultStorage, err = NewLocalStorage(basePath)
	case "aliyun-oss":
		// TODO: Implement Aliyun OSS storage initialization once aliyun-oss-go-sdk is vendored.
		// defaultStorage, err = NewAliyunOSSStorage(AliyunOSSConfig{
		//     Bucket:    confy.GetWithDefault[string]("storage.aliyun_oss.bucket", ""),
		//     AccessKey: confy.GetWithDefault[string]("storage.aliyun_oss.access_key", ""),
		//     SecretKey: confy.GetWithDefault[string]("storage.aliyun_oss.secret_key", ""),
		//     Endpoint:  confy.GetWithDefault[string]("storage.aliyun_oss.endpoint", ""),
		//     Region:    confy.GetWithDefault[string]("storage.aliyun_oss.region", ""),
		//     Folder:    confy.GetWithDefault[string]("storage.aliyun_oss.folder", ""),
		// })
		mlog.Errorf("aliyun-oss storage backend is not yet implemented, falling back to local")
		basePath := confy.GetWithDefault[string]("storage.local.path", "./storage")
		defaultStorage, err = NewLocalStorage(basePath)
	default:
		mlog.Errorf("unsupported storage type: %s, falling back to local", storageType)
		basePath := confy.GetWithDefault[string]("storage.local.path", "./storage")
		defaultStorage, err = NewLocalStorage(basePath)
	}

	if err != nil {
		mlog.Errorf("failed to initialize storage: %v", err)
	}
}

// Default returns the default storage provider instance.
func Default() StorageProvider {
	if defaultStorage == nil {
		mlog.Errorf("storage not initialized, using local fallback")
		s, _ := NewLocalStorage("./storage")
		return s
	}
	return defaultStorage
}

// Save saves data to the default storage.
func Save(filename string, data []byte) error {
	return Default().Save(filename, data)
}

// LoadOnce loads entire file content from the default storage.
func LoadOnce(filename string) ([]byte, error) {
	return Default().LoadOnce(filename)
}

// LoadStream returns a streaming reader from the default storage.
func LoadStream(filename string) (io.ReadCloser, error) {
	return Default().LoadStream(filename)
}

// Exists checks if a file exists in the default storage.
func Exists(filename string) (bool, error) {
	return Default().Exists(filename)
}

// Delete removes a file from the default storage.
func Delete(filename string) error {
	return Default().Delete(filename)
}
