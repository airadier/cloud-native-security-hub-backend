.PHONY: test build push

test:
	go test -v ./...

dbimport: cmd/dbimport/main.go pkg/dbmodel/dbmodel.go
	go build -o dbimport cmd/dbimport/main.go

dev:
	RESOURCES_PATH=test/fixtures/resources VENDOR_PATH=test/fixtures/vendors go run cmd/server/main.go

db-dev:
	RESOURCES_PATH=test/fixtures/resources VENDOR_PATH=test/fixtures/vendors go run cmd/dbserver/main.go

watch:
	ag -l | entr -c go test -v ./...

build:
	docker build -t gcr.io/mateo-burillo-ns/securityhub-backend .

push: build
	docker push gcr.io/mateo-burillo-ns/securityhub-backend

deploy: deploy-backend deploy-frontend

deploy-backend:
	kubectl -n securityhub patch deployment backend -p "{\"spec\": {\"template\": {\"metadata\": { \"labels\": {  \"redeploy\": \"$(shell date +%s)\"}}}}}"

deploy-frontend:
	kubectl -n securityhub patch deployment frontend -p "{\"spec\": {\"template\": {\"metadata\": { \"labels\": {  \"redeploy\": \"$(shell date +%s)\"}}}}}"
