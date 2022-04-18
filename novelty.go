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

const BasePath = "api/v1/novelty"

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

type Component struct {
	Index   int     `json:"index"`
	Value   string  `json:"value"`
	Novelty float32 `json:"novelty"`
}

type NoveltyClient struct {
	serverURL string
	client    http.Client
}

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

	var nResp NoveltyResponse
	err = json.Unmarshal(body, &nResp)
	if err != nil {
		return nil, err
	}
	log.Printf("novelty response: %+v", string(body))
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
	req.Header.Add("cookie", cookie)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
