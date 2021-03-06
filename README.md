# cloudfunctions-example

## これは

Cloud Functions の各種実装例です. 2020/01/13 時点で以下の実装を含みます.

* HTTP Trigger -> Google Cloud Storage Bucket
* Google Cloud Storage Bucket -> BigQuery

## Files

```sh
$ tree .
.
├── README.md
├── gcs-to-bigquery
│   ├── go.mod
│   └── handler.go
├── httptrigger-to-gcs
│   ├── go.mod
│   ├── go.sum
│   ├── handler.go
│   └── handler_test.go
├── output
│   ├── gcs-to-bigquery.zip
│   └── httptrigger-to-gcs.zip
└── terraform
    ├── Makefile
    ├── bigquery
    │   └── schema.json
    ├── main.tf
    ├── resources.tf
    └── variables.sample.tf
```

## Cloud Functions

### httptrigger-to-gcs

HTTP Trigger な関数です. デプロイすると URL が払い出されます. 以下のような JSON ボディを POST するとポストした内容が Google Cloud Storage (以後, GCS) に保存されます.

```json
[
  {
    "name": "hoge",
    "event": "test",
    "timestamp": "1234567890"
  }
]
```

以下, 実行例です.

```sh
$ curl -X POST https://asia-northeast1-your-sample-pj.cloudfunctions.net/httptrigger-to-gcs -d '[{"name":"hoge", "event":"test", "timestamp":"1234567890"}]'
ok
```

尚, 意図しない JSON や JSON 以外を POST するとエラーが返ります.

```sh
$ curl -X POST https://asia-northeast1-your-sample-pj.cloudfunctions.net/httptrigger-to-gcs -d 'foo' -v
... 略 ...
< HTTP/2 400
< content-type: text/plain; charset=utf-8
... 略 ...
error
```

### gcs-to-bigquery

GCS にオブジェクトが保存されたことをトリガーとする関数です. 上記の HTTP Trigger な関数で POST されたデータを GCS を介して BigQuery に保存する流れです. 下図のように BigQuery に POST されたデータがストアされます.

![画像](https://raw.githubusercontent.com/inokappa/cloudfunctions-example/master/docs/images/2020011401.png)

## インフラ

### Google Cloud Platform の準備

インフラは Terraform を利用してまるっと作成します. ただし, Terraform は v0.12 には対応していません. すいません.

1. プロジェクトの作成にする
1. Cloud Functions API を有効にする
1. Terraform 用のサービスアカウントを作成し, JSON フォーマットの鍵を作成してダウンロードする
1. サービスアカウントの鍵を credentials ディレクトリに保存する
1. 環境変数 `GOOGLE_APPLICATION_CREDENTIALS` にサービスアカウントの鍵を設定する

```sh
export GOOGLE_APPLICATION_CREDENTIALS=path/to/credentials/xxxxxxxxxxxxxx.json
```

### tfstate 用のバケットを作成

tfstate ファイルは GCS のバケットに保存する為, 手動でバケットを作成します.

### variables.tf の修正

環境に合わせて, variables.sample.tf の内容を修正して variables.tf というファイル名で保存します. variables.tf にはプロジェクト ID や BigQuery のテーブルセット等を設定する必要があります.

### あとはいつもの Terraform

```sh
$ cd terraform
$ terraform init -backend-config="bucket=your-tfstate-bucket-name"
$ make plan
$ make apply
```

### BigQuery のスキーマ

BigQuery のスキーマは以下の通りです.

```sh
$ cat terraform/bigquery/schema.json
[
    {
        "name": "name",
        "type": "STRING",
        "mode": "NULLABLE",
        "description": "Sample Name"
    },
    {
        "name": "event",
        "type": "STRING",
        "mode": "NULLABLE",
        "description": "Sample Event"
    },
    {
        "name": "timestamp",
        "type": "INT64",
        "mode": "NULLABLE",
        "description": "Sample timestamp"
    }
]
```

## Github Actions

Cloud Functions コードの CI は Github Actions を見様見真似で追加しています. 以下のジョブが実行されます.

* Lint (go vet)
* Test (go test)
* Build (go build ※ ただビルドするだけです)

## todo

* インフラ構成を Terraform v0.12 に対応させる
* gcloud コマンドを利用する
