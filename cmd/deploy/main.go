package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
)

const DIGITAL_OCEAN_TOKEN = "DIGITAL_OCEAN_TOKEN"
const DEPLOYMENT_CONFIG_FILE = "deployment.config.json"

type DeploymentConfig struct {
	Name   string            `json:"name"`
	Region string            `json:"region"`
	Size   string            `json:"size"`
	Image  string            `json:"image"`
	Env    map[string]string `json:"env"`
	Port   string            `json:"port"`
}

func loadDeploymentConfig() DeploymentConfig {
	file, err := os.ReadFile(DEPLOYMENT_CONFIG_FILE)
	if err != nil {
		panic("")
	}

	var deploymentConfig DeploymentConfig
	if err := json.Unmarshal(file, &deploymentConfig); err != nil {
		panic("")
	}

	return deploymentConfig
}

func getOrCreateDropletIfNotExists(client *godo.Client, ctx context.Context, config DeploymentConfig) {

}

func main() {
	// load digital ocean token
	doToken, exists := os.LookupEnv(DIGITAL_OCEAN_TOKEN)
	if !exists {
		panic(fmt.Sprintf("Missing %v", DIGITAL_OCEAN_TOKEN))
	}
	doClient := godo.NewFromToken(doToken)

	// load deployment config
	config := loadDeploymentConfig()
	fmt.Println("config: ")
	fmt.Println(config)

	ctx := context.Background()

	// check if droplet exists
	droplets, _, err := doClient.Droplets.ListByName(ctx, config.Name, &godo.ListOptions{PerPage: 5})

	// load ssh keys
	keys, _, err := doClient.Keys.List(ctx, &godo.ListOptions{PerPage: 5})
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		fmt.Println(key)
	}

	// doClient.Droplets.Create(ctx, &godo.DropletCreateRequest{

	// })

	// droplets, _, err := doClient.Droplets.List(ctx, &godo.ListOptions{})
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }

	// fmt.Println("droplets: ", droplets)
}
