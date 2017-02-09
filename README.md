## app.options.cafe

This is the options trading application we offer to customers. Angular2 front-end, Golang backend, Electron desktop application. 

## Deploying With Ansible (and other actions by hand)

* Login to the new host as root and run this command : ```apt-get install python```

* Cd into the ansible directory and run this command : ```ansible-playbook -i inventory bootstrap.yml```

* Log back into the server as root and reboot : ```reboot```

* The server is running on a new ssh port now. We now run this : ```ansible-playbook --ask-sudo-pass app.yml```

* Log back into the server and setup mysql by hand.

* Setup ```/home/deploy/.env``` (make sure to close the perms off)

* Run this ```sudo systemctl enable supervisor.service```

* Run this ```sudo mkdir /etc/letsencrypt``` ```sudo chmod 700 /etc/letsencrypt``` ```sudo chown deploy:deploy /etc/letsencrypt/```