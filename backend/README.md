## Getting Up And Running For Development

* Make sure the GoPATH is setup. This is what I have set ```export GOPATH=$HOME/Development/golang:$HOME/Development```

* In the case of ```$HOME/Development``` Make sure you add a symlink with src. ```ln -s /Users/spicer/Development/app.options.cafe /Users/spicer/Development/src/github/app.options.cafe```

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

* ```cd cron && go run *.go --action=symbol-import```


## Redis 

* We use Redis 4.0

* We use Redis. On OSX you install it with `brew install redis`

* OSX: `brew services list` | `brew services redis start` | `brew services restart redis`