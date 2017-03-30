#!/bin/bash

ACCESS_TOKEN=20c0603580344cd68a99260b634be30e
ENVIRONMENT=production
LOCAL_USERNAME=`whoami`
REVISION=`git log -n 1 --pretty=format:"%H"`

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

# Tell Rollbar about this.

echo "Telling rollbar about this deploy"

curl https://api.rollbar.com/api/1/deploy/ \
  -F access_token=$ACCESS_TOKEN \
  -F environment=$ENVIRONMENT \
  -F revision=$REVISION \
  -F local_username=$LOCAL_USERNAME

echo ""

echo "Deploy Done....."