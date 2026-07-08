# AI Code Review SaaS MVP

Go製のAIコードレビューSaaS MVPのバックエンドAPIです。

現時点ではOpenAI API、DB、GitHub連携、フロントエンドは実装せず、ダミーのレビュー結果を返します。

設計資料は `docs/` 配下にあります。API仕様は `docs/api-spec.md` を参照してください。

## 起動方法

```sh
go run ./cmd/server
```

サーバーはデフォルトで `http://127.0.0.1:8080` で起動します。

待ち受けアドレスを変更したい場合は `SERVER_ADDR` を指定します。

```sh
SERVER_ADDR=127.0.0.1:8081 go run ./cmd/server
```

## API

### ヘルスチェック

```sh
curl http://127.0.0.1:8080/health
```

レスポンス例:

```json
{
  "status": "ok"
}
```

### コードレビュー

```sh
curl -X POST http://127.0.0.1:8080/review \
  -H "Content-Type: application/json" \
  -d '{"language":"go","code":"package main\n\nfunc main(){println(\"hello\")}" }'
```

レスポンス例:

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
