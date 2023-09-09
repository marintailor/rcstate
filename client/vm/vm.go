package vm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// List return a string with table formatted list of virtual machines.
func List(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/vm/list", host)
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

// Start will start a virtual machine.
func Start(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/vm/start", host)
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

// Status returns the status of the virtual machine.
func Status(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/vm/status", host)
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

// Stop will stop a virtual machine.
func Stop(json string, host string) (string, error) {
	path := fmt.Sprintf("http://%s/v1/vm/stop", host)
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
