#!/bin/bash
sleep 5 # 让hotstuff server们冷静一下

hsOutput=./log/hotstuff/1.out

cpuresult=`ps aux | grep hotstuffserver | grep -v 'grep' | awk '{print $10}' | tr "\n" ","|sed -e 's/,$/\n/'|tr "," " "`

echo "cpus(hh:mm:ss):" $cpuresult >> $hsOutput
./getCPU.sh $cpuresult >> $hsOutput
