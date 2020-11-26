#!/bin/bash

function getSecond()
{
    string=$1  
    
    #对IFS变量 进行替换处理
    OLD_IFS="$IFS"
    IFS=":"
    array=($string)
    IFS="$OLD_IFS"

    cnt=0

    for var in ${array[@]}
    do
        let cnt+=1
    done

    second=0

    if [ $cnt -eq 3 ]; then
        second=$(($[10#${array[0]}]*3600 + $[10#${array[1]}]*60 + $[10#${array[2]}]))
    fi

    if [ $cnt -eq 2 ]; then
        second=$(($[10#${array[0]}]*60 + $[10#${array[1]}]))
    fi

    return $second
}


if [ $# -lt 2 ]
then
    echo "请输入至少2个值: leader cpu时间, 多个 follower cpu时间"
    echo "例如: ./getCPU.sh 01:2:3 1:2 01:02"
    exit
fi

getSecond $1

leaderCPU=$second # 第一个CPU就是leader的cpu
#min=$second

for i in $@
do
    getSecond $i
    # [ $max -lt $second ] && max=$second
    #[ $min -gt $second ] && min=$second
    let sum+=second
done

echo -e "leader cpu: $leaderCPU   \c"
echo -e "follower cpu: \c"
echo "ibase=10; scale=2; ($sum-$leaderCPU)/($#-1)" | bc

echo -e "leader/follower: \c"
echo "ibase=10; scale=2; $leaderCPU/(($sum-$leaderCPU)/($#-1))" | bc
