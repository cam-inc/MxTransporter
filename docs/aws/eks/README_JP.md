# MxTransporter を EKS 環境で構築する

ここでは MxTransporter を EKS 環境で構築するための手順を紹介します。
各 AWS リソースを作成するコマンドはすべて ```Makefike``` でラップされています。

## 準備
### コマンド
以下のコマンドが使用できる必要があります。記載のバージョンは検証済みです。

```
aws      v2.1.30
docker   v20.10.8
eksctl   v0.70.0
helm     v3.6.3
kubectl  v1.22.1
make     v3.81
```

### 環境変数
構築を始める前に、カレントディレクトリにある```.env.template``` と ```secrets.env.template``` を参考にして、```.env``` と ```.secrets.env```を作成してください。

もし、Kinesis Data Streams に Change Streams をエクスポートしたい場合は、```.secrets.env```に以下のような記載をしましょう。

```
EXPORT_DESTINATION=kinesisStream
```

## 手順
**1. ノードインスタンス(EC2)にキーペアを作成します**

EC2インスタンスにSSH接続するために用いられます。

<br>

**2. EKS クラスター用の KMS キーを作成します**

Kubernetes Secrets を複合化するために用いられます。

```
$ make kms
```

<br>

**3. EKS クラスターとノードグループを作成します**

```
$ make build
```

・クラスター作成

```.env```にある環境変数を参照して、```cluster.yaml```を作成します。その後に、```eksctl create cluster```コマンドが実行されます。

・ノードグループ作成

```.env```にある環境変数を参照して、```nodegroup.yaml```を作成します。その後に、```eksctl create nodegroup```コマンドが実行されます。

<br>

**4. ノードグループロールに Kinesis ポリシーを付与する**

MxTransporter コンテナが Change Streams を Kinesis Data Streams へエクスポートするために、ノードグループロールに Kinesis ポリシーを付与する必要があります。ノードグループロールは```eksctl create cluster```コマンドが実行されたときに作成され、例えば```AmazonKinesisFullAccess```のようなポリシーを付与します。

<br>

**5. Kubernetes Secrets 作成**

```secrets.env```にある環境変数を収集し、それらを Kubernetes Secrets としてクラスターに作成します。

```
$ make secrets
```

**6. Kubernetes リソースをデプロイする**
もしオプションの環境変数を```.secrets.env```セットしたのであれば、```./templates/stateless.yaml```内の env パラメータを編集します。Kubernetes で実行されているコンテナに必要な環境変数のみを設定します。

Kubernetes リソースは helm によって作成されます。

以下のコマンドは StatefulSet、HeadlessService、Horizontal Pod Autoscaler、PVC を作成します。

```
$ make deploy
```

このコマンドでは、以下のような処理が行われます。<br>
・Docker イメージビルド<br>
・ECR リポジトリへのログイン<br>
・ECR リポジトリへ Docker イメージを送信<br>
・helm の variables を作成します<br>
・helm テンプレートをデプロイする<br>

<br>

**7. Kubernetes リソースをアップグレードする**

以下のコマンドで、Kubernetes リソースをアップグレードできます。

環境変数の```ECR_REPO_TAG```を更新しないと、最新の Docker イメージが参照されず、新しいコンテナが作成されない点に注意してください。

```
$ make upgrade
```

<br>

# アーキテクチャ

![image](https://user-images.githubusercontent.com/37132477/141406354-2616bdf9-8f19-4d3f-b752-23ecaeae2611.png)

MongoDB コレクションごとに Pod が作成され、それぞれの Pod に 永続ボリュームが紐付いています。StatefulSet が作成されるため、ポッドが停止した場合でも、永続ボリュームに保存されている再開トークンを再度参照することで、変更ストリームを取得できます。