#!/usr/bin/env bash

# Fail2ban jails
# jails="nginx-bad-request nginx-botsearch nginx-forbidden nginx-http-auth nginx-limit-req sshd"
jails="sshd"

# Loop through each jail in the list
for jail in $jails; do
  echo "-------------------------------------------------------------------------"

  # Use fail2ban-client to get the status of the current jail
  if sudo fail2ban-client status "$jail"; then
      echo "Success"
    else
    echo "Error: Could not get status for jail: $jail.  It may not exist or Fail2ban may not be running."
  fi
  echo "---------------------------------------------------------------------"
done

echo "Finished checking Fail2ban jail statuses."
