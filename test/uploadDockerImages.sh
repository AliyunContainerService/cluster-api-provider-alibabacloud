#!/bin/bash

image=${1}
sshhost=${2}
sshkey=${3}

docker save ${image} | bzip2 | ssh -i ${sshkey} ec2-user@${sshhost} "bunzip2 > /tmp/tempimage.bz2 && sudo docker load -i /tmp/tempimage.bz2"
echo "image: $image"
echo "sshhost: $sshhost"
