#!/bin/bash

if [ $# -lt 4 ]
then
    echo "请输入至少4个参数: 1: 算法的名字  2: 集群的大小  3: 每个实例使用的CPU个数(0表示不使用taskset) 4: CPU总个数(请确保CPU的总个数大于实例总共需要的CPU个数)"
    echo "例如: ./run_bft_server.sh pbft 4 2 16"
    exit
fi

algorithm=$1
bin=./$1server
echo $bin
clusterSize=$2
cpuPer=$3
cpuTotal=$4

log_dir="./log/${algorithm}"
if [ ! -d $log_dir ]; then
  mkdir -p -m 755 $log_dir
  echo "mkdir -p -m 755 ${log_dir} done"
fi


trap 'trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT

# 启动程序
for((i=1;i<=$clusterSize;i++));  
do
    $bin --self-id $i --privkey keys/r$i.key --cluster-size $clusterSize  > $log_dir/$i.out 2>&1 "$@" &
    echo $!
    if [ $cpuPer -ne 0 ]; then
        start=$(expr $cpuTotal - $cpuPer \* $(expr $clusterSize - $i + 1));
        end=$(expr $start + $cpuPer - 1);
        echo $start
        echo $end
        taskset -pac $start-$end $! # 用taskset指定使用的cpu
    fi    
done  


for((i=1;i<=$clusterSize;i++));  
do
    wait;
done
