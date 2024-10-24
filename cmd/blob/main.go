package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	BlobAPIVersion = "7" // ref: https://github.com/vercel/storage/blob/main/packages/blob/src/api.ts const BLOB_API_VERSION
	DefaultBaseURL = "https://blob.vercel-storage.com"
)

func main() {
	ctx := context.Background()

	body := bytes.NewReader([]byte("test"))

	pathname := uuid.NewString() + ".txt"

	put, err := Put(ctx, pathname, body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("response: %+v\n", put)

	head, err := Head(ctx, put.URL)
	if err != nil {
		panic(err)
	}

	fmt.Printf("head: %+v\n", head)

	download, err := Download(ctx, put.DownloadURL)
	if err != nil {
		panic(err)
	}

	fmt.Printf("download: %s\n", download)

	if err := Del(ctx, put.URL); err != nil {
		panic(err)
	}
}

type PutResponse struct {
	URL                string `json:"url"`
	DownloadURL        string `json:"downloadUrl"`
	Pathname           string `json:"pathname"`
	ContentDisposition string `json:"contentDisposition"`
}

func Put(
	ctx context.Context,
	pathname string,
	body io.Reader,
) (*PutResponse, error) {
	if pathname == "" {
		return nil, fmt.Errorf("pathname is required")
	}

	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BLOB_READ_WRITE_TOKEN is required")
	}

	base, _ := url.Parse(DefaultBaseURL)

	base.Path = pathname

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, base.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-version", BlobAPIVersion)

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	var putRes PutResponse

	if err := json.NewDecoder(res.Body).Decode(&putRes); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return &putRes, nil
}

type HeadResponse struct {
	HasMore bool       `json:"hasMore"`
	Blobs   []HeadBlob `json:"blobs"`
}

type HeadBlob struct {
	URL         string    `json:"url"`
	DownloadURL string    `json:"downloadUrl"`
	Pathname    string    `json:"pathname"`
	Size        int       `json:"size"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

func Head(
	ctx context.Context,
	url string,
) (*HeadResponse, error) {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BLOB_READ_WRITE_TOKEN is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DefaultBaseURL+"?"+url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-version", BlobAPIVersion)

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	var headRes HeadResponse

	if err := json.NewDecoder(res.Body).Decode(&headRes); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return &headRes, nil
}

func Download(
	ctx context.Context,
	url string,
) ([]byte, error) {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BLOB_READ_WRITE_TOKEN is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-version", BlobAPIVersion)

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	return io.ReadAll(res.Body)
}

type DelRequest struct {
	URLs []string `json:"urls"`
}

func Del(
	ctx context.Context,
	url string,
) error {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return fmt.Errorf("BLOB_READ_WRITE_TOKEN is required")
	}

	body := DelRequest{
		URLs: []string{url},
	}

	buf := bytes.Buffer{}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, DefaultBaseURL+"/delete", &buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-version", BlobAPIVersion)

	req.Header.Set("Authorization", "Bearer "+token)

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	return nil
}
