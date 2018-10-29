## App: app.options.cafe

This is the options trading application we offer to customers. Angular2 front-end, Golang backend, Electron desktop application. 

## Deploying With Ansible (and other actions by hand)

We assume the server is already setup for us. At Cloudmanic we have a different git repo for setting up docker servers.

* To deploy we do this ```cd scripts && ./deploy.sh```

* We mostly use Ansible for manage our deploys.

* We use ```ansible-vault``` to encrypt our files. (ask someone for the password). We can use ```ansible-vault edit filename``` to edit these files. 

* You need to add ```ansible/.vault_pass``` with the value password as the only string in the file. https://www.digitalocean.com/community/tutorials/how-to-use-vault-to-protect-sensitive-ansible-data-on-ubuntu-16-04


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