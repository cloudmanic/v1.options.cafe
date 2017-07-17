## Getting Up And Running For Development

* Make sure the GoPATH is setup. This is what I have set ```export GOPATH=$HOME/Development/golang:$HOME/Development```

* In the case of ```$HOME/Development``` Make sure you add a symlink with src. ```ln -s /Users/spicer/Development /Users/spicer/Development/src```

## Note On Backend Docker

* ```go run *.go``` the docker way : ```cd backend/docker && docker-compose run --rm -p 7652:7652 backend-dev bash```

* Docker .env file example located in ```backend/docker/.env```

```
P_UID=501
P_GID=20
RESTART=always
VIRTUAL_HOST=app.options.dev
COMPOSE_PROJECT_NAME=app.options.dev
GO_PATH=/Users/spicer/Development/golang
```

## Backend Unit Testing....

* ```go test ./...``` from the root of the project
* ```go test app.options.cafe/backend/library/archive``` (as an example)
* ```go test -v app.options.cafe/backend/library/archive``` (if you want to see the output of say a println)