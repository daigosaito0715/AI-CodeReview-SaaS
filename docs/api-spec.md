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
| `language` | string | レビュー対象コードの言語 |
| `review_version` | string | レビュー処理のバージョン |
| `score` | number | 0〜100の総合スコア |
| `summary` | string | レビュー全体の短い要約 |
| `strengths` | array | 良かった点の一覧 |
| `issues` | array | 問題点一覧 |
| `suggestions` | array | 次に行う改善案の一覧 |
| `codex_prompt` | string | Codexへ貼り付けるための修正依頼プロンプト |

### Review Issue

| Field | Type | Description |
| --- | --- | --- |
| `severity` | string | 重要度。現時点では `low` を返す。将来的に `high`, `medium`, `low` を使用する |
| `title` | string | 指摘事項のタイトル |
| `description` | string | 指摘事項の説明 |
| `location` | string | 該当箇所 |
| `reason` | string | 改善理由 |

### Response Example

```json
{
  "language": "go",
  "review_version": "dummy-v1",
  "score": 85,
  "summary": "全体的にシンプルで読みやすいコードです。現時点では大きな問題はありませんが、出力方法と将来の拡張性に改善余地があります。",
  "strengths": [
    "コード量が少なく、処理の流れを追いやすい",
    "関数の役割が明確で、最小構成として理解しやすい",
    "不要な複雑さが少なく、改善方針を立てやすい"
  ],
  "issues": [
    {
      "severity": "low",
      "title": "println の使用",
      "description": "本格的なアプリケーションでは fmt.Println やログライブラリの使用を検討してください。",
      "location": "コード全体",
      "reason": "出力方法を標準的なAPIに寄せることで、後からログ出力やテストに置き換えやすくなります。"
    }
  ],
  "suggestions": [
    "fmt.Println または用途に合ったログ出力へ置き換える",
    "今後処理が増える場合に備えて、入出力処理とビジネスロジックを分ける"
  ],
  "codex_prompt": "以下のレビュー結果をもとに、go のコードを改善してください。\n\n目的:\n- 既存の挙動を変えずに、読みやすさ・保守性・実運用時の扱いやすさを改善する\n- レビュー指摘を反映し、必要以上に大きな設計変更は避ける\n\nレビュー指摘:\n- [low] println の使用: 本格的なアプリケーションでは fmt.Println やログライブラリの使用を検討してください。\n  該当箇所: コード全体\n  改善理由: 出力方法を標準的なAPIに寄せることで、後からログ出力やテストに置き換えやすくなります。\n\n改善案:\n- fmt.Println または用途に合ったログ出力へ置き換える\n- 今後処理が増える場合に備えて、入出力処理とビジネスロジックを分ける\n\n作業条件:\n- 既存の動作を維持する\n- 変更理由を簡潔に説明する\n- 可能であれば小さな差分で修正する\n- 不要な依存関係を追加しない\n"
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

現行MVPのAPIは、総合スコア、要約、良かった点、問題点一覧、改善案、Codex用プロンプトを返します。

`strengths` は画面仕様の「良かった点」に対応し、`suggestions` は「改善案」に対応します。`issues.location` は「該当箇所」、`issues.reason` は「改善理由」に対応します。

## 11. Future Extensions

将来的には以下のフィールド追加を検討します。

```json
{
  "review_id": "rev_123",
  "created_at": "2026-07-08T00:00:00Z",
  "model": "gpt-5",
  "repository": "owner/repo",
  "pull_request": 123
}
```

これらは履歴保存、OpenAI API連携、GitHub連携に合わせた拡張候補であり、現時点のAPIレスポンスには含めません。
