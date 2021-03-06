.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o helm-broker main.go
	docker build  . -t polskikiel/hb-test:history2

.PHONY: deploy
deploy: build
	docker push polskikiel/hb-test:history2