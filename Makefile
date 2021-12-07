.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: wire
wire:
	wire ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: protoc
protoc:
	for file in $$(find api -name '*.proto'); do \
		protoc -I $$(dirname $$file) \
		-I ./third_party \
		--go_out=:$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=:$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$$(dirname $$file)" --validate_opt=paths=source_relative \
		$$file; \
	done

.PHONY: user-server
user-server:
	go run ./cmd/user/

.PHONY: post-server
post-server:
	go run ./cmd/post/

.PHONY: comment-server
comment-server:
	go run ./cmd/comment/

.PHONY: auth-server
auth-server:
	go run ./cmd/auth/

.PHONY: docker-build
docker-build:
	docker build -t blog/user-server:latest -f ./build/docker/user/Dockerfile .
	docker build -t blog/auth-server:latest -f ./build/docker/auth/Dockerfile .
	docker build -t blog/post-server:latest -f ./build/docker/post/Dockerfile .
	docker build -t blog/comment-server:latest -f ./build/docker/comment/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -f ./deployments/
	kubectl apply -f ./deployments/user/
	kubectl apply -f ./deployments/post/
	kubectl apply -f ./deployments/auth/
	kubectl apply -f ./deployments/comment/

.PHONY: kube-delete
kube-delete:
	kubectl delete -f ./deployments/
	kubectl delete -f ./deployments/user/
	kubectl delete -f ./deployments/post/
	kubectl delete -f ./deployments/auth/
	kubectl delete -f ./deployments/comment/


.PHONY: kube-redeploy
kube-redeploy:
	./scripts/kube-redeploy.sh

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: migrate-up
migrate-up:
	./scripts/migrate-up.sh

.PHONY: migrate-down
migrate-down:
	./scripts/migrate-down.sh