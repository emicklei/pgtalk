rmdb:
	docker rm -f pgtalk-db

dbsrv:
	docker run --name pgtalk-db -e POSTGRES_PASSWORD=pgtalk  -p 7432:5432 -d postgres

db:
	docker exec -it pgtalk-db createdb -U postgres pgtalk

gen:
	cd ../cmd/pgtalk-gen && go install
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:7432/postgres pgtalk-gen -v -o . -include things,products,categories -mapping types/mapping.json
	# gofmt -w .

test:
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:7432/postgres go test ${FLAGS}

