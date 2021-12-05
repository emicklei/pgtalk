# pgtalk/test

## start local postgres db

    make db

## create database

Use your tool to create the databse "pgtalk" in the public schema.

## running the db tests

The test will create the tables if missing.

    make test