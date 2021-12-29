db:
	docker run --name pgtalk-db -e POSTGRES_PASSWORD=pgtalk  -p 5432:5432 -d postgres

test:
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:5432/pgtalk gotip test -v -gcflags=-G=3

gen:
	cd ../cmd/pgtalk-gen && go install
	PGTALK_CONN=postgres://postgres:pgtalk@localhost:5432/pgtalk pgtalk-gen -o .
	# TODO
	/Users/emicklei/sdk/gotip/bin/gofmt -w .