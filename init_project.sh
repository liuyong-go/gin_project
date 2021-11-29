#!/bin/bash
help(){
	echo "使用方式"
	echo "https://raw.githubusercontent.com/liuyong-go/gin_project/main/init_project.sh mod_name dir"
	echo "eg https://raw.githubusercontent.com/liuyong-go/gin_project/main/init_project.sh 'github.com/liuyong-go' new_project" 
	exit 0
}
case $1 in
	-h)
	help
	;;
esac

if [ $# -lt 2 ];then
	echo "参数数量不对，请输入 -h寻求帮助"
	exit 0
fi
#modname=${1//\//\\/}
modname=${1}
modnameparse=${modname//\//\\/}
dir=${2}
echo $modname
echo $dir
# if [ ! -d $dir ];then
# 	mkdir $dir
# fi
echo "拉取代码"
git clone https://github.com/liuyong-go/gin_project.git $dir

cd $dir

rm -rf .git
echo "s/github\.com\/liuyong-go\/gin_project/${modname}/g"

grep -rl "github.com/liuyong-go/gin_project" . --exclude=*sh| xargs  perl -pi -e "s/github\.com\/liuyong-go\/gin_project/${modnameparse}/g"

go mod init ${modname}
go mod tidy
gitpath=${3}
if [ -z "$gitpath" ];then
	echo "done"
	exit 0
fi
echo "初始化git"


git init
git add *
git commit -m '初始化'
git remote add origin $gitpath
git push -u origin master


echo "done"
