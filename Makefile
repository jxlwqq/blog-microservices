redeployed-at:=$(shell date +%s)

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing PGV can currently only be done from source. See: https://github.com/envoyproxy/protoc-gen-validate#installation"

.PHONY: protoc
protoc:
	for file in $$(find api -name '*.proto'); do \
		protoc \
		-I $$(dirname $$file) \
		-I ./third_party \
		--go_out=:$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=:$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$$(dirname $$file)" --validate_opt=paths=source_relative \
		--grpc-gateway_out=:$$(dirname $$file) --grpc-gateway_opt=paths=source_relative \
		$$file; \
	done

.PHONY: wire
wire:
	wire ./...

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations/user -database "mysql://root:@tcp(localhost:3306)/users" -verbose up
	migrate -path ./migrations/post -database "mysql://root:@tcp(localhost:3306)/posts" -verbose up
	migrate -path ./migrations/comment -database "mysql://root:@tcp(localhost:3306)/comments" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations/user -database "mysql://root:@tcp(localhost:3306)/users" -verbose down
	migrate -path ./migrations/post -database "mysql://root:@tcp(localhost:3306)/posts" -verbose down
	migrate -path ./migrations/comment -database "mysql://root:@tcp(localhost:3306)/comments" -verbose down


.PHONY: blog-server
blog-server:
	go run ./cmd/blog/

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

.PHONY: dmt-server
dmt-server:
	@echo "Quick start DTM server. See: https://github.com/dtm-labs/dtm"

.PHONY: docker-build
docker-build:
	docker build -t blog/blog-server:latest -f ./build/docker/blog/Dockerfile .
	docker build -t blog/user-server:latest -f ./build/docker/user/Dockerfile .
	docker build -t blog/auth-server:latest -f ./build/docker/auth/Dockerfile .
	docker build -t blog/post-server:latest -f ./build/docker/post/Dockerfile .
	docker build -t blog/comment-server:latest -f ./build/docker/comment/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -f ./deployments/
	kubectl apply -f ./deployments/dtm/
	kubectl apply -f ./deployments/blog/
	kubectl apply -f ./deployments/user/
	kubectl apply -f ./deployments/post/
	kubectl apply -f ./deployments/auth/
	kubectl apply -f ./deployments/comment/
	kubectl apply -f ./deployments/addons/

.PHONY: kube-delete
kube-delete:
	kubectl delete -f ./deployments/
	kubectl delete -f ./deployments/dtm/
	kubectl delete -f ./deployments/blog/
	kubectl delete -f ./deployments/user/
	kubectl delete -f ./deployments/post/
	kubectl delete -f ./deployments/auth/
	kubectl delete -f ./deployments/comment/
	kubectl apply -f ./deployments/addons/

.PHONY: kube-redeploy
kube-redeploy:
	@echo "redeployed at ${redeployed-at}"
	kubectl patch deployment blog-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment user-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment auth-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment post-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment comment-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'