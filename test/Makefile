db:
	# make sure to create an empty database called "pgtalk"
	#
	docker run --name pgtalk-db -e POSTGRES_PASSWORD=pgtalk  -p 7432:5432 -d postgres

test:
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:7432/postgres go test -v

gen:
	cd ../cmd/pgtalk-gen && go install
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:7432/postgres pgtalk-gen -v -o . -include things,products,categories
	gofmt -w .