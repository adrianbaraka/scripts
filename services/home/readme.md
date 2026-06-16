# My homelab config and setup

- My homelab running a media server(jellyfin), Books servers(kavita) and a couple of tools.
- All of them run via docker containers orchestrated by docker compose.
- For external access I use [tailscale](https://tailscale.com/) and a free domain from [dnsexit](https://dnsexit.com/).

### Prerequisites
- Install [docker](), [docker compose]().
- Optionally setup tailscale and a domain name.
### Setup

1. Copy the `.env.example` file to `.env` and populate the values.
2. Create the shared Docker network that allows all containers to communicate.
    ````bash
    docker network create home
    ````
3. Run `docker compose -f compose.apps.yaml up -d`. This starts the core apps.
4. Run `docker compose -f compose.utils.yaml up -d`. This starts the utility stack.
5. For beszel add a new system in the beszel hub dashboard and populate the given crendentials in the given file and finally run `docker compose -f compose.beszel.agent.yaml up -d`


### Additional notes
- The caddy image needs to be rebuilt using the [extension](https://caddyserver.com/docs/modules) specific to the domain registrar of choice. 
- A dockerfile with dnsexit is in the files directory.
    ````bash
    # build the image from the dockerfile
    docker build -t adrianbaraka/caddy-dnsexit:latest .
    ````

### File Structure