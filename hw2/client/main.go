package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const serverURL = "http://localhost:8082"

type VersionResponse struct {
	Version string `json:"version"`
}

type DecodeRequest struct {
	InputString string `json:"inputString"`
}

type DecodeResponse struct {
	OutputString string `json:"outputString"`
}

type HardOpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	version, err := getVersion()
	if err != nil {
		log.Fatalf("Error getting version: %v", err)
	}

	testString := "Hello, World!"
	encodedString := base64.StdEncoding.EncodeToString([]byte(testString))

	decodedString, err := decodeString(encodedString)
	if err != nil {
		log.Fatalf("Error decoding string: %v", err)
	}

	success, statusCode, err := hardOp()

	if err != nil {
		fmt.Printf("%s\t%s\t%v\n", version, decodedString, err)
	} else {
		fmt.Printf("%s\t%s\t%v,%d\n", version, decodedString, success, statusCode)
	}
}

func getVersion() (string, error) {
	resp, err := http.Get(serverURL + "/version")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return "", err
	}

	return versionResp.Version, nil
}

func decodeString(encodedString string) (string, error) {
	reqBody := DecodeRequest{
		InputString: encodedString,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(serverURL+"/decode", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var decodeResp DecodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodeResp); err != nil {
		return "", err
	}

	return decodeResp.OutputString, nil
}

func hardOp() (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/hard-op", nil)
	if err != nil {
		return false, 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, 0, fmt.Errorf("timeout exceeded")
		}
		return false, 0, err
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	success := resp.StatusCode < 500
	return success, resp.StatusCode, nil
}
