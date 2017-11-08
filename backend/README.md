## Getting Up And Running For Development

* Make sure the GoPATH is setup. This is what I have set ```export GOPATH=$HOME/Development/golang:$HOME/Development```

* In the case of ```$HOME/Development``` Make sure you add a symlink with src. ```ln -s /Users/spicer/Development /Users/spicer/Development/src```

## Note On Backend Development Docker

* ```go run *.go``` the docker way : ```cd backend/docker && docker-compose run --rm -p 7080:7080 app bash```

* Add a Docker .env file example located at ```backend/docker/.env```

```
P_UID=501
P_GID=20
RESTART=always
VIRTUAL_HOST=app.options.dev
COMPOSE_PROJECT_NAME=app.options.dev
GO_PATH=/Users/spicer/Development/golang
```

* When connecting to a docker host we have to do this in the ```backend/.env``` file. 

```DB_HOST=tcp(mysql-5.5:3306)```
 

## Backend Unit Testing....

* ```go test ./...``` from the root of the project
* ```go test app.options.cafe/backend/library/archive``` (as an example)
* ```go test -v app.options.cafe/backend/library/archive``` (if you want to see the output of say a println)

## Notes On Billings

Every user must have a subscription (table: ```User::StripeSubscription```) if this field is empty the user can not use the app. They should be presented with a screen to select a plan. New users have a trail period. So new users are assigned to a default subscription see .env ```STRIPE_DEFAULT_PLAN```. If after their trail ends and they do not add a payment source their subscription will be deleted via webhooks making ```User::StripeSubscription``` empty. 

## Best Way To Test Webhooks Locally

You can use https://ngrok.com to send webhooks in locally. The free account changes the url every time you run it. Here is how you start it ```ngrok http 7080```

* At stripe add this url ```https://0c1a1ee8.ngrok.io/webhooks/stripe``` (or whatever ngrok changes it to)

## Cron Jobs

Checkout ```cron/README.md```

## Getting Started With Data

Some data needs to be imported to kick the application off. For example import a list of symbols to populate the symbols table.

* ```go run *.go --action=symbol-import```

