#!/bin/bash
app_name="ceshi"
app_path="/Users/liuyong/Downloads/demo/"
build_path="/Users/liuyong/Downloads/demo/"
git_path="https://github.com/liuyong-go/gin_project.git"
nginx_conf="/Users/liuyong/Downloads/www.gotest.com.conf"
listen_adress="0.0.0.0"
port=(58322 58323)
# for p in ${port[*]}
# do
#     echo "p is $p"
# done

#app_name ps进程获取当前执行进程名

ps aux | grep php | awk '{print $2,$11}' > /tmp/deploy_$app_name
i=0
while read a b
do
    i=`expr $i + 1`
    pids[$i]=$a
    appnames[$i]=$b
done < /tmp/deploy_$app_name
#echo ${pids[@]}
appnames[0]="ceshi_58322"
#echo ${appnames[@]}
#needkill

for p in ${port[*]}
do
    index=0
    for am in ${appnames[*]}
    do
        if [ ${app_name}"_"${p} == $am ];then
            needkillport[${#needkillport[*]}]=$p
            needkillpid[${#needkillpid[*]}]=${pids[$index]}
            break
        fi
        index=$index + 1
    done
done
killpid=0
killport=0
runport=0
if [ ${#needkillpid[*]} -gt 1 ];then
    echo "need kill count gt 1"
    exit 0
fi
if [ ${#needkillpid[*]} == 1 ];then
    killpid=${needkillpid[0]}
    killport=${needkillport[0]}
fi
for p in ${port[*]}
do
    if [ $p != $killpid ];then
        runpid=$p
        break
    fi
done
echo "killpid: "$killpid
echo "killport: "$killport
echo "run: "$runpid
#判断引用目录，不存在git clone,存在 pull
if [ ! -d $app_path ];then
    git clone $git_path $app_path
    cd $app_path
else
    cd $app_path
    git pull
fi
#到应用目录执行go build -o app_name_{数组中未运行中端口号}，启动进程
go mod tidy
go build -o ${build_path}${app_name}"_"${runpid}
sudo chmod 755 ${build_path}${app_name}"_"${runpid}
nohup ${build_path}${app_name}"_"${runpid} -l ${listen_adress}":"${runpid} &
#修改nginx配置文件映射端口号 sudo nginx -s reload

#5分钟后kill掉之前进程



