#!/bin/bash
go build src/gapp/gluster/gluster.go
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.147:/home/john/gluster
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.186:/home/john/gluster
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.196:/home/john/gluster

sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.147:/home/john/gluster.json
sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.186:/home/john/gluster.json
sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.196:/home/john/gluster.json