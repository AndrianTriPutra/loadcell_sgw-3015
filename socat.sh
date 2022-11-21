#!/bin/bash

sleep 5
socat TCP-LISTEN:9090,fork,reuseaddr FILE:/dev/ttyUSB0,b9600,raw