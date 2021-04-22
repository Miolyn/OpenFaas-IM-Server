#!/bin/bash

# start 1
nohup ./connect/main > /tmp/connect.log 2>&1 &
nohup ./logic/main > /tmp/logic.log 2>&1 &
nohup ./user/main > /tmp/user.log 2>&1 &
nohup ./proxy/main > /tmp/proxy.log 2>&1 &
# just keep this script running
#while [[ true ]]; do
#    sleep 1
#done
