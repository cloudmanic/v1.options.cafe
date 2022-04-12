## Getting Up And Running For Development

* Go Version - 1.10.3

* Make sure the GoPATH is setup. This is what I have set ```export GOPATH=$HOME/Development/golang:$HOME/Development```

* In the case of ```$HOME/Development``` Make sure you add a symlink with src. ```ln -s /Users/spicer/Development/app.options.cafe /Users/spicer/Development/src/github/app.options.cafe```

## Backend Unit Testing....

* ```go test ./...``` from the root of the project
* ```go test app.options.cafe/backend/library/archive``` (as an example)
* ```go test -v app.options.cafe/backend/library/archive``` (if you want to see the output of say a println)
* ```go test -run TestCreateBacktests01 ./controllers/...``` (to run one particular test)

We have a special database docker image we use for unit testing. This is so we can run many tests all at once. Each test creates
its own DB. We start this special image with ```scripts/start_testing_db.sh```. The database listens on port 9906. In our unit tests we should get a DB using this.

```
db, dbName, _ := models.NewTestDB("")
defer models.TestingTearDown(db, dbName)
```

When testing by hand you can pass in a non-empty string and set the DB to a string you know, so you can look in your db.

Every so often it is good to run this command ```docker restart options_cafe_testing```. It will clear out any old databases.

If you wanted to keep the database around so you can inspect it. Useful for when you are developing you can do this in your test. Just make sure you put it back before committing code.

```
db, _, _ := models.NewTestDB("oc_testing")
//defer models.TestingTearDown(db, dbName)
```

## Database tables

We use GORM to create our database tables. Checkout ```models/db.go```. To add a new table you must update ```doMigrations()```. Also in ```models/test.go``` you must update ```TruncateAllTables()```

## Notes On Billings

Every user must have a subscription (table: ```User::StripeSubscription```) if this field is empty the user can not use the app. They should be presented with a screen to select a plan. New users have a trail period. So new users are assigned to a default subscription see .env ```STRIPE_DEFAULT_PLAN```. If after their trail ends and they do not add a payment source their subscription will be deleted via webhooks making ```User::StripeSubscription``` empty.

## Best Way To Test Webhooks Locally

You can use https://ngrok.com to send webhooks in locally. The free account changes the url every time you run it. Here is how you start it ```ngrok http 7080```

* At stripe add this url ```https://0c1a1ee8.ngrok.io/webhooks/stripe``` (or whatever ngrok changes it to)

## Cron Jobs

Checkout ```cron/README.md```

## Getting Started With Data

Some data needs to be imported to kick the application off. For example import a list of symbols to populate the symbols table.

* ```go run main.go --cmd=symbol-import```


## Redis

* We use Redis 4.0

* We use Redis. On OSX you install it with `brew install redis`

* OSX: `brew services list` | `brew services redis start` | `brew services restart redis`

## User Status

In the users table a user has a status. Here is a summary of what that means.
'Active','Disable','Delinquent','Expired','Trial'

* ```Active```: User is outside of the free trail. User has a valid credit card on file and it charges. This is the state we want all users in.

* ```Disabled```: User has been disabled by an admin. This does not happen much. Mainly used if a user is abusing the system.

* ```Delinquent```: User is outside of the free trail, has added a credit card in the past, but for some reason we can't charge that card at the time of subscription renewal (monthly or yearly). Background feeds disable. Users are redirected to this url `/settings/account/upgrade`. Their strips subscriotion has been deleted. They have an customer profile in strope but not a subscription.

* ```Expired```: User's free trail has come to an end and they have not added a credit card yet to charge. Background feeds continue to work but the UI blocks the user from access the application. They are redirected to this page `/settings/account/expired`.

* ```Trial```: User is currently on the free trial every new user gets. They have full access to everything.
