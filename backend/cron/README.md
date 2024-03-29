## Cron - About 

Cron is designed to be a standalone app. It runs on its own and does utility activities such as downloading and storing data. This shares a code base with the core app so libraries can be shared.

Note: This is not run from typical unix cron. It has its own scheduling built in. Just run the app and leave it running.

The ```.env``` file is shared with the core application.

This app should be run from docker like all our other apps.

## Note On Backend Development Docker

* ```go run *.go``` the docker way : ```cd ../backend/docker && docker-compose run --rm cron bash```

* This assumes all the .env stuff was setup from the main app. ```../README.md```

* Might need to run ```go get``` from within the docker container

## To Run One Command At A Time.

You can simply pass ```--action={action name}```. Here is a list of commands you can run.

* ```go run *.go --action=symbol-import```