#!/bin/bash
sudo chmod a+x ./bin/webhandler
sudo ln -s -T $PWD/bin/webhandler  /bin/webhandler
godep go build -o ./bin/websnapshot ./
sudo ln -s -T $PWD/bin/websnapshot /bin/websnapshot
