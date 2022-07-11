{-# LANGUAGE OverloadedStrings #-}

import Database.MySQL.Simple

main :: IO ()
main = do
  conn <- connect defaultConnectInfo {connectHost = "127.0.0.1", connectUser = "user", connectPassword = "password", connectDatabase = "rdb"}
  [Only i] <- query_ conn "select * from rdb.users"
  print (i :: Int)
