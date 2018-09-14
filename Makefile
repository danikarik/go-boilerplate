all: key cert assets build

help: ## Show usage
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

key: ## Generate private key
	@openssl genrsa -out certs/localhost.key 2048

cert: key ## Generate SSL certificate
	@openssl req -new -x509 -key certs/localhost.key -out certs/localhost.cert -days 3650 -subj /CN=localhost

assets: ## Compile assets
	@npm run prod

build: ## Compile binary
	@packr build -o bin/app

run: ## Build and run application
	@bin/app

watch: ## Watch static files
	@npm run watch

dev: ## Hot reload
	@refresh run
