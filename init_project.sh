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
#git clone https://github.com/liuyong-go/gin_project.git $dir

cd $dir

rm -rf .git
echo "s/github\.com\/liuyong-go\/gin_project/${modname}/g"

grep -rl "github.com/liuyong-go/gin_project" . --exclude=*sh| xargs  perl -pi -e "s/github\.com\/liuyong-go\/gin_project/${modnameparse}/g"

go mod init ${modname}
go mod tidy
echo "done"
echo " 对新建的文件夹进行初始化：  git init，之后该项目下会出现一个  .git 的隐藏文件，该项目变成本地仓库。"
echo "添加远端仓库：git remote add opstech 远端仓库地址

   添加之后，可以用  git remote -v 进行查看 ，然后git  commit -m "xx"
"
echo "git push --set-upstream xxxx"



