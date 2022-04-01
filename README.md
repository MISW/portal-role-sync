# portal-role-sync
- portalを参照して、メンバーとそのロール情報を取得し、auth0に渡す。

# 仕様
1. [portal](https://github.com/MISW/Portal)のapiを叩いてユーザー(とそのロール)一覧の情報を取得する。
2. auth0のmanagement-apiを叩いてrules configに```key: "members", value: ${JSON ユーザーとそのロール情報}```をセットする。 
(Auth0ではこの```members```を認証に用いている。例えばロールが```member```と```admin```のユーザーのみ許可するなど。)

# 環境変数
- [.env.template](./.env.template)
    - auth0の変数: auth0ログイン後、```Applications > APIs > Auth0 Management API > Machine to Machine Applications > MISW Portal``` 
    - portal_token: [みすportal](https://github.com/MISW/Portal)に設定した```EXTERNAL_INTEGRATION_TOKENS``` 
- ```.env```ファイルを作って[.envrc](./.envrc)を利用すると便利。
