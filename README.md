# AI Code Review SaaS MVP

Go製のAIコードレビューSaaS MVPのバックエンドAPIです。

現時点ではOpenAI API、DB、GitHub連携、フロントエンドは実装せず、ダミーのレビュー結果を返します。

設計資料は `docs/` 配下にあります。API仕様は `docs/api-spec.md` を参照してください。

## 起動方法

```sh
go run ./cmd/server
```

サーバーは `http://localhost:8080` で起動します。

## API

### ヘルスチェック

```sh
curl http://localhost:8080/health
```

レスポンス例:

```json
{
  "status": "ok"
}
```

### コードレビュー

```sh
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{"language":"go","code":"package main\n\nfunc main(){println(\"hello\")}" }'
```

レスポンス例:

```json
{
  "score": 85,
  "summary": "全体的にシンプルで読みやすいコードです。",
  "issues": [
    {
      "severity": "low",
      "title": "println の使用",
      "description": "本格的なアプリケーションでは fmt.Println やログライブラリの使用を検討してください。"
    }
  ],
  "codex_prompt": "以下のレビュー内容をもとに、既存の挙動を変えずにコードを改善してください。"
}
```

## バリデーション

- `language` が空の場合はエラー
- `code` が空の場合はエラー
- `code` が100,000バイトを超える場合はエラー
- 不正なJSONの場合はエラー

エラー時もJSONでレスポンスします。

```json
{
  "error": {
    "code": "validation_error",
    "message": "language は必須です。"
  }
}
```

## ディレクトリ構成

```txt
cmd/server/          APIサーバーのエントリーポイント
internal/handler/    HTTPハンドラーとレスポンス処理
internal/model/      リクエスト/レスポンスモデル
internal/service/    レビュー処理のインターフェースと実装
```
