## App: app.options.cafe

This is the options trading application we offer to customers. Angular2 front-end, Golang backend, Electron desktop application. 

## Deploying With Ansible (and other actions by hand)

We assume the server is already setup for us. At Cloudmanic we have a different git repo for setting up docker servers.

* To deploy we do this ```cd scripts && ./deploy.sh```

* We mostly use Ansible for manage our deploys.

* We use ```ansible-vault``` to encrypt our files. (ask someone for the password). We can use ```ansible-vault edit filename``` to edit these files. 

* You need to add ```ansible/.vault_pass``` with the value password as the only string in the file. https://www.digitalocean.com/community/tutorials/how-to-use-vault-to-protect-sensitive-ansible-data-on-ubuntu-16-04


## Getting Up And Running

- Run `go run main.go -cmd=symbol-import` to import all the possible symbols into the symbols db table.