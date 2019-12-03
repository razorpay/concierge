#!/bin/sh
set -euo pipefail

echo "Waiting for MySQL service start...";
while ! nc -z $DB_HOST $DB_PORT;
do
    sleep 1;
done;
echo "Connected!";