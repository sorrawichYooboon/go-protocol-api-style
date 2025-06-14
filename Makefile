run-dev:
	go run main.go

generate-random-key:
	openssl rand -base64 32

docker-build:
	docker build -t memoviz-user-service .

docker-run-dev:
	docker run -d -p 8081:8081 -e APP_ENV=dev memoviz-user-service

docker-gcloud-build:
	docker build -t <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO_NAME>/<IMAGE_NAME>:<TAG> .

docker-gcloud-build-example:
	docker build -t asia-southeast1-docker.pkg.dev/memoviz-staging/memoviz-docker/user-service:0 .

docker-gcloud-push:
	docker push <REGION>-docker.pkg.dev/<PROJECT_ID>/<REPO_NAME>/<IMAGE_NAME>:<TAG>

docker-gcloud-push-example:
	docker push asia-southeast1-docker.pkg.dev/memoviz-staging/memoviz-docker/user-service:0

docker-log:
	docker logs -f <container_id>

docker-compose-up:
	docker-compose up -d

generate-api-key:
	openssl rand -hex 32

generate-unix-timestamp:
	date +%s

migrate-version:
	migrate -path ./migrations -database "postgres://memoviz_user_service_user:memoviz_user_service_password@localhost:5432/memoviz_user_service_db?sslmode=disable" version

migrate-version-placeholders:
	migrate -path ./migrations -database "postgres://$(user):$(password)@$(host):$(port)/$(dbname)?sslmode=disable" version

migrate-force-version:
	migrate -path ./migrations -database "postgres://memoviz_user_service_user:memmemoviz_user_service_password@localhost:5432/memoviz_user_service_db?sslmode=disable" force $(version)

migrate-up:
	migrate -path ./migrations -database "postgres://memoviz_user_service_user:memoviz_user_service_password@localhost:5432/memoviz_user_service_db?sslmode=disable" up

rate-limit-test:
	i=1
    while [ $i -le 105 ]; do
            curl -X GET http://localhost:8081/health                      
            ((i++))
    done

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-steps-up-one:
	go run cmd/migrate/main.go steps 1

migrate-steps-down-one:
	go run cmd/migrate/main.go steps -1