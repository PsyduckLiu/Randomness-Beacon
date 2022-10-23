#!/bin/bash
for i in $(seq 0 6)
do
go run main.go $i > result/result$i.txt &
sleep 1
echo "node $i is running"
port=3000$i
PID=$(sudo netstat -nlp | grep "$port" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')
echo ${PID} >> result/running.pid
done

# go run test/client.go 0 > result/client.txt &

wait
echo "all nodes are closed"