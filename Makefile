# ---
# https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects
# ---

# include .env
# include .env-secret
#
app-name := rosso-go-hetzner
app-buildnr := $(shell date +%s)
build-day   := $(shell date +%Y-%m-%d)
go-bin := $(shell go env GOBIN)
pg_user := ${PGUSER}
# PGPORT := ${PGPORT}

help:           ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

info: ## show the acual values
	@echo @@@ info
	@echo build-day $(build-day)
	@echo app-name: $(app-name)
	@echo go-bin: $(go-bin)
	@echo app-buildnr: $(app-buildnr)
	@echo pg_user: $(pg_user)
	@echo PGHOST: ${PGHOST}
	@echo PGPORT: ${PGPORT}
	@echo NO_COLOR: ${NO_COLOR}
	@echo BOOTSTRAP_VERSION: ${BOOTSTRAP_VERSION}
	@echo HTMX_VERSION: ${HTMX_VERSION}

clean: ## clean
	@echo "@@@ clean"
	@go clean
	@rm -rf static/fonts
	@rm -rf static/css/bootstrap*.css
	@rm -f static/js/*.js

setup_npm: ## setup node static files
	@echo "@@@ setup node"
	@mkdir static || true
	@mkdir static/js || true
	@mkdir static/css || true
	@mkdir static/css/fonts || true
	@cd static && \
		npm install bootstrap@${BOOTSTRAP_VERSION} && \
		cp node_modules/bootstrap/dist/js/bootstrap.bundle.js js && \
		cp node_modules/bootstrap/dist/css/bootstrap.css css
	@cd static && \
		npm install htmx.org@${HTMX_VERSION} && \
	  	cp node_modules/htmx.org/dist/htmx.js js && \
	  	cp node_modules/htmx.org/dist/ext/json-enc.js js
	@cd static && \
		npm i htmx-ext-form-json && \
		cp node_modules/htmx-ext-form-json/form-json.js js
	@cd static && \
		npm i bootstrap-icons@1.11.3 && \
		cp node_modules/bootstrap-icons/font/bootstrap-icons.css css/ && \
		mkdir fonts || true && \
		cp node_modules/bootstrap-icons/font/fonts/bootstrap-icons.woff2 css/fonts/ && \
		cp node_modules/bootstrap-icons/font/fonts/bootstrap-icons.woff css/fonts/
	@rm -rf static/node_modules
	@rm -f static/package*
	@echo done

setup_go: ## setup node static files
	@echo "@@@ setup go"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install golang.org/x/tools/gopls@latest

# woodpecker: ## update woodpecker secrts etc
#   @echo woodpecker secret
#   @woodpecker-cli repo secret rm \
#     -repository gitea_admin/rosso0815-hetzner \
#     -name ssh_key || true
#   @echo $(shell echo  "$$GIT_SSH_KEY" > ./test )
#   @woodpecker-cli repo secret add \
#     -repository gitea_admin/rosso0815-hetzner \
#     -name ssh_key \
#     -value @test
#   @rm test
#   @woodpecker-cli repo secret ls --repository gitea_admin/rosso0815-hetzner

# ansible:
#   @ansible-playbook \
#     --inventory playbooks/inventory/hosts \
#     --extra-vars my_version=1.0.6 \
#     --private-key ~/.ssh/id_ed25519 \
#     --user rosso0815 --connection ssh \
#     --ssh-common-args '-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \
#     playbooks/deploy-goapp.yml

update: clean lstatic
	@echo "@@@ update"
	@echo "drop table users cascade;" | psql || true
	@echo "drop table customers cascade;" | psql || true

audit: ## Quality Check
	# @yamllint .
	# @go mod tidy -diff
	# @go mod verify
	# @test -z "$(shell gofmt -l .)"
	# @go vet ./...
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	# @go run golang.org/x/vuln/cmd/govulncheck@latest ./...

sql_dump: ## dump actual database
	# 	FIXME add sequeces
	# 	https://dba.stackexchange.com/questions/294212/export-postgres-sequences-only
	#	pg_dump -h localhost -p 5432 -d namedb -U postgres -t '*_id_seq' > dump-seq.sql
	#
	@pg_dump -p $(PGPORT) -h $(PGHOST) -U $(PGUSER) -d $(PGNAME) -a --insert -t customer > dump_$(build-day).sql
	@pg_dump -p $(PGPORT) -h $(PGHOST) -U $(PGUSER) -d $(PGNAME) -a --insert -t invoice >> dump_$(build-day).sql
	@pg_dump -p $(PGPORT) -h $(PGHOST) -U $(PGUSER) -d $(PGNAME) -a --insert -t invoiceentry >> dump_$(build-day).sql
	@pg_dump -p $(PGPORT) -h $(PGHOST) -U $(PGUSER) -d $(PGNAME) -a --insert -t '*_seq' >> dump_$(build-day).sql
	@ls -ltr

sql: ## run sql
	@echo "@@@ sql"
	@goose reset
	@goose up
	@echo add customer-data
	@psql < db/data/001_customer.sql
	@echo add invoice-data
	@psql < db/data/002_invoice.sql
	@echo add userkv-data
	@psql < db/data/003_userkv.sql
	@sqlc generate
	@templ generate
	@go build
	@echo "@@@ sql done"

run: run_sql ## run local
	@echo run
	@sqlc generate
	@templ generate
	@air
	# @rm -rf tmp
	# @mkdir tmp
	# @wgo -file=.go -file=.templ -xfile=_templ.go templ generate :: go run main.go web
	# @wgo -cd tmp  -file=.go -file=.templ -xfile=_templ.go templ generate .. :: go build -o main ..  :: ./main web

	# @kill -9 $$(lsof -t -i :3000) || true
	# @rm -rf tmp
	# @templ generate
	# @go build -o tmp/main .
	# @cd tmp && ./main web

build: update ## build all envoronments
	@echo @@@ build
	@cd db && sqlc generate
	@templ generate
	# @go fmt
	# @go vet
	# @go test -race
	@CGO_ENABLED=1 go build -o hetzner_httpd

local_386: update ## build 386 envoronments
	@echo @@@ build
	@sqlc generate
	@templ generate
	@go fmt
	@go vet
	# @go test -race
	# @CGO_ENABLED=1 GOOS=linux GOARCH=386  go build -o hetzner_httpd_386
	@GOOS=linux GOARCH=386  go build -o hetzner_httpd_386

lint: ## lint the stuff
	@go build
	@golangci-lint run
	@go vet

docker_build: ## build all envoronments
	@echo @@@ docker_build
	@docker build . -t test:test

docker_run: docker_build ## run inside docker
	# docker -H master run -p 3000:3000 -d test:test
	@echo @@@ docker_run
	@docker run -it --rm --name test -e GOAPP_WEB_PREFIX=/test -p 8080:8080 test:test

build_gitea: build ## build goes into gitea images
	@echo @@@ build_docker
	@docker build . -t gitea.localnet/goapp:1.0.0

k8s-build: debug ## build a k8s-image
	@nerdctl build --namespace k8s.io -t $(app-name):$(app-buildnr) .
	@nerdctl images --namespace k8s.io
	@ansible-playbook deploy.yml \
		  -e DOCKER_TAG=$(app-name) \
		  -e DOCKER_APP=$(app-name) \
		  -e DOCKER_VERSION=$(app-buildnr)

schema_database: ## build apng from database
	@planter postgres://pfistera@localhost/pfistera?sslmode=disable -o db_schema.puml

# --- EOF
