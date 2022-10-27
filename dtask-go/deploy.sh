#!/bin/bash

# --------- remote
remote="alphaboi@curtisnewbie.com"
remote_path="~/services/dtaskgo/build/"
remote_config_path="~/services/dtaskgo/config/"
# ---------

#scp "run.sh" "${remote}:${remote_path}"
#scp "Dockerfile" "${remote}:${remote_path}"

# copy the config file just in case we updated it
# scp "app-conf-dev.json" "${remote}:${remote_path}"

scp -r dtask/* "${remote}:${remote_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi
scp dtask/app-conf-prod.json "${remote}:${remote_config_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "mv ~/services/dtaskgo/logs/dtaskgo.log ~/services/dtaskgo/logs/dtaskgo-$(date +'%Y%m%d_%H%M%S').log"
ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build dtaskgo"
