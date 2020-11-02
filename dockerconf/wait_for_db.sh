#!/bin/sh
set -euo pipefail

echo "Waiting for MySQL service start...";
while ! nc -z $DB_HOST 3306;
do
    sleep 1;
done;
echo "Connected!";
