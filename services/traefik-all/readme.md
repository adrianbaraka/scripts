- In `/etc/nsswitch.conf` edit this line `hosts: files dns myhostname mdns4_minimal [NOTFOUND=return] mymachines` to ensure dns is before mdns to correctly handle the .local domain.
https://en.wikipedia.org/wiki/.local

https://doc.traefik.io/traefik/setup/docker/

mkcert -cert-file certs/local.crt -key-file certs/local.key "home.arpa" "*.home.arpa"

htpasswd -nb admin "P@ssw0rd" | sed -e 's/\$/\$\$/g'
admin:$$apr1$$BapxfnUp$$yH.D48WbpzPRAPR5Y9Hoc0


mkdir -p dynamic
cat > dynamic/tls.yml << EOF
tls:
  certificates:
    - certFile: /certs/local.crt
      keyFile: /certs/local.key
EOF
