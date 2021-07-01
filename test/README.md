# pgtalk/test

## running the db tests

    docker run --name pgtalk-db -e POSTGRES_PASSWORD=pgtalk  -p 5432:5432 -d postgres