dbsrv:
	docker run --name pgtalk-srv -e POSTGRES_PASSWORD=pgtalk  -p 5432:5432 -d postgres

db:
	docker exec -it pgtalk-srv createdb -U postgres pgtalk-db

gen:
	cd ../cmd/pgtalk-gen && go install
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:5432/pgtalk-db pgtalk-gen -v -o . -include things,products,categories -mapping types/mapping.json -cache pgtalk-tables-cache.json
	# gofmt -w .

test:
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:5432/pgtalk-db go test ${FLAGS}

rmdb:
	docker rm -f pgtalk-srv