package env

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Down stops all resources in the specified environment(s).
func Down(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/env/down", host)
	payload := bytes.NewBuffer([]byte(json))
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response: ", err.Error())
		return "", err
	}

	return string(b), nil
}

// Show return a string with environments details.
func Show(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/env/show", host)
	payload := bytes.NewBuffer([]byte(json))
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response: ", err.Error())
		return "", err
	}

	return string(b), nil
}

// Up starts all resources in the specified environment(s).
func Up(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/env/up", host)
	payload := bytes.NewBuffer([]byte(json))
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response: ", err.Error())
		return "", err
	}

	return string(b), nil
}
