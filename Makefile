.PHONY: build-mailman build-postgres-dbinit \
minikube minikube-clean minikube-start minikube-stop minikube-tunnel \
postgres postgres-clean postgres-reset \
minikube-build-postgres-dbinit \
postgres-config postgres-storage postgres-deployment postgres-service \
postgres-service-clean postgres-deployment-clean postgres-storage-clean postgres-config-clean \
mailman mailman-clean mailman-rebuild \
minikube-build-mailman \
mailman-config mailman-deployment mailman-service \
mailman-config-clean mailman-deployment-clean mailman-service-clean


BUILD_VERSION ?= dev

#
# Local image builds
#

build-mailman:
	@docker build -t mailman:$(BUILD_VERSION)

build-postgres-dbinit:
	@docker build -t postgres-dbinit:$(BUILD_VERSION) db


#
# Minikube commands
#

minikube: minikube-start postgres mailman

minikube-clean: mailman-clean postgres-clean minikube-stop

minikube-start:
	@minikube start

minikube-stop:
	@minikube stop

# Start minikube tunnel, which allows for connecting to the mailman service on localhost from the host machine
minikube-tunnel:
	@minikube tunnel


#
# Mailman
#

mailman: minikube-build-mailman mailman-config mailman-deployment mailman-service
mailman-clean: mailman-service-clean mailman-deployment-clean mailman-config-clean
# Rebuild and redeploy mailman, update configuration
mailman-rebuild: minikube-build-mailman mailman-deployment-clean mailman-config-clean mailman-config mailman-deployment

minikube-build-mailman:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t mailman:$(BUILD_VERSION) .; \
	}

MAILMAN_YAMLS_DIR=deployments/kubernetes/local/mailman
MAILMAN_CONFIG=$(MAILMAN_YAMLS_DIR)/config.yaml
MAILMAN_DEPLOYMENT=$(MAILMAN_YAMLS_DIR)/deployment.yaml
MAILMAN_SERVICE=$(MAILMAN_YAMLS_DIR)/service.yaml

mailman-config:
	@kubectl create -f $(MAILMAN_CONFIG)

mailman-config-clean:
	@kubectl delete -f $(MAILMAN_CONFIG) --ignore-not-found

mailman-deployment:
	@kubectl create -f $(MAILMAN_DEPLOYMENT)

mailman-deployment-clean:
	@kubectl delete -f $(MAILMAN_DEPLOYMENT) --ignore-not-found

mailman-service:
	@kubectl create -f $(MAILMAN_SERVICE)

mailman-service-clean:
	@kubectl delete -f $(MAILMAN_SERVICE) --ignore-not-found


#
# Postgres
#

postgres: minikube-build-postgres-dbinit postgres-config postgres-storage postgres-deployment postgres-service
postgres-clean: postgres-service-clean postgres-deployment-clean postgres-storage-clean postgres-config-clean
# Gracefully delete DB content and reinitialize the schema
postgres-reset: postgres-deployment-clean postgres-clear-volume postgres-deployment

minikube-build-postgres-dbinit:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t postgres-dbinit:$(BUILD_VERSION) db; \
	}

POSTGRES_YAMLS_DIR=deployments/kubernetes/local/postgres
POSTGRES_CONFIG=$(POSTGRES_YAMLS_DIR)/config.yaml
POSTGRES_DEPLOYMENT=$(POSTGRES_YAMLS_DIR)/deployment.yaml
POSTGRES_SERVICE=$(POSTGRES_YAMLS_DIR)/service.yaml
POSTGRES_STORAGE=$(POSTGRES_YAMLS_DIR)/storage.yaml
POSTGRES_CLEANER_DEPLOYMENT=$(POSTGRES_YAMLS_DIR)/cleaner-deployment.yaml

postgres-config:
	@kubectl create -f $(POSTGRES_CONFIG)

postgres-config-clean:
	@kubectl delete -f $(POSTGRES_CONFIG) --ignore-not-found

postgres-storage:
	@kubectl create -f $(POSTGRES_STORAGE)

postgres-storage-clean:
	@kubectl delete -f $(POSTGRES_STORAGE) --ignore-not-found

postgres-deployment:
	@kubectl create -f $(POSTGRES_DEPLOYMENT)

postgres-deployment-clean:
	@kubectl delete -f $(POSTGRES_DEPLOYMENT) --ignore-not-found

postgres-service:
	@kubectl create -f $(POSTGRES_SERVICE)

postgres-service-clean:
	@kubectl delete -f $(POSTGRES_SERVICE) --ignore-not-found

# Delete DB content
postgres-clear-volume:
	@kubectl delete -f $(POSTGRES_CLEANER_DEPLOYMENT) --ignore-not-found
	@kubectl create -f $(POSTGRES_CLEANER_DEPLOYMENT)
	@kubectl wait --for condition=complete -f $(POSTGRES_CLEANER_DEPLOYMENT)
