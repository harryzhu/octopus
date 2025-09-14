#!/bin/bash
cmd_1=$1
echo "$cmd_1"

function cmdStop() {
  currentPID=$(pidof octopus)
  if [ x"$currentPID" != x ]; then
    echo kill "$currentPID"
    kill -9 "$currentPID"
  fi
}

function cmdStart() {
source /etc/profile
/usr/local/bin/octopus server --debug=false --host=localhost --port=9090 --admin-user=admin --admin-password=123
#& >/dev/null 2>&1
}

if [ "$cmd_1" == "stop" ]; then
  cmdStop
fi

#
if [ "$cmd_1" == "start" ]; then
  cmdStart
fi

#
if [ "$cmd_1" == "restart" ]; then
  cmdStop
  cmdStart
fi

