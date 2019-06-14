package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Container struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IPAddress  string `json:"local_ip"`
	Parameters struct {
		Environment struct {
			RootPassword string `json:"MYSQL_ROOT_PASSWORD"`
		} `json:"environment"`
	} `json:"parameters"`
}

type ContainerImage struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	Role     string `json:"role"`
	Category string `json:"role_class"`
}

type Service struct {
	ID         int            `json:"id"`
	Label      string         `json:"label"`
	Image      ContainerImage `json:"container_image"`
	Containers []Container    `json:"containers"`
}

type ListServicesResponse struct {
	Service Service `json:"deployment_container_service"`
}

func loadContainers() ([]Instance, error) {

	var services []ListServicesResponse

	var instances []Instance

	response, err := Get("/projects/" + os.Getenv("PROJECT_ID") + "/services")

	if err != nil {
		return instances, err
	}

	responseData, _ := ioutil.ReadAll(response.Body)
	jsonErr := json.Unmarshal(responseData, &services)
	if jsonErr != nil {
		return instances, jsonErr
	}

	for _, value := range services {
		for _, c := range value.Service.Containers {
			if value.Service.Image.Role != "mysql" {
				continue
			}
			instances = append(instances, Instance{
				IPAddress: c.IPAddress,
				Password: c.Parameters.Environment.RootPassword,
			})
		}
	}

	return instances, nil

}

func Get(path string) (*http.Response, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(3*time.Second, func() {
		cancel()
	})
	client := &http.Client{}
	endpoint := os.Getenv("API_HOST") // https://dev.computestacks.net/api
	fullPath := endpoint + path
	request, reqError := http.NewRequest("GET", fullPath, bytes.NewBuffer([]byte{}))
	if reqError != nil {
		return nil, reqError
	}
	request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json; api_version=50")
	request.Header.Set("Authorization", csAuth())
	response, err := client.Do(request)

	if err != nil {
		return response, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return response, nil
	} else {
		return response, errors.New("error status code received")
	}
}

func csAuth() string {
	auth := os.Getenv("API_KEY") + ":" + os.Getenv("API_SECRET")
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
