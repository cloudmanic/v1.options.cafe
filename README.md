## App: app.options.cafe

This is the options trading application we offer to customers. Angular2 front-end, Golang backend, Electron desktop application. 

## Deploying With Ansible (and other actions by hand)

We assume the server is already setup for us. At Cloudmanic we have a different git repo for setting up docker servers.

* To deploy we do this ```cd scripts && ./deploy.sh```

* We mostly use Ansible for manage our deploys.

* We use ```ansible-vault``` to encrypt our files. (ask someone for the password). We can use ```ansible-vault edit filename``` to edit these files. 

* You need to add ```ansible/.vault_pass``` with the value password as the only string in the file. https://www.digitalocean.com/community/tutorials/how-to-use-vault-to-protect-sensitive-ansible-data-on-ubuntu-16-04

# Deploying Servers

* When deploying a server with Digital Ocean copy the following into the `User-Data` filed. It will run Cloud Init when the VPS boots up.

```
#cloud-config
users:
  - name: spicer
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh-authorized-keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAsw21gLc2CaJN8hJB7zWDYWLF5gqWl6t8ozgso8aOrq8rz7P8ji3MwvHEelEe6UMNg4CxWTGYIWvFptlfCRvy9d94RBy9AAdb4pEBmSOyxPf8sJ+xD+V3TFJfmMOAm4049cBLN9b7+PRkUjl4jC3zTch5tQ+5lG7v04tWwzCaSCSD2HNuw2qKK3FpaLA6EIw+ieueBkgNgRnwMvgVO8nmyOkR5b3WUoL4vow3heNHV00V4M0yhBHLHDIFkXMgMztpLm3Dki1ZplUF0EyPH5llj5a4n2RMR5c7B1wAiXuUPO0oQTw9ItS5SZl9zKu9ZuIvqeXWsz/0NqRdEMIKqvxIZQ== spicer@cloudmanic.com
packages:
  - python
```

* Once a fresh server is up and running configure it with `ansible-playbook server.yml`



## Getting Up And Running

* Run `go run main.go -cmd=symbol-import` to import all the possible symbols into the symbols db table.


## Process To Add A New Options Strategy 

* Add a new trade type in `backend/library/archive/trade_types`

* Add the new trade type in `backend/library/archive/trade_group_classify.go` both functions: ClassifyTradeGroup, and GetAmountRiskedInTrade

* Update the ENUM values in `TradeGroup:type`. Typically I will update the model code, then update the ENUM locally, in testing DB, and then on production by hand.

* Add option in `frontend/src/app/trading/trades/home.component.html` 

* Add new file for strategy at `backend/screener`

* Update switch statement in `backend/screener/base.go::PrimeAllScreenerCaches` 

* Add option to screener `frontend/src/app/trading/screener/add-edit/add-edit.component.html`

* Update `backend/controllers/screeners.go` to include any new screener keys. The var at the top `screenerItemKeys`.

* Update the screener type in `ScreenFuncs` in `backend/screener/base.go` in init()

* Add new tradegroup to `frontend/src/app/models/trade-group.ts` in `TradeGroupsCont`
