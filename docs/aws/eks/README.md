# Build MxTransporter in EKS 

This section describes the procedure for building MxTransporter in the EKS environment.<br>
All commands to create each AWS resources are wrapped in ```Makefile```.

## Prepare
### Command
Need to use following commands. The version listed is the verified version.
```
aws      v2.1.30
docker   v20.10.8
eksstl   v0.70.0
helm     v3.6.3
kubectl  v1.22.1
make     v3.81
```

### Environment variables files
Before starting the construction, create ```.env``` and ```.secrets.env``` in the current directory by referring to ```.env.template``` and ```secrets.env.template```.

If you want to export change streams to Kinesis Data Streams, write the following description in ```.env```.

```
EXPORT_DESTINATION=kinesisStream
```

## Procedure
**1. Create Key pair for node instance (EC2).**

It is used to do ssh to EC2 instance.

<br>

**2. Create KMS key for EKS cluster.**

It is used to encrypt kubernetes secrets.

<br>

**3. Create EKS cluster and node group.**

```
$ make build
```

・Create cluster

Create ```cluster.yaml``` by referring to the environment variables written in .env. Then create a cluster with the ```eksctl create cluster``` command.

・Create node group

Create ```nodegroup.yaml``` by referring to the environment variables written in .env. Then create a node group with the ```eksctl create nodegroup``` command.

<br>

**4. Attach Kinesis policy to node group role.**

To Mxtransporter container export change streams to Kinesis Data Streams, you have to attach Kinesis policy to node group role.<br>
Node group role is created by ```eksctl create cluster``` command, just attach , for example ```AmazonKinesisFullAccess``` policy, to that role.

<br>

**5. Create kubernetes secrets.**

Collect the environment variables written in secrets.env and create them in a cluster as kubernetes secrets.

```
$ make secrets
```

<br>

**6. Deploy kubernetes resources.**

This command creates a StatefulSet, HeadlessService, Horizontal Pod Autoscaler, and PVC.

```
$ make deploy
```

The following processing is performed here.<br>
・build docker image.<br>
・login ecr repository.<br>
・push docker image to ecr repository.<br>
・build helm variables.<br>
・deploy with helm templates.<br>

<br>

**7. Upgrade kubernetes resources.** 

You can upgrade kubernetes resources with the following command.

Note that if you do not update ```ECR_REPO_TAG```, the new docker image will be referenced and the container will not be created.

```
$ make upgrade
```