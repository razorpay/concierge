#!/bin/bash
set -euo pipefail

ALOHOMORA_BIN=$(which alohomora)
cd /app
# casting alohomora to unlock the secrets
echo "<?php return 'production';" > "environment/env.php"
$ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/env.sh.j2"
# TODO: try setting up unicreds once we resolve
# https://github.com/Versent/unicreds/issues/75
source "dockerconf/env.sh"
php artisan custom:leasemanager
