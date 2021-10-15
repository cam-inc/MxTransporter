include .env

.PHONY: all help \
			build-image

all: help

help:
	@echo ""
	@echo "================================"
	@echo "Usage"
	@echo "================================"
	@echo "$$ make build-image"
	@echo "create transporter docker image"


##########################
# main command
##########################

build-image:
	@echo "build image..."
	docker build -t mxtransporter .
	docker tag mxtransporter $(GCR_REPO):$(GCR_REPO_TAG)
	@echo "push image to gcr..."
	docker push $(GCR_REPO):$(GCR_REPO_TAG)

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

