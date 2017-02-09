#!/bin/bash

echo "Building app.options.cafe"
./build.sh

echo "Copying binary to app.options.cafe"
scp -P 9022 ../builds/options_cafe.linux.amd64 spicer@138.197.50.228:/tmp/options_cafe.linux.amd64

echo "Configuring app.options.cafe on server"
ssh -t -p 9022 spicer@138.197.50.228 "sudo -- sh -c 'supervisorctl stop app.options.cafe; mv /tmp/options_cafe.linux.amd64 /home/deploy/options_cafe.linux.amd64; chown deploy:deploy /home/deploy/options_cafe.linux.amd64; chmod 500 /home/deploy/options_cafe.linux.amd64; chmod 600 /home/deploy/.env; chown deploy:deploy /home/deploy/.env; setcap CAP_NET_BIND_SERVICE=+eip /home/deploy/options_cafe.linux.amd64; supervisorctl start app.options.cafe'"

echo "Tailing the server log just to make sure everything went ok. (Control C when done)"
./tail_server.sh