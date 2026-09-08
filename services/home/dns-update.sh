#!/usr/bin/env bash

env_exists() {
    local var=$1

    if ! [[ -v $var ]]; then
        echo "The env var '$var' is not set." >&2
        exit 1
    fi
}

# load env file
source .env
env_exists DOMAIN
env_exists TAILSCALE_IP
env_exists DNS_EXIT_TOKEN

echo "Updating $DOMAIN DNS to $TAILSCALE_IP"

curl -fsS "https://api.dnsexit.com/dns/ud/" \
    -d "apikey=$DNS_EXIT_TOKEN" \
    -d "host=$DOMAIN,*.$DOMAIN" \
    -d "ip=$TAILSCALE_IP"