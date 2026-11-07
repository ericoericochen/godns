package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const DIGITAL_OCEAN_TOKEN = "DIGITAL_OCEAN_TOKEN"
const DEPLOYMENT_CONFIG_FILE = "deployment.config.json"

type DeploymentConfig struct {
	Droplet string            `json:"droplet"`
	Env     map[string]string `json:"env"`
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

func main() {
	// load digital ocean token
	doToken, exists := os.LookupEnv(DIGITAL_OCEAN_TOKEN)
	if !exists {
		panic(fmt.Sprintf("Missing %v", DIGITAL_OCEAN_TOKEN))
	}

	fmt.Println(doToken)

	// load deployment config
	config := loadDeploymentConfig()
	fmt.Println(config)
}
