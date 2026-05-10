#!/usr/bin/env bash
# to auto add custom dns configs
# not currently needed

declare -A hostnames


domain="home.arpa"
hostIp=192.168.1.100 # or tailscale ip


hostnames["jellyfin"]="$hostIp"
hostnames["hub"]="$hostIp"
hostnames["dozzle"]="$hostIp"
hostnames["beszel"]="$hostIp"
hostnames["uptime"]="$hostIp"
hostnames["pihole"]="$hostIp"
hostnames["kavita"]="$hostIp"
hostnames["portainer"]="$hostIp"

printf "CUSTOM_DNS='"
for host in "${!hostnames[@]}"; do
    printf "%s %s.%s;" "${hostnames["$host"]}" "${host}" "$domain"
done
printf "'\n"