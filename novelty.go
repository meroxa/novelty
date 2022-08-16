package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// BasePath for the Novelty REST API
const BasePath = "api/v1/novelty"

// NoveltyResponse models the response from the Novelty Server API
type NoveltyResponse struct {
	Observation        interface{} `json:"observation"`
	Score              float32     `json:"score"`
	TotalObsScore      float32     `json:"totalObsScore"`
	Sequence           int         `json:"sequence"`
	Probability        float32     `json:"probability"`
	Uniqueness         float32     `json:"uniqueness"`
	InfoContent        float32     `json:"infoContent"`
	MostNovelComponent Component   `json:"mostNovelComponent"`
}

// Component ...
type Component struct {
	Index   int     `json:"index"`
	Value   string  `json:"value"`
	Novelty float32 `json:"novelty"`
}

// NoveltyClient wraps http.Client for convenience
type NoveltyClient struct {
	serverURL string
	client    http.Client
}

// NewNoveltyClient constructor
func NewNoveltyClient(serverURL string) (NoveltyClient, error) {
	// validate serverURL
	_, err := url.Parse(serverURL)
	if err != nil {
		return NoveltyClient{}, err
	}
	c := http.Client{
		Timeout: 3 * time.Second,
	}
	return NoveltyClient{
		serverURL: serverURL,
		client:    c,
	}, nil
}

// Observe calls the Novelty Server API submitting an "observation". The server
// synchronously returns a NoveltyResponse with the novelty details embedded.
func (c NoveltyClient) Observe(name string, data []string) (*NoveltyResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/observe", c.serverURL, BasePath, name)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body, err := c.post(path, jsonData)
	if err != nil {
		return nil, err
	}
	log.Printf("novelty response: %+v", string(body))

	var nResp NoveltyResponse
	err = json.Unmarshal(body, &nResp)
	if err != nil {
		return nil, err
	}
	return &nResp, nil
}

func (c NoveltyClient) post(path string, data []byte) ([]byte, error) {
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	cookie := os.Getenv("NOVELTY_AUTH")
	if cookie != "" {
		req.Header.Add("cookie", cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
