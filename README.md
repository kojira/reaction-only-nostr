# Reaction Only

Nostrのリアクションのみで会話する実験的な静的Webクライアント。GoをWebAssemblyにコンパイルし、ブラウザ上で動作します。

NIP-07対応ブラウザ拡張でログインし、対象イベントIDと作者の公開鍵を小文字64桁hexで指定して送信します。フィードのボタンからリアクション自体へ返信できます。NIP-25 kind 7を使用し、NIP-30カスタム絵文字は `:shortcode:` と `emoji` タグで表現します。秘密鍵は拡張機能が管理します。

公開リレーに送信され、誰でも読めます。フィードは署名検証を行わないため、悪意あるリレーのデータを信用しないでください。画像URLは第三者サーバーに接続します。公開リレーの受理・保持は保証されません。リレー認証、npub/nevent、削除、スレッド整理は未対応です。

## ローカルビルド

Go 1.25以降を使用し、`GOOS=js GOARCH=wasm go build -o app.wasm .` と `cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .` を実行し、HTTPサーバーで配信してください。GitHub Actionsは同じ処理を実行してPagesへデプロイします。
