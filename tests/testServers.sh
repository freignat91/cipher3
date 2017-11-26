#!/bin/bash

cipher3 stopServer :3001
cipher3 stopServer :3002
sleep 1
cipher3 launchServer :3001 test1
cipher3 launchServer :3002 test2
sleep 1
cipher3 handcheckServer :3001 :3002
cipher3 sendData :3001 test2 "Hey, how are you test2?"
cipher3 sendData :3002 test1 "good!"
