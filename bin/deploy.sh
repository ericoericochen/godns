#!/bin/bash

# Fail if any command in a pipeline fails
set -o pipefail

# Read the deployment config
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo_dir="$script_dir/.."
deployment_config_json_path="$script_dir/../deployment.config.json"
config=$(<"$deployment_config_json_path")

# Load env variable
env_file="$script_dir/../.env"
source "$env_file"

# Load the VM config
name=$(jq -r  '.name' <<< "$config")
region=$(jq -r  '.region' <<< "$config")
size=$(jq -r '.size' <<< "$config")
image=$(jq -r '.image' <<< "$config")
env=$(jq '.env' <<< "$config")

# Check if we already have a VM deployed
if droplet=$(doctl compute droplet list --output json | 
    jq --arg name "$name" '.[] | select(.name == $name)'); then

    if [[ -n "$droplet" ]]; then
        echo "Droplet Found"
    else
        # Deploy the VM if it doesn't exist
        echo "Droplet Not Found"
        echo "Creating Droplet '$name'..."
        droplet=$(doctl compute droplet create "$name" \
            --region "$region" \
            --size "$size" \
            --image "$image" \
            --ssh-keys "$DIGITAL_OCEAN_SSH_KEY_ID" \
            --wait --output | jq -e '.[0]' || exit 1
        
        )
        echo "Droplet '$name' created"
    fi
    
    echo "$droplet"

    # SSH into the VM and run setup
    droplet_ipv4=$(jq -r '.networks.v4[] | select(.type == "public") | .ip_address' <<< "$droplet")

    ssh "root@$droplet_ipv4" 'bash -s' <<'REMOTE'
    # install docker and start it
    apt install docker.io -y
    systemctl start docker
    docker info
REMOTE

    # Build the docker image and upload it to the VM
    godns_image="godns:latest"
    docker build --platform linux/amd64 -t "$godns_image" "$repo_dir" || exit 1
    docker save "$godns_image" | ssh "root@$droplet_ipv4" 'docker load' || exit 1

    # Stop the current running container and start the latest one
    printf -v remote_args '%q ' "$droplet_ipv4" "$env"
    ssh "root@$droplet_ipv4" "bash -s -- $remote_args" <<'REMOTE'
    droplet_ipv4="$1"
    env="$2"

    container_id=$(docker ps -aq --filter 'name=^godns$') || exit 1
    if [[ -n "$container_id" ]]; then
        docker stop godns
        docker rm godns
    fi

    env_args=()
    while IFS= read -r pair; do
        env_args+=(-e "$pair")
    done < <(jq -r 'to_entries[] | "\(.key)=\(.value)"' <<< "$env")

    echo "Starting server..."
    echo "Binding port $droplet_ipv4:53:53/udp"
    docker run -d -p "$droplet_ipv4:53:53/udp" "${env_args[@]}" --name godns godns:latest
REMOTE
else
    echo "❌ Droplet lookup failed" >&2
    exit 1
fi