RAS_MONITORING_IMAGE ?= "artifactory-kfs.habana-labs.com/k8s-infra-docker-dev/habana/ras-monitoring"
RAS_MONITORING_VERSION ?= 2.0.9
ENV ?= prod
tidy:
	@go mod tidy

test/code:
	@go test -cover ./...

test:  ## test the habana_bmc_exporter image
	@docker rm -f habana_bmc_exporter-test > /dev/null 2>&1
	@echo start the habana_bmc_exporter container in detached mode
	@docker run -itd --name habana_bmc_exporter-test -p 5000:5000 -v `pwd`:/tmp $(RAS_MONITORING_IMAGE):$(RAS_MONITORING_VERSION) -config /tmp/config-test.json -exporter g3-red-fish && \
	sleep 10s && \
	curl localhost:5000/debug/readiness
	curl localhost:5000/metrics
	@docker rm -f habana_bmc_exporter-test

mock:
	@which mockgen || go install go.uber.org/mock/mockgen@latest
	@mockgen -destination pkg/mock/exporter.go -package mock -source pkg/bmc-monitoring/bmc-monitoring.go Exporter

build/bin:
	@go build  -o bin/habana_bmc_exporter ./app/services/habana_bmc_exporter

run: build/bin
	@./bin/habana_bmc_exporter

## build: build habana_bmc_exporter docker image
.PHONY: build
build:
	docker build \
	-f zarf/docker/dockerfile.habana_bmc_exporter \
		-t $(RAS_MONITORING_IMAGE):$(RAS_MONITORING_VERSION) \
		--build-arg BUILD_VERSION=$(RAS_MONITORING_VERSION) \
		--build-arg ENV=$(ENV) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

## push/habana_bmc_exporter: push the habana_bmc_exporter to artifactory
.PHONY: push/habana_bmc_exporter
push/habana_bmc_exporter: build
	@docker push $(RAS_MONITORING_IMAGE):$(RAS_MONITORING_VERSION)

## kustomize/habana_bmc_exporter: build habana_bmc_exporter deployment yaml dynamically, and print to screen
.PHONY: kustomize/habana_bmc_exporter
kustomize/habana_bmc_exporter:
	@cd zarf/k8s; kustomize edit set image habana_bmc_exporter_image=$(RAS_MONITORING_IMAGE):$(RAS_MONITORING_VERSION)
	@kustomize build zarf/k8s


# prepare-dashboards will prepare the dashboards for provisionning.
# it will remove the variable values that we had when saving the dashboard.
# it will wrap the dashbord json in the dahboard object that we need for provisioning.
prepare-dashboards:
	@for file in zarf/grafana/dashboards/*.json; do \
	cat $$file | jq '.templating.list[].current.text=""' | jq '.templating.list[].current.value=""' | jq 'del(.id)' > tmp.json && \
	cat tmp.json | jq '{"dashboard": .}'  > $$file ; \
	done
	@rm tmp.json


coverage: ## Run the tests of the project and export the coverage
	go test -coverpkg=./... -cover -covermode=count -coverprofile=cover.out.tmp `go list ./... | grep -v ./pkg/mock`
	cat cover.out.tmp | grep -v "exporter.go" > profile.cov
	go tool cover -func profile.cov


