#!/bin/bash

echo "Building options_cafe.linux.amd64..."

./build.sh

echo "Copying binary to app.options.cafe"

scp -P 9022 ../builds/options_cafe.linux.amd64 spicer@app.options.cafe:/tmp/options_cafe.linux.amd64

echo "Configuring app.options.cafe on server"

ssh -t -p 9022 spicer@app.options.cafe "sudo -- sh -c '
  supervisorctl stop app.options.cafe; 
  mv /tmp/options_cafe.linux.amd64 /home/deploy/options_cafe.linux.amd64; 
  chown deploy:deploy /home/deploy/options_cafe.linux.amd64; 
  chmod 500 /home/deploy/options_cafe.linux.amd64; 
  chmod 600 /home/deploy/.env; 
  chown deploy:deploy /home/deploy/.env; 
  setcap CAP_NET_BIND_SERVICE=+eip /home/deploy/options_cafe.linux.amd64; 
  supervisorctl start app.options.cafe;
  sleep 12;
  echo \"\";
  echo \"Checking out the server logs making sure there was no error on startup.\";
  echo \"\";
  tail /home/deploy/logs/app.options.cafe.out.log;
  echo \"\";
  echo \"\";
'"

echo "Deploy Done....."