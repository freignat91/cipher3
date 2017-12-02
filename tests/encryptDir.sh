#!/bin/bash

#$1: source directory path
#$2: target encrypted file path
#$3: key  path

tar cvf /tmp/$2.tmp $1
if [ $? -eq 0 ]
then
   cipher3 encryptFile /tmp/$2.tmp $2 $3
fi
rm /tmp/$2.tmp
