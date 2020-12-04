#!/bin/bash

calCPU() {
    sleep 5 # 让server们冷静一下

    cpuresult=`ps aux | grep ${algorithm}server | grep -v 'grep' | awk '{print $10}' | tr "\n" ","|sed -e 's/,$/\n/'|tr "," " "`

    echo "cpus(hh:mm:ss):" $cpuresult >> $clientOutput
    ./getCPU.sh $cpuresult >> $clientOutput
}

if [ $# -lt 4 ]
then
    echo "请输入至少4个参数: 1: 算法的名字  2: 集群的大小  3: 使用的CPU个数(0表示不使用taskset) 4: prepare-num 5: viewchange-num(默认是10000)"
    echo "例如: ./cmdclient.sh pbft 4 15 10 10000"
    exit
fi

algorithm=$1
bin=./$1vcclient
clusterSize=$2
cpuPer=$3
prepareNum=$4

if [ $# -eq 5 ]
then
    viewchangeNum=$5
else    
    viewchangeNum=10000
fi

log_dir="./log"
if [ ! -d $log_dir ]; then
  mkdir -p -m 755 $log_dir
  echo "mkdir -p -m 755 ${log_dir} done"
fi

#clientOutput=$log_dir/${algorithm}client.txt
clientOutput=$log_dir/vcclient.txt

# 启动程序

trap 'calCPU && trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT

cmd="$bin --tls=true --cluster-size $clusterSize --payload-size 0 --prepare-num $prepareNum --viewchange-num $viewchangeNum $@"

echo "$cmd"
echo >> $clientOutput # 输入一个空行
echo "$cmd" >> $clientOutput

eval "$cmd $@ >> $clientOutput 2>&1 &"

if [ $cpuPer -ne 0 ]; then
    start=0
    end=$(expr $start + $cpuPer - 1);
    taskset -pac $start-$end $! # 用taskset指定使用的cpu
fi

wait;
