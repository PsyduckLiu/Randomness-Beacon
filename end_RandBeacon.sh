#!/bin/bash
cd ./ConsensusNode
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ../EntropyNode
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ..
cat configInit.yml > config.yml
# kill -9 `ps -ef |grep main|awk '{print $2}'`
for i in $(seq 0 9)
do
kill -9 `ps -ef |grep exe/main\ $i|awk '{print $2}'`
done