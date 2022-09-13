#!bin/bash

MYSQL_PWD=password
DATABASE=rdb
DIR=`dirname $0`

export MYSQL_PWD

# TODO: ユーザー作成

# TODO: Database作成
echo Create database
mysql -u root -e "CREATE DATABASE IF NOT EXISTS ${DATABASE};"

# usersテーブル作成
echo Create users table
mysql -u root ${DATABASE} < ${DIR}/sql/users.ddl

# tasksテーブル作成
echo Create tasks table
mysql -u root ${DATABASE} < ${DIR}/sql/tasks.ddl

# ダミーデータ作成
DML=`find ${DIR}/sql -name "*.dml"`
for D in ${DML}
do
    echo Insert dml file: `basename ${D}`
    mysql -u root ${DATABASE} < ${D}
done
