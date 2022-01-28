include .env

.PHONY: all help \
		build create-service-account create-custome-role attach-role-to-service-account build-gke-variables build-gke create-namespace\
		deploy build-image build-helm-variables deploy-gke-resources \
		secret create-secret \
		stop pods-stop start pods-start \
		upgrade upgrade-gke-resources update-gke \
		create-bigquery-dataset create-bigquery-table set-bigquery-table-expiration-date

all: help

help:
	@echo ""
	@echo "================================"
	@echo "Usage"
	@echo "================================"
	@echo "$$ make build"
	@echo "create GKE cluster & service account & namespace"
	@echo ""
	@echo "$$ make secret"
	@echo "create GKE secret"
	@echo ""
	@echo "$$ make deploy"
	@echo "create GKE resources"
	@echo ""
	@echo "$$ make upgrade"
	@echo "upgrade GKE resources"


##########################
# main command
##########################
build: create-service-account \
		create-custome-role \
		attach-role-to-service-account \
		build-gke-variables \
		build-gke \
		create-namespace

secrets: create-secrets

deploy: build-image \
		puch-to-gcr \
		build-helm-variables \
		deploy-gke-resources

stop: pods-stop

start: pods-start

upgrade: build-image \
         build-helm-variables \
         upgrade-gke-resources

##########################
# sub command
##########################

create-service-account:
	@echo "create service account..."
	gcloud iam service-accounts create mxtransporter --display-name="mxtransporter" \
			--project $(PROJECT_NAME) \
			--description="service account for node instance"

create-custome-role:
	@echo "create custome role..."
	@echo "create bigquery role..."
	gcloud iam roles create "BigqueryForMxtransporter" \
			--project=$(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) \
			--title="Bigquery For mxtransporter" \
			--description="Role that mxtransporter access bigquery resources." \
			--permissions=bigquery.datasets.get,bigquery.tables.get,bigquery.tables.getData,bigquery.tables.list,bigquery.tables.update,bigquery.tables.updateData
	@echo "create pubsub role..."
	gcloud iam roles create "PubSubForMxtransporter" \
			--project=$(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) \
			--title="PubSub For mxtransporter" \
			--description="Role that mxtransporter access pubsub resources." \
			--permissions=pubsub.subscriptions.get,pubsub.subscriptions.list,pubsub.subscriptions.create,pubsub.subscriptions.update,pubsub.topics.get,pubsub.topics.list,pubsub.topics.attachSubscription,pubsub.topics.create,pubsub.topics.publish,pubsub.topics.update

attach-role-to-service-account:
	@echo "attach role to service account in $(PROJECT_NAME)..."
	gcloud projects add-iam-policy-binding $(PROJECT_NAME) \
			--member="serviceAccount:mxtransporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="roles/monitoring.editor"
	gcloud projects add-iam-policy-binding $(PROJECT_NAME) \
			--member="serviceAccount:mxtransporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="roles/logging.admin"
	gcloud projects add-iam-policy-binding $(PROJECT_NAME) \
			--member="serviceAccount:change-streams-exporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="roles/storage.admin"

	@echo "attach role to service accountin $(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) ..."
	gcloud projects add-iam-policy-binding $(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) \
			--member="serviceAccount:mxtransporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="projects/$(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS)/roles/BigqueryForMxtransporter"
	gcloud projects add-iam-policy-binding $(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) \
			--member="serviceAccount:mxtransporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="roles/stackdriver.accounts.editor"
	gcloud projects add-iam-policy-binding $(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS) \
			--member="serviceAccount:mxtransporter@$(PROJECT_NAME).iam.gserviceaccount.com" \
			--role="projects/$(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS)/roles/PubSubForMxtransporter"

build-gke-variables:
	@echo  "create gke variables template..."
	sh ./deployment-manager/gke.sh

build-gke:
	@echo "build gke cluster..."
	gcloud deployment-manager deployments create mxtransporter-cluster \
			--config ./deployment-manager/gke.yaml \
			--project $(PROJECT_NAME)
	gcloud container clusters get-credentials mxtransporter-cluster \
			--project $(PROJECT_NAME)
	@echo "selected context:"
	@kubectl config view | grep current-context

create-namespace:
	@echo "create namespace..."
	kubectl create namespace mxtransporter

update-gke:
	@echo "build gke cluster..."
	gcloud deployment-manager deployments update mxtransporter-cluster \
			--config ./deployment-manager/gke.yaml \
			--project $(PROJECT_NAME)

create-secrets:
	@echo "create secret in kubernetes secrets..."
	@if [ -e secrets.env ]; then \
	  kubectl create secret generic $(KUBERNETES_SECRET_NAME) --from-env-file=secrets.env -n mxtransporter; \
    else \
      echo '"secrets.env" file is not exists. Create secrets.env file in current directory.'; \
    fi

build-image:
	@echo "build image..."
	docker build -t mxtransporter ../../../

puch-to-gcr:
	docker tag mxtransporter $(GCR_REPO):$(GCR_REPO_TAG)
	@echo "push image to gcr..."
	docker push $(GCR_REPO):$(GCR_REPO_TAG)

build-helm-variables:
	@echo "create helm variables template..."
	sh ./helm/values.sh

deploy-gke-resources:
	@echo "deploy gke resources..."
	helm install -f helm/values.yaml -n mxtransporter mxtransporter-for-$(PROJECT_NAME) ./helm

pods-stop:
	@kubectl config view | grep current-context
	@read -p "Are you sure you want to stop the pod? (y/N): " yn; \
	case "$$yn" in \
	  [yY]*) kubectl scale statefulset --all --replicas=0 -n mxtransporter;; \
	  *) echo "Rejected to stop pods in the specified cluster.";; \
	esac

pods-start:
	@echo "start pods..."
	kubectl scale statefulset --all --replicas=1 -n mxtransporter

upgrade-gke-resources:
	@echo "deploy gke resources..."
	helm upgrade -f helm/values.yaml --set-string image.tag=$(GCR_REPO_TAG) -n mxtransporter mxtransporter-for-$(PROJECT_NAME) ./helm


##########################
# bigquery command
##########################

create-bigquery-dataset:
	@echo "create bigquery dataset..."
	@read -p "Specify the dataset name you want to create : " dataset; \
	bq mk --dataset \
		--location=$(BIGQUERY_DATASET_LOCATION) \
		$$dataset

create-bigquery-table:
	@echo "create bigquery table..."
	@read -p "Specify the dataset name you want to create table in : " dataset; \
	read -p "Specify the table name you want to create : " table; \
	bq mk --table \
		--time_partitioning_type=DAY \
		--time_partitioning_field=clusterTime \
		--schema=id:STRING,operationType:STRING,clusterTime:TIMESTAMP,fullDocument:STRING,ns:STRING,documentKey:STRING,updateDescription:STRING \
		$$dataset.$$table

set-bigquery-table-expiration-date:
	@read -p "Specify the dataset name you want to create table in : " dataset; \
	read -p "Specify the table name you want to create : " table; \
	bq update --time_partitioning_expiration $(BIGQUERY_TABLE_PARTITIONING_EXPIRATION_TIME) $(PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS):$$dataset.$$table