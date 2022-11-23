#!/bin/bash
cd ../ConsensusNode
for i in $(seq 0 6)
do
go run main.go $i > result/result$i.txt &
sleep 1
echo "consensus node $i is running"
port=3000$i
PID=$(sudo netstat -nlp | grep "$port" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')
echo ${PID} >> result/running.pid
done

sleep 5

cd ../EntropyNode
for i in $(seq 0 59)
do
go run main.go $i > result/result$i.txt &
sleep 1
echo "entropy node $i is running"
echo $! >> result/running.pid
done

wait
echo "all nodes are closed"
