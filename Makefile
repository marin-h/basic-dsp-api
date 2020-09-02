swagger.validate:
	docker run --rm -it -e GOPATH=${HOME}/go:/go -v ${HOME}:${HOME} -w ${HOME}/simple-dsp quay.io/goswagger/swagger validate api/swagger/swagger.yml

swagger.doc:
	docker run -i yousan/swagger-yaml-to-html < api/swagger/swagger.yml > doc/index.html

start:
	go run .

benchmarks:
	go test ./benchmark/ -bench=. -v

tests:
	go test ./test/ -v