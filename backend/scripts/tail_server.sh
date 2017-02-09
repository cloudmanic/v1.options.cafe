#!/bin/bash

echo ""
echo ""
echo \"Tailing the server log just to make sure everything went ok. (Control-C when done)\";
echo ""
echo ""

ssh -t -p 9022 spicer@app.options.cafe "sudo -- sh -c 'sudo tail -f  /home/deploy/logs/app.options.cafe.out.log'"