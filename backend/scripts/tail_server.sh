#!/bin/bash

ssh -t -p 9022 spicer@app.options.cafe "sudo -- sh -c 'sudo tail -f  /home/deploy/logs/app.options.cafe.out.log'"