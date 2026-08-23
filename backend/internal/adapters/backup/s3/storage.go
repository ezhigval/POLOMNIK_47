package s3

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"polomnik/internal/ports"
)

type Storage struct {
	endpoint  string
	region    string
	bucket    string
	accessKey string
	secretKey string
	prefix    string
	client    *http.Client
}

func New(endpoint, region, bucket, accessKey, secretKey, prefix string, timeout time.Duration) *Storage {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Storage{
		endpoint:  strings.TrimRight(strings.TrimSpace(endpoint), "/"),
		region:    strings.TrimSpace(region),
		bucket:    strings.TrimSpace(bucket),
		accessKey: strings.TrimSpace(accessKey),
		secretKey: strings.TrimSpace(secretKey),
		prefix:    strings.Trim(strings.TrimSpace(prefix), "/"),
		client:    &http.Client{Timeout: timeout},
	}
}

var _ ports.BackupStoragePort = (*Storage)(nil)

func (s *Storage) Configured() bool {
	return s != nil && s.endpoint != "" && s.bucket != "" && s.accessKey != "" && s.secretKey != ""
}

func (s *Storage) Upload(ctx context.Context, object ports.BackupObject) error {
	if !s.Configured() {
		return ports.ErrBackupStorageNotConfigured
	}
	name := strings.TrimLeft(object.Name, "/")
	if name == "" {
		return fmt.Errorf("backup object name is required")
	}
	if s.prefix != "" {
		name = path.Join(s.prefix, name)
	}

	payload, err := io.ReadAll(object.Body)
	if err != nil {
		return fmt.Errorf("read backup object: %w", err)
	}

	objectURL, err := url.Parse(s.endpoint + "/" + s.bucket + "/" + name)
	if err != nil {
		return err
	}

	contentType := object.ContentType
	if contentType == "" {
		contentType = "application/gzip"
	}

	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	payloadHash := sha256Hex(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, objectURL.String(), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", objectURL.Host)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("Authorization", s.authorization(req, payloadHash, amzDate, dateStamp))

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("s3 upload failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s *Storage) authorization(req *http.Request, payloadHash, amzDate, dateStamp string) string {
	region := s.region
	if region == "" {
		region = "ru-central1"
	}
	canonicalHeaders := "content-type:" + req.Header.Get("Content-Type") + "\n" +
		"host:" + req.URL.Host + "\n" +
		"x-amz-content-sha256:" + payloadHash + "\n" +
		"x-amz-date:" + amzDate + "\n"
	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.EscapedPath(),
		"",
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	scope := dateStamp + "/" + region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	signingKey := hmacSHA256([]byte("AWS4"+s.secretKey), []byte(dateStamp))
	signingKey = hmacSHA256(signingKey, []byte(region))
	signingKey = hmacSHA256(signingKey, []byte("s3"))
	signingKey = hmacSHA256(signingKey, []byte("aws4_request"))
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	return fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey,
		scope,
		signedHeaders,
		signature,
	)
}

func sha256Hex(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}
