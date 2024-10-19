package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const DefaultBaseURL = "https://edge-config.vercel.com"

func main() {
	ctx := context.Background()

	value, err := Get(ctx, "test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("value: %s\n", value)

	all, err := GetAll(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("all: %v\n", all)
}

func Get(
	ctx context.Context,
	key string,
) (string, error) {
	id := os.Getenv("EDGE_CONFIG_ID")
	if id == "" {
		return "", fmt.Errorf("EDGE_CONFIG_ID is required")
	}

	token := os.Getenv("EDGE_CONFIG_TOKEN")
	if token == "" {
		return "", fmt.Errorf("EDGE_CONFIG_TOKEN is required")
	}

	req, err := http.NewRequest(http.MethodGet, DefaultBaseURL+"/"+id+"/item/"+key+"?token="+token, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	var value string

	if err := json.NewDecoder(res.Body).Decode(&value); err != nil {
		return "", err
	}

	return value, nil
}

func GetAll(
	ctx context.Context,
) (map[string]string, error) {
	id := os.Getenv("EDGE_CONFIG_ID")
	if id == "" {
		return nil, fmt.Errorf("EDGE_CONFIG_ID is required")
	}

	token := os.Getenv("EDGE_CONFIG_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("EDGE_CONFIG_TOKEN is required")
	}

	req, err := http.NewRequest(http.MethodGet, DefaultBaseURL+"/"+id+"/items?token="+token, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			println(err)
		}
	}()

	println(res.StatusCode)

	var all map[string]string

	if err := json.NewDecoder(res.Body).Decode(&all); err != nil {
		return nil, err
	}

	return all, nil
}
