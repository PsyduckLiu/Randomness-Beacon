#!/bin/bash
cd ./ConsensusNode
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ../EntropyNode
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ..
cat configInit.yml > config.yml
kill -9 `ps -ef |grep main|awk '{print $2}'`