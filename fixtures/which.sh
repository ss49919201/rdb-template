#!/bin/bash

MYSQL_PWD=password
DATABASE=rdb

export MYSQL_PWD

echo Which general_log file
mysql -u root -e "SHOW VARIABLES LIKE 'general_log%'"
