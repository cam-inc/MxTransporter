include .env

.PHONY: all help \
			create-kms-key \
			build create-cluster create-nodegroup \
			secret create-secret \
			deploy build-image login-ecr push-to-ecr build-helm-variables deploy-eks-resources

help:
	@echo ""
	@echo "================================"
	@echo "Usage"
	@echo "================================"
	@echo "$$ make build"
	@echo "create EKS cluster & node group & namespace"
	@echo ""
	@echo "$$ make secret"
	@echo "create EKS secret"
	@echo ""
	@echo "$$ make deploy"
	@echo "create EKS resources"
	@echo ""
	@echo "$$ make upgrade"
	@echo "upgrade EKS resources"
	@echo ""


##########################
# main command
##########################
kms: create-kms-key

build: build-eks-variables build-cluster \
			create-namespace \
			build-eks-node-group-variables build-nodegroup

secrets: create-secrets

deploy: build-image \
			login-ecr \
			push-to-ecr \
			build-helm-variables \
			deploy-eks-resources

stop: pods-stop

start: pods-start


upgrade: build-image \
			login-ecr \
			push-to-ecr \
			build-helm-variables \
			upgrade-eks-resources


##########################
# sub command
##########################
create-kms-key:
	aws kms create-alias --alias-name alias/mxtransporter-key-for-eks-cluster --target-key-id `aws kms create-key | jq -r '.KeyMetadata.KeyId'`

build-eks-variables:
	@echo  "create eks variables template..."
	sh ./cluster.sh

build-cluster:
	@echo "create cluster..."
	eksctl create cluster -f ./cluster.yaml --profile $(AWS_PROFILE)
	eksctl utils write-kubeconfig --cluster=mxtransporter-cluster --profile $(AWS_PROFILE)
	@echo "selected context:"
	@kubectl config view | grep current-context

create-namespace:
	@echo "create namespace..."
	kubectl create namespace mxtransporter

build-eks-node-group-variables:
	@echo  "create eks node group variables template..."
	sh ./nodegroup.sh

build-nodegroup:
	@echo "create nodegroup"
	eksctl create nodegroup -f ./nodegroup.yaml --profile $(AWS_PROFILE)

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

login-ecr:
	@echo "login ecr..."
	read -p "Specify aws region name : " region; \
	read -p "Specify aws account id : " accountId; \
	aws ecr get-login-password \
		--profile $(AWS_PROFILE) \
		--region $$region \
		| \
		docker login --username AWS \
		--password-stdin $$accountId.dkr.ecr.ap-northeast-1.amazonaws.com

push-to-ecr:
	@echo "Warning: You have to authenticate to ECR"
	docker tag mxtransporter $(ECR_REPO):$(ECR_REPO_TAG)
	@echo "push image to ecr ..."
	docker push $(ECR_REPO):$(ECR_REPO_TAG)

build-helm-variables:
	@echo "create helm variables template..."
	sh ./helm/values.sh

deploy-eks-resources:
	@echo "deploy eks resources..."
	helm install -f helm/values.yaml -n mxtransporter mxtransporter-for-$(AWS_ACCOUNT_NAME) ./helm

pods-stop:
	@kubectl config view | grep current-context
	@read -p "Are you sure you want to stop the pod? (y/N): " yn; \
	case "$$yn" in \
	  [yY]*) kubectl scale statefulset --all --replicas=0 -n change-streams-exporter;; \
	  *) echo "Rejected to stop pods in the specified cluster.";; \
	esac

pods-start:
	@echo "start pods..."
	kubectl scale statefulset --all --replicas=1 -n change-streams-exporter

upgrade-eks-resources:
	@echo "deploy eks resources..."
	helm upgrade -f helm/values.yaml --set-string image.tag=$(ECR_REPO_TAG) -n mxtransporter mxtransporter-for-$(AWS_ACCOUNT_NAME) ./helm