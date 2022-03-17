#!/bin/bash
app_name="ceshi"
app_path="/Users/liuyong/Downloads/demo/"
build_path="/Users/liuyong/Downloads/demo/"
git_path="https://github.com/liuyong-go/gin_project.git"
nginx_conf="/Users/liuyong/Downloads/www.gotest.com.conf"
nginx_pattern="localhost:"
listen_adress="0.0.0.0"
port=(58322 58323)
# for p in ${port[*]}
# do
#     echo "p is $p"
# done
#app_name ps进程获取当前执行进程名

ps aux | grep $app_name | awk '{print $2,$11}' > /tmp/deploy_$app_name
i=0
while read a b
do
    pids[$i]=$a
    appnames[$i]=$b
    i=`expr $i + 1`
done < /tmp/deploy_$app_name
echo "pids: "${pids[@]}
#appnames[0]="ceshi_58322"
#echo ${appnames[@]}
#needkill

for p in ${port[*]}
do
    index=0
    for am in ${appnames[*]}
    do
        if [ ${build_path}${app_name}"_"${p} == $am ];then
            needkillport[${#needkillport[*]}]=$p
            needkillpid[${#needkillpid[*]}]=${pids[$index]}
            break
        fi
        index=`expr $index + 1`
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
    if [ $p != $killport ];then
        runport=$p
        break
    fi
done
echo "killpid: "$killpid
echo "killport: "$killport
echo "run: "$runport
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
go build -o ${build_path}${app_name}"_"${runport}
sudo chmod 755 ${build_path}${app_name}"_"${runport}
echo "执行命令"${build_path}${app_name}"_"${runport} -l=${listen_adress}":"${runport}
nohup ${build_path}${app_name}"_"${runport} -l=${listen_adress}":"${runport} &
#修改nginx配置文件映射端口号 sudo nginx -s reload
echo "修改nginx配置"
grep -rn $nginx_pattern $nginx_conf 
sys=`uname  -a`
if [[ $sys =~ "Darwin" ]];then
    sed -i '' "s/${nginx_pattern}${killport}/${nginx_pattern}${runport}/g" $nginx_conf
else
    sed -i  "s/${nginx_pattern}${killport}/${nginx_pattern}${runport}/g" $nginx_conf
fi
grep -rn $nginx_pattern $nginx_conf 
#sudo nginx -s reload
echo "重启nginx完成"
#5分钟后kill掉之前进程
echo "休眠1分钟"
sleep 60
kill -9 $killpid


