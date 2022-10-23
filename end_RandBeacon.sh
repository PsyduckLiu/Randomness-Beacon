#!/bin/bash
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid
cat configInit.yml > config.yml