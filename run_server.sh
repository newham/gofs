#!/bin/bash

f=gofs.service
name=gofs
exc=gofs

sudo rm ./$exc

wget https://github.com/newham/gofs/releases/download/v2.0/$exc

sudo chmod +x $exc

echo "[Unit]
Description=$name

[Service]
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/$exc
User=$(whoami)

[Install]
WantedBy=multi-user.target
" > $f

sudo mv $f /lib/systemd/system

sudo systemctl enable $f

sudo systemctl restart $f

sudo systemctl status $f
