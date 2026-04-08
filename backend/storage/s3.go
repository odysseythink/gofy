package storage

// TODO: Implement S3 storage backend once github.com/aws/aws-sdk-go-v2 is added to vendor.
//
// This file will provide S3Storage implementing StorageProvider, backed by AWS S3
// (or any S3-compatible object store such as MinIO).
//
// Required dependencies (not yet vendored):
//   github.com/aws/aws-sdk-go-v2/aws
//   github.com/aws/aws-sdk-go-v2/config
//   github.com/aws/aws-sdk-go-v2/credentials
//   github.com/aws/aws-sdk-go-v2/service/s3
//
// Configuration struct:
//
//   type S3Config struct {
//       Bucket       string
//       AccessKey    string
//       SecretKey    string
//       Region       string
//       Endpoint     string
//       UsePathStyle bool
//   }
//
// Example initialization (to be uncommented when SDK is available):
//
// import (
//     "bytes"
//     "context"
//     "fmt"
//     "io"
//     "os"
//
//     "github.com/aws/aws-sdk-go-v2/aws"
//     "github.com/aws/aws-sdk-go-v2/config"
//     "github.com/aws/aws-sdk-go-v2/credentials"
//     "github.com/aws/aws-sdk-go-v2/service/s3"
// )
//
// type S3Storage struct {
//     client *s3.Client
//     bucket string
// }
//
// func NewS3Storage(cfg S3Config) (*S3Storage, error) {
//     opts := []func(*config.LoadOptions) error{
//         config.WithRegion(cfg.Region),
//     }
//     if cfg.AccessKey != "" && cfg.SecretKey != "" {
//         opts = append(opts, config.WithCredentialsProvider(
//             credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
//         ))
//     }
//     awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
//     if err != nil {
//         return nil, fmt.Errorf("failed to load AWS config: %w", err)
//     }
//     clientOpts := []func(*s3.Options){}
//     if cfg.Endpoint != "" {
//         clientOpts = append(clientOpts, func(o *s3.Options) {
//             o.BaseEndpoint = aws.String(cfg.Endpoint)
//         })
//     }
//     if cfg.UsePathStyle {
//         clientOpts = append(clientOpts, func(o *s3.Options) {
//             o.UsePathStyle = true
//         })
//     }
//     client := s3.NewFromConfig(awsCfg, clientOpts...)
//     return &S3Storage{client: client, bucket: cfg.Bucket}, nil
// }
//
// func (s *S3Storage) Save(filename string, data []byte) error {
//     _, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
//         Bucket: aws.String(s.bucket),
//         Key:    aws.String(filename),
//         Body:   bytes.NewReader(data),
//     })
//     return err
// }
//
// func (s *S3Storage) LoadOnce(filename string) ([]byte, error) {
//     resp, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
//         Bucket: aws.String(s.bucket),
//         Key:    aws.String(filename),
//     })
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()
//     return io.ReadAll(resp.Body)
// }
//
// func (s *S3Storage) LoadStream(filename string) (io.ReadCloser, error) {
//     resp, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
//         Bucket: aws.String(s.bucket),
//         Key:    aws.String(filename),
//     })
//     if err != nil {
//         return nil, err
//     }
//     return resp.Body, nil
// }
//
// func (s *S3Storage) Download(filename string, targetPath string) error {
//     data, err := s.LoadOnce(filename)
//     if err != nil {
//         return err
//     }
//     return os.WriteFile(targetPath, data, 0644)
// }
//
// func (s *S3Storage) Exists(filename string) (bool, error) {
//     _, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
//         Bucket: aws.String(s.bucket),
//         Key:    aws.String(filename),
//     })
//     if err != nil {
//         return false, nil
//     }
//     return true, nil
// }
//
// func (s *S3Storage) Delete(filename string) error {
//     _, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
//         Bucket: aws.String(s.bucket),
//         Key:    aws.String(filename),
//     })
//     return err
// }
