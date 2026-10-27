.PHONY: run test verify terraform container
run:
	go run ./cmd/service
test:
	go test -race ./...
terraform:
	terraform -chdir=infra fmt -check
	terraform -chdir=infra init -backend=false
	terraform -chdir=infra validate
verify: test terraform
container:
	docker build -t platform-blueprint .
