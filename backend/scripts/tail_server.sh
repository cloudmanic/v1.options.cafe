#!/bin/bash

ssh -t -p 9022 spicer@138.197.50.228 "sudo -- sh -c 'sudo tail -f  /home/deploy/logs/app.options.cafe.out.log'"