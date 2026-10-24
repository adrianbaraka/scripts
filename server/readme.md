# Server configuration checklist

This guide assumes a Debian/Ubuntu server and an existing root login. Work through it in order. Keep **two independent root sessions** open until key-based login, `sudo`, SSH validation, and the firewall have all been tested. Never close the known-good session while changing remote access.

## 1. Prerequisites

- A Debian/Ubuntu server with a working root console or SSH session.
- A local computer on which you can create and retain an SSH private key securely.
- The server's public IP address or hostname.
- Your normal public IP address, only if you intend to use it in Fail2ban's `ignoreip` setting.
- For Beszel, a reachable Tailscale address for the monitored server and a running Beszel Hub.

Replace `abc`, `server.example`, `203.0.113.10`, and paths with your values. Do not copy example IP addresses or private-key paths literally.

## 2. Create an administrator account

Update packages first. `full-upgrade` is separate from `update`; do not combine them as `apt update upgrade`.

```bash
apt update && apt full-upgrade
adduser --disabled-password --gecos "" abc
usermod -aG sudo abc
```

Create a passwordless-sudo rule only if that is an intentional trade-off for this server. It grants this account unrestricted root access without asking for its password.

```bash
visudo -f /etc/sudoers.d/00-abc-sudo.conf
```

Add this line in the editor:

```sudoers
abc ALL=(ALL:ALL) NOPASSWD: ALL
```

`visudo` validates the file before saving. Confirm the resulting permissions and syntax explicitly:

```bash
chmod 440 /etc/sudoers.d/00-abc-sudo.conf
visudo -c
```

## 3. Set up and test SSH keys

On the **local computer**, create a dedicated modern key pair if you do not already have one:

```bash
ssh-keygen -t ed25519 -a 100 -f ~/.ssh/server_ed25519 -C "abc@server"
```

Protect the private key; copy only the `.pub` public-key file to the server. While still logged in as root on the server:

```bash
install -d -m 700 -o abc -g abc /home/abc/.ssh
nano /home/abc/.ssh/authorized_keys
chown abc:abc /home/abc/.ssh/authorized_keys
chmod 600 /home/abc/.ssh/authorized_keys
```

Paste the contents of `~/.ssh/server_ed25519.pub` into `authorized_keys` as one complete line. From a **new local terminal**, test the new account before changing SSH settings:

```bash
ssh -i ~/.ssh/server_ed25519 abc@127.0.0.1
sudo -v
sudo apt update
```

Stop here if either SSH or `sudo` fails. Keep the root session open and fix the problem first.

### Optional local SSH shortcut

Add this on the local machine as `~/.ssh/config` (not on the server):

```sshconfig
Host my-server
    HostName 127.0.0.1
    Port 22
    User abc
    IdentityFile ~/.ssh/server_ed25519
    IdentitiesOnly yes
```

Then connect with `ssh my-server`.

## 4. Lock down SSH safely

Only do this after the key login above works in a separate terminal. Put the intended settings in a drop-in file, for example `/etc/ssh/sshd_config.d/00-abc-ssh.conf`:

````bash
sudo nano /etc/ssh/sshd_config.d/00-abc-ssh.conf
````

```sshconfig
PermitRootLogin no
PasswordAuthentication no
KbdInteractiveAuthentication no
PubkeyAuthentication yes
X11Forwarding no
MaxAuthTries 3

# AllowUsers abc
```

Before applying anything, validate the complete configuration:

```bash
sudo sshd -t
```

On Debian/Ubuntu, the SSH service is normally named `ssh`, not `sshd`. Reload it only after `sshd -t` succeeds:

```bash
sudo systemctl reload ssh
sudo systemctl status ssh --no-pager
```

If `reload` is unavailable, use `sudo systemctl restart ssh` only after validation and while the existing working session remains open. In another terminal, confirm all three outcomes:

```bash
ssh root@203.0.113.10
ssh abc@203.0.113.10
ssh -i ~/.ssh/server_ed25519 abc@203.0.113.10
```

The first two should fail (root and password-only access); the third should succeed. Do not disconnect the original session until it does.

## 5. Configure Fail2ban

Install Fail2ban and copy your maintained jail configuration into the drop-in directory:

```bash
sudo apt update && sudo apt install fail2ban
sudo install -m 644 00-abc-fail2ban.conf /etc/fail2ban/jail.d/00-abc-fail2ban.conf
```

Enable only jails for services you actually run—for example, `sshd`. If you set `ignoreip`, understand that those addresses can never be banned by Fail2ban. Whitelist only a stable, trusted address or subnet; do not add a broad range merely to avoid a lockout. A simple SSH jail might include:

```ini
[sshd]
enabled = true
```

Validate and start it:

```bash
sudo fail2ban-client -t
sudo systemctl restart fail2ban
sudo systemctl status fail2ban --no-pager
sudo fail2ban-client status
sudo fail2ban-client status sshd
```

The repository includes `scripts/jails.sh`, use it as a convenient status overview after the service is configured.

## 6. Configure UFW

Install UFW and inspect its application profiles:

```bash
sudo apt update && sudo apt install ufw
sudo ufw app list
sudo ufw app info OpenSSH
sudo ufw status numbered
sudo ufw show added
```

**Allow SSH before enabling the firewall**, otherwise you can lock yourself out. Prefer the profile when it is available:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status verbose
```

Only open ports required by this host. For example, also allow `443/udp` only if you intentionally serve HTTP/3/QUIC. If SSH runs on a nonstandard port, allow that exact port instead of `OpenSSH`. Changes made after UFW is enabled can be applied with:

```bash
sudo ufw reload
```

Test a fresh SSH connection after enabling UFW before closing the working session.

## 7. Add swap (if needed)

Check current memory and swap first:

```bash
swapon --show
free -h
```

For a 2 GiB swap file, run the following. Do not create a second `/swapfile` if one already exists.

```bash
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

Make it persistent without duplicating the `/etc/fstab` entry:

```bash
grep -qE '^/swapfile[[:space:]]' /etc/fstab || echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
sudo swapon --show
```

`fallocate` works on typical ext4/xfs filesystems. On filesystems that do not support it reliably, create the file with `dd` instead.

## 8. Install apt-listbugs

`apt-listbugs` warns about Debian bugs affecting packages before installation or upgrades. It is primarily useful on Debian and may not be available or relevant on Ubuntu.

```bash
sudo apt update && sudo apt install apt-listbugs
```

Read its warnings during future upgrades rather than blindly accepting a problematic package.

## 9. Add the server to Beszel

The Beszel Hub creates the agent credentials. In the Hub:

1. Add a new system.
2. Copy the SSH key shown for that system into the agent deployment's `keys.txt` file.
3. Set the monitored server's Tailscale IP in the agent `.env` file.
4. From the directory containing the agent's Compose file, start it:

    ```bash
    docker compose up -d
    docker compose ps
    ```

Use the Tailscale IP only when the Hub can reach it over your tailnet. Keep `keys.txt` and `.env` private: they are credentials, not files to commit or share. Confirm in Beszel Hub that the system connects and reports metrics.

## Final checks

Before ending the root session, verify:

- Key-only SSH works for `abc` in a brand-new terminal.
- `sudo -v` works for `abc`.
- Root SSH and password authentication are refused.
- `systemctl status ssh`, `fail2ban`, and `ufw status verbose` show the expected state.
- The host remains reachable after UFW is enabled.
- Beszel shows the new system as online.
