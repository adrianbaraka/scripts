# HomeLab

A containerized homelab environment featuring a secure reverse proxy with local SSL, centralized DNS, and a unified dashboard.

## Components
- **Caddy**: Reverse proxy with internal CA (automatic SSL for `.home.arpa`).
- **Pi-hole**: Network-wide ad-blocking and local DNS records.
- **Homepage**: Unified dashboard with service widgets.
- **Beszel**: Lightweight server resource monitoring.
- **Uptime Kuma**: Service availability monitoring.
- **Dozzle**: Real-time container log viewer.
- **Jellyfin**: Media streaming server.

## Setup Instructions

### 1. Networking
Create the shared Docker network that allows all containers to communicate:
```bash
docker network create home
```

### 2. Environment Configuration
Copy the example environment file and populate it with your specific keys and paths:
```bash
cp .env.example .env
nano .env
```


### 3. Deployment
The stack is split into core utilities and applications.

**Start the Utility Stack:**
Contains Pi-hole, Caddy, Homepage, Uptime Kuma, and Beszel Hub.
```bash
docker compose -f compose.utils up -d
```

**Start the Applications Stack:**
Contains Jellyfin and other services.
```bash
docker compose -f compose.apps.yaml up -d
```

### 4. Service Specific Notes

####  Pi-hole (DNS)
- **Initial Boot:** Pi-hole requires an active internet connection to pull initial blocklists.
- **Switching:** Once the container is healthy, you can change your system/router DNS settings to point to this Pi-hole instance.

#### Beszel (Monitoring)
1. Log in to the Beszel Web UI (`https://beszel.home.arpa`).
2. Add a new system.
3. Copy the generated credentials into your `.env` file:
   - `BESZEL_HUB_URL`
   - `BESZEL_LISTEN`
   - `BESZEL_TOKEN`
   - `BESZEL_KEY`
4. Start the agent to begin monitoring the host:
```bash
docker compose -f compose.beszel.agent.yaml up -d
```

#### Tailscale

- For tailscale get a auth key from the tailscale console and copy it to the .env file.
- Then in the DNS tab add a custom nameserver with the ip to that of the tailcale ip and restrict the domain to `home.arpa`. 
- This way any devices on the tailscale will have **.home.arpa* routed to pihole and successfully resolved to the actual service by caddy.
- Easier to just use tailscale on the host rather than as a docker container.
- Install [tailscale](https://tailscale.com/docs/install).
    ``` bash
        sudo tailscale up --hostname homelab --ssh --accept-dns=true
    ````

#### SSL Certificates
This setup uses Caddy's internal CA for `.home.arpa` domains. To trust the certificates:
1. Copy `root.crt` from the Caddy container:
   `docker cp caddy:/data/caddy/pki/authorities/local/root.crt ./`
2. Import and trust this certificate in your OS or Browser's Certificate Authority store.
    ```bash
    # Update System Store linux
    sudo cp root.crt /usr/local/share/ca-certificates/caddy-home.crt
    sudo update-ca-certificates --fresh

    # browser stores
    sudo apt install libnss3-tools -y
    certutil -d sql:$HOME/.pki/nssdb -A -t "C,," -n "Caddy Local CA" -i /usr/local/share/ca-certificates/caddy-home.crt

### Adding a New Service
To add a new application (e.g., `grafana`):
1.  **Caddy**: Add the subdomain to the `Caddyfile`.
2.  **DNS**: Add a local DNS record in Pi-hole pointing `grafana.home.arpa` to your server IP.
3.  **Homepage**: Add the service to `services.yaml`. If it requires an API key:
    * Add `HOMEPAGE_VAR_GRAFANA_KEY=${GRAFANA_KEY}` to your `utils-compose.yaml` environment.
    * Add `GRAFANA_KEY=your_key` to your `.env` file.
    * In `services.yaml`, use: `key: "{{HOMEPAGE_VAR_GRAFANA_KEY}}"`

## File Structure
- `utils-compose.yaml`: Core infrastructure (Caddy, Pi-hole, Homepage, etc).
- `apps-compose.yaml`: Applications (Jellyfin, etc).
- `./files/Caddyfile`: Reverse proxy and SSL configuration.
- `./files/services.yaml`: Homepage dashboard layout.
- `.env`: API keys and environment variables.
