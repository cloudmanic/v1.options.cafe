## App: app.options.cafe

This is the options trading application we offer to customers. Angular2 front-end, Golang backend, Electron desktop application. 

## Getting Up And Running For Development

* Make sure the GoPATH is setup. This is what I have set ```export GOPATH=$HOME/Development/golang:$HOME/Development```

* In the case of ```$HOME/Development``` Make sure you add a symlink with src. ```ln -s /Users/spicer/Development /Users/spicer/Development/src```

## Deploying With Ansible (and other actions by hand)

* Login to the new host as root and run this command : ```apt-get install python```

* Cd into the ansible directory and run this command : ```ansible-playbook -i inventory bootstrap.yml```

* Log back into the server as root and reboot : ```reboot```

* The server is running on a new ssh port now. We now run this : ```ansible-playbook --ask-sudo-pass app.yml```

* Log back into the server and setup mysql by hand.

* Setup ```/home/deploy/.env``` (make sure to close the perms off)

* Run this ```sudo systemctl enable supervisor.service```

* Run this ```sudo mkdir /etc/letsencrypt``` ```sudo chmod 700 /etc/letsencrypt``` ```sudo chown deploy:deploy /etc/letsencrypt/```

* When deploying this is useful for kick starting https://skitch.cloudmanic.com/Lightsail_1E4C5E91.png

* Run this ```sudo rm -rf /root/.ssh```

* Run this ```sudo passwd root``` This sets root's password. Maybe we do not want this. Who knows.....

## Note On Backend Docker

* ```go run *.go``` the docker way : ```cd backend/docker && dc run --rm backend```

* Docker .env file example located in ```backend/docker/.env```

```
P_UID=501
P_GID=20
RESTART=always
VIRTUAL_HOST=app.options.dev
COMPOSE_PROJECT_NAME=app.options.dev
GO_PATH=/Users/spicer/Development/golang
```