package storage

// TODO: Implement Aliyun OSS storage backend once github.com/aliyun/aliyun-oss-go-sdk is added to vendor.
//
// This file will provide AliyunOSSStorage implementing StorageProvider, backed by
// Alibaba Cloud Object Storage Service.
//
// Required dependencies (not yet vendored):
//   github.com/aliyun/aliyun-oss-go-sdk/oss
//
// Configuration struct:
//
//   type AliyunOSSConfig struct {
//       Bucket    string
//       AccessKey string
//       SecretKey string
//       Endpoint  string
//       Region    string
//       Folder    string
//   }
//
// Example implementation (to be uncommented when SDK is available):
//
// import (
//     "bytes"
//     "fmt"
//     "io"
//     "path"
//
//     "github.com/aliyun/aliyun-oss-go-sdk/oss"
// )
//
// type AliyunOSSStorage struct {
//     bucket *oss.Bucket
//     folder string
// }
//
// func NewAliyunOSSStorage(cfg AliyunOSSConfig) (*AliyunOSSStorage, error) {
//     client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
//     if err != nil {
//         return nil, fmt.Errorf("failed to create OSS client: %w", err)
//     }
//     bucket, err := client.Bucket(cfg.Bucket)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get OSS bucket: %w", err)
//     }
//     return &AliyunOSSStorage{bucket: bucket, folder: cfg.Folder}, nil
// }
//
// func (s *AliyunOSSStorage) key(filename string) string {
//     if s.folder != "" {
//         return path.Join(s.folder, filename)
//     }
//     return filename
// }
//
// func (s *AliyunOSSStorage) Save(filename string, data []byte) error {
//     return s.bucket.PutObject(s.key(filename), bytes.NewReader(data))
// }
//
// func (s *AliyunOSSStorage) LoadOnce(filename string) ([]byte, error) {
//     body, err := s.bucket.GetObject(s.key(filename))
//     if err != nil {
//         return nil, err
//     }
//     defer body.Close()
//     return io.ReadAll(body)
// }
//
// func (s *AliyunOSSStorage) LoadStream(filename string) (io.ReadCloser, error) {
//     return s.bucket.GetObject(s.key(filename))
// }
//
// func (s *AliyunOSSStorage) Download(filename string, targetPath string) error {
//     return s.bucket.GetObjectToFile(s.key(filename), targetPath)
// }
//
// func (s *AliyunOSSStorage) Exists(filename string) (bool, error) {
//     return s.bucket.IsObjectExist(s.key(filename))
// }
//
// func (s *AliyunOSSStorage) Delete(filename string) error {
//     return s.bucket.DeleteObject(s.key(filename))
// }
