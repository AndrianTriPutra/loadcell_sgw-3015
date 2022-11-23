#!/bin/bash

sleep 2
stty -F /dev/ttyUSB0 9600 cs7 -cstopb -parenb
sleep 3
socat TCP-LISTEN:9090,fork,reuseaddr FILE:/dev/ttyUSB0,b9600,raw