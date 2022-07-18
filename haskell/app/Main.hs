{-# LANGUAGE OverloadedStrings #-}

import Database.MySQL.Simple

main = do
  conn <- connect defaultConnectInfo {connectHost = "127.0.0.1", connectPort = 3306, connectUser = "user", connectPassword = "password", connectDatabase = "rdb"}
  [Only i] <- ping conn
  [Only i2] <- insert conn
  print i
  print i2

ping :: Connection -> IO [Only Int]
ping conn = query_ conn "select 2 + 2"

insert :: Connection -> IO [Only Int]
insert conn = query_ conn "INSERT INTO users (id) VALUES (`1`), (`2`), (`3`)"
