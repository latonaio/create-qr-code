# create-qr-code
## 概要
create-qr-codeは、受け取ったJSON文字列をもとにQRコードの生成を行うマイクロサービスです。

## 動作環境
create-qr-codeは、aion-coreのプラットフォーム上での動作を前提としています。 使用する際は、事前に下記の通りAIONの動作環境を用意してください。 
- ARM CPU搭載のデバイス(NVIDIA Jetson シリーズ等) 
- OS: Linux Ubuntu OS 
- CPU: ARM64 
- Kubernetes 
- [AION のリソース](https://github.com/latonaio/aion-core)

## セットアップ
```
git clone git@github.com:latonaio/create-qr-code.git
cd path/to/create-qr-code.git
make docker-build
```

## 起動方法
```
kubectl apply -f create-qr-code.yaml
```

## kanban との通信
### kanbanから受信するデータ
kanban から受信する metadata に下記の情報を含む必要があります。

| key | value |
| --- | --- |
| size | QRコードのサイズ |
| json_str | QRコードに埋め込む文字列 |
| output_path | 出力ファイルパス |

具体例 : 
```example
# metadata (map[string]interface{}) の中身

"size": "200"
"json_str": "{id: xxxxx, name: yyyyy}"
"output_path": "/aaa/bbb/ccc/sample.png"
```

### kanban に送信するデータ
kanban に送信する metadata は下記の情報を含みます。

| key | type | description |
| --- | --- | --- |
| file_path | string | ファイルパス |

具体例: 
```example
# metadata (map[string]interface{}) の中身

"file_path": "/aaa/bbb/ccc/sample.png"
```
