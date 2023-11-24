test:
	go test -race -coverprofile="coverage.out" -covermode=atomic ./...
	go tool cover -html="coverage.out"

lint:
	golangci-lint run

BIN:=$(CURDIR)/bin
install-deps:
	GOBIN=$(BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


#generate-note-api:
#	mkdir -p pkg/auth_v1  # Создаем куда всё сгенеренное положим
#	protoc --proto_path api/auth_v1 \  # Вызываем проток и говорим протоку гдеу нас находиться прото файл, точка входа
#	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \  # параметр который говорит, что мы хотим сгенерить всё что есть в протофайле, в гошные структуры, и этим всем занимается бинарник protoc-gen-go
#	--plugin=protoc-gen-go=bin/protoc-gen-go \ 				# говорим что бинарники нужно искать в локальном bin а не в том который у вас по умолчанию
#	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
#	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
#	api/auth_v1/auth.proto   				# прямой путь дор прото файла


gen-auth: ## Генерация proto-файлов
		protoc 	--proto_path=api/auth_v1 \
				--proto_path=proto \
				--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
				--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
				--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
				--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
				api/auth_v1/auth.proto

install:
	GOBIN=$(BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	GOBIN=$(BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1


gen:
	protoc 	--proto_path=api/admin_v1 \
			--proto_path=proto \
			--go_out=pkg/admin_v1 --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
			--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
			--go-grpc_out=pkg/admin_v1 --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway.exe \
			--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2.exe \
			--grpc-gateway_out=pkg/admin_v1 --grpc-gateway_opt=paths=source_relative \
			--validate_out lang=go:pkg/admin_v1 --validate_opt=paths=source_relative \
			--plugin=protoc-gen-validate=bin/protoc-gen-validate.exe \
			--openapiv2_out=allow_merge=true,merge_file_name=api_admin_v1:docs \
			--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2.exe \
			api/admin_v1/admin.proto