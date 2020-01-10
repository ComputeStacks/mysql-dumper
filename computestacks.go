package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Project struct {
	ID		int			`json:"id"`
	Name	string		`json:"name"`
}

type ContainerImage struct {
	ID			int			`json:"id"`
	Label		string		`json:"label"`
	Role		string		`json:"role"`
	Category	string		`json:"category"`
}

type Container struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
	IP		string	`json:"ip"`
}

type IngressRule struct {
	Proto				string		`json:"proto"`
	Port				int			`json:"port"`
	ExternalAccess		bool		`json:"external_access"`
	BackendSSL			bool		`json:"backend_ssl"`
	TcpProxyOpt			string		`json:"tcp_proxy_opt"`
}

type ContainerService struct {
	ID				int					`json:"id"`
	Name			string				`json:"name"`
	Label			string				`json:"label"`
	CreatedAt		time.Time			`json:"created_at"`
	UpdatedAt		time.Time			`json:"updated_at"`
	Domains			[]string			`json:"domains"`
	Image			ContainerImage		`json:"image"`
	Containers		[]Container			`json:"containers"`
	IngressRules	[]IngressRule		`json:"ingress_rules"`
	Settings		[]ServiceParameter	`json:"settings"`
}

type ServiceParameter struct {
	ID			int			`json:"id"`
	Name		string		`json:"name"`
	Label		string		`json:"label"`
	ParamType	string		`json:"param_type"`
	Value		string		`json:"decrypted_value"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type ProjectMetadataList struct {
	CurrentServiceID	int					`json:"current_service_id"`
	Project				Project				`json:"project"`
	Services			[]ContainerService	`json:"services"`
}

func loadContainers() ([]Instance, error) {

	var currentService ProjectMetadataList

	var instances []Instance

	response, err := LoadMetaData()

	if err != nil {
		return instances, err
	}

	responseData, _ := ioutil.ReadAll(response.Body)
	jsonErr := json.Unmarshal(responseData, &currentService)
	if jsonErr != nil {
		return instances, jsonErr
	}

	for _, value := range currentService.Services {

		if value.Image.Role != "mysql" {
			continue
		}

		var ipAddr string
		var rootPassword string

		for _, setting := range value.Settings {
			if setting.Name == "mysql_password" {
				rootPassword = setting.Value
			}
		}

		for _, c := range value.Containers {
			ipAddr = c.IP
		}

		instances = append(instances, Instance{
			IPAddress: ipAddr,
			Password: rootPassword,
		})
	}

	return instances, nil

}

func LoadMetaData() (*http.Response, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(3*time.Second, func() {
		cancel()
	})
	client := &http.Client{}
	endpoint := os.Getenv("METADATA_URL")
	endpointAuth := os.Getenv("METADATA_AUTH")
	request, reqError := http.NewRequest("GET", endpoint, bytes.NewBuffer([]byte{}))
	if reqError != nil {
		return nil, reqError
	}
	request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer " + endpointAuth)
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
