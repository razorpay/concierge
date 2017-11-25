#!/bin/sh
echo "[+] $(date) Creating commit.txt"
echo "$GIT_COMMIT_HASH" > /app/public/commit.txt

echo "[+] $(date) Creating Log Directory"
mkdir -p /var/log/nginx

echo "[+] $(date) Ensuring nginx user has permissions"
chown -R nginx.nginx /app

echo  "[+] $(date) Fix Storage permissions"
cd /app/ && chmod -R o+wx storage/

echo "[+] $(date) Create /run/"
mkdir -p /run/
chown 0755 /run/
