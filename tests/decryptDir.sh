#!/bin/bash

#$1: encrypted file source
#$2: targeted directory path
#$3: key  path

cipher3 decryptFile $1 /tmp/$1.tmp $3
if [ $? -eq 0 ]
then
  tar xvf /tmp/$1.tmp 
fi
rm /tmp/$1.tmp
