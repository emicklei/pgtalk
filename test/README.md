# pgtalk/test

## start local postgres db

    make dbsrv
    
## create database

    make db

## running the db tests

The test will create the tables if missing.

    make test