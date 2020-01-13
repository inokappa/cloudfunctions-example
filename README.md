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
    └── variables.tf
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

### gcs-to-bigquery

GCS にオブジェクトが保存されたことをトリガーとする関数です. 上記の HTTP Trigger な関数で POST されたデータを GCS を介して BigQuery に保存する流れです.

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

環境に合わせて, variables.tf の内容を修正します. プロジェクト ID や BigQuery のテーブルセット等を設定する必要があります.

### あとはいつもの Terraform

```sh
$ cd terraform
$ make init
$ make plan
$ make apply
```

## todo

* インフラ構成を Terraform v0.12 に対応させる
