# API Specification

## 1. 目的

このドキュメントは、レビュー画面仕様書 `docs/review-screen-spec.md` に合わせて、AI Code Review SaaS MVP のAPI仕様を明確にするためのものです。

現時点では、Go製バックエンドAPIがダミーのレビュー結果を返します。OpenAI API連携、レビュー履歴保存、GitHub連携、認証機能は今後の拡張対象です。

## 2. Base URL

ローカル開発環境では以下を使用します。

```txt
http://localhost:8080
```

## 3. Endpoints

| Method | Path | Description |
| --- | --- | --- |
| GET | `/health` | APIサーバーのヘルスチェック |
| POST | `/review` | コードレビューを実行し、レビュー結果を返す |

## 4. Health Check API

### Request

```http
GET /health
```

### Response

```json
{
  "status": "ok"
}
```

## 5. Review API

### Request

```http
POST /review
Content-Type: application/json
```

### Review Request

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `language` | string | Yes | レビュー対象コードの言語。例: `go`, `javascript`, `python` |
| `code` | string | Yes | レビュー対象のコード本文 |

### Request Example

```json
{
  "language": "go",
  "code": "package main\n\nfunc main(){println(\"hello\")}"
}
```

### Review Response

現行MVPの `/review` API は以下の形式を返します。

| Field | Type | Description |
| --- | --- | --- |
| `score` | number | 0〜100の総合スコア |
| `summary` | string | レビュー全体の短い要約 |
| `issues` | array | 問題点一覧 |
| `codex_prompt` | string | Codexへ貼り付けるための修正依頼プロンプト |

### Review Issue

| Field | Type | Description |
| --- | --- | --- |
| `severity` | string | 重要度。現時点では `low` を返す。将来的に `high`, `medium`, `low` を使用する |
| `title` | string | 指摘事項のタイトル |
| `description` | string | 指摘事項の説明 |

### Response Example

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

## 6. Error Response

エラー時もJSONでレスポンスします。

### Error Response Format

| Field | Type | Description |
| --- | --- | --- |
| `error` | object | エラー情報 |
| `error.code` | string | アプリケーション内で扱うエラーコード |
| `error.message` | string | ユーザーまたは開発者向けのエラーメッセージ |

### Error Response Example

```json
{
  "error": {
    "code": "validation_error",
    "message": "language は必須です。"
  }
}
```

## 7. HTTP Status

| Status | Code | Condition |
| --- | --- | --- |
| `200 OK` | - | リクエストが正常に処理された |
| `400 Bad Request` | `invalid_json` | リクエストボディが不正なJSON、またはJSONオブジェクトが複数送信された |
| `400 Bad Request` | `validation_error` | `language` または `code` が空、または `code` が長すぎる |
| `404 Not Found` | `not_found` | 存在しないエンドポイントへアクセスした |
| `405 Method Not Allowed` | `method_not_allowed` | 許可されていないHTTPメソッドでアクセスした |
| `500 Internal Server Error` | `review_failed` | レビュー処理中にサーバー内部でエラーが発生した |

## 8. Validation Rules

### `language`

- 必須
- 空文字不可
- 空白のみ不可

### `code`

- 必須
- 空文字不可
- 空白のみ不可
- 100,000バイト以下

### JSON

- リクエストボディは有効なJSONである必要がある
- 1リクエストにつきJSONオブジェクトは1つだけ送信する
- 未定義フィールドは許可しない

## 9. Curl Examples

### Health Check

```sh
curl http://localhost:8080/health
```

### Review

```sh
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{"language":"go","code":"package main\n\nfunc main(){println(\"hello\")}" }'
```

### Validation Error

```sh
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{"language":"","code":"package main"}'
```

Response:

```json
{
  "error": {
    "code": "validation_error",
    "message": "language は必須です。"
  }
}
```

### Invalid JSON

```sh
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{"language":"go","code":'
```

Response:

```json
{
  "error": {
    "code": "invalid_json",
    "message": "リクエストボディが不正なJSONです。"
  }
}
```

## 10. Review Screenとの対応

`docs/review-screen-spec.md` では、レビュー画面に以下の表示項目を持たせる想定です。

- 総合スコア
- 要約
- 良かった点
- 問題点一覧
- 改善案
- Codex用プロンプト

現行MVPのAPIは、総合スコア、要約、問題点一覧、Codex用プロンプトを返します。

良かった点、改善案、問題点の該当箇所、改善理由は、今後のAPI拡張で追加する候補です。

## 11. Future Extensions

将来的には以下のフィールド追加を検討します。

```json
{
  "positives": [
    "関数の役割が分かりやすい"
  ],
  "suggestions": [
    "エラー処理を追加する"
  ],
  "issues": [
    {
      "severity": "medium",
      "title": "エラー処理が不足している",
      "description": "失敗時の処理が不足しています。",
      "location": "main.go:10",
      "reason": "エラーを無視すると原因調査が難しくなるためです。"
    }
  ]
}
```

これらは画面仕様に合わせた拡張候補であり、現時点のAPIレスポンスには含めません。
