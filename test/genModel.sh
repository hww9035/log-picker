#!/usr/bin/env bash

#Usage：
# 默认当前目录： ./genModel.sh db table
# 指定目录：./genModel.sh db table ./mysqlDatabase/tango

if [ $# -lt 2 ]; then
  echo "Usage:"
  echo -e "\t"'./genModel.sh ${database} ${table}'
  echo -e "\t"'./genModel.sh ${database} ${table} ${path}'
  exit 1
fi

modelPkgName="model"
dbname=$1
tables=$2
base_dir=$(cd $(dirname $0) && pwd)
model_dir=${base_dir}

model_path=$3
if [ "$model_path" ]; then
  model_dir=$model_path
fi


# 数据库配置
host="127.0.0.1"
port=3306
username="root"
passwd="12345"


echo "开始创建库：$dbname 表：$2 outPath: $model_dir"
gentool -dsn "${username}:${passwd}@tcp(${host}:${port})/${dbname}?charset=utf8mb4&parseTime=True&loc=Local" \
-db mysql \
-tables "${tables}" \
-modelPkgName="${modelPkgName}" \
#-onlyModel \
-outFile gen.go \
-outPath="${model_dir}"
