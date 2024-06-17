# 発火村 (はっかむら)

![](./frontend/public/og.png)

## メンバー

- ikura-hamu (リーダー / バックエンド)
- cp20 (フロントエンド)
- tidus (バックエンド)
- zoi_dayo (フロントエンド)
- Alietty (フロントエンド)

## 開発者クイックスタート

### バックエンド

#### 動かすために最低限必要

- Docker

#### 開発に必要

- Go 1.22 以上
- [Task](https://github.com/go-task/task)

コマンドの使い方は`task --list`から確認。

#### 説明

初めて使うとき、backendフォルダに`.env`という空のファイルを置く必要がある。

8080ポートでサーバーアプリケーション、8081ポートでadminer、3306ポートでMySQLが立ち上がる。アプリは[Air](https://github.com/air-verse/air)を使ってホットリロードを設定しているので、Goのコードを変更しても`go run main.go`とかする必要は無く、ちょっと待つとすぐ変更が反映される。他をいじったら一旦`task down`したり `task clean` したり必要。

メッセージの変換は、手元ではGPTを使わず、` (converted by mock)`と後ろにつけるようにしている。

### フロントエンド

Node.js v20 を想定しています

```sh
cd frontend
npm i
```

開発用サーバーは次のコマンドで起動できます

```sh
npm run dev
```

次のコマンドでビルドできます

```sh
npm run build
```

ビルド後 `frontend/dist` ファイルを適切に配信してください

### VSCode用セットアップ

ワークスペースで推奨されている拡張機能を入れてください

`.vscode/settings-template.json` をコピーして `.vscode/settings.json` を作ってください

