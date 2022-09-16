#!/bin/bash

# --------- remote
remote="alphaboi@curtisnewbie.com"
remote_path="~/services/dtaskgo/build"
# ---------

#scp "run.sh" "${remote}:${remote_path}"
#scp "Dockerfile" "${remote}:${remote_path}"

# copy the config file just in case we updated it
# scp "app-conf-dev.json" "${remote}:${remote_path}"

scp -r dtask/* "${remote}:${remote_path}/"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build dtaskgo"
