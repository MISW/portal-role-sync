# portal-role-sync

- portal を参照して、メンバーとそのロール情報を取得し、auth0 に渡す。

# 仕様

1. [portal](https://github.com/MISW/Portal)の api を叩いてユーザー(とそのロール)一覧の情報を取得する。
2. auth0 の management-api を叩いて user の app_metadata に

   ```app_metadata.json
   {
       "app_metadata": {
           "portal_role": ${ロール}
       }
   }
   ```

   をセットする。
   (Auth0 ではこの`portal_role`を認証に用いている。例えばロールが`member`と`admin`のユーザーのみ許可するなど。)
   なお、初めから全ユーザの meta_data 更新をした場合 api を叩く limit を超えるため、一旦 get してから、update の必要があるものだけ update する。

# 環境変数

- [.env.template](./.env.template)
  - auth0 の変数: auth0 ログイン後、`Applications > APIs > Auth0 Management API > Machine to Machine Applications > MISW Portal`
  - portal_token: [みす portal](https://github.com/MISW/Portal)に設定した`EXTERNAL_INTEGRATION_TOKENS`
- `.env`ファイルを作って[.envrc](./.envrc)を利用すると便利。

# 参考

- [Auth0: Manage Metadata Using the Management API](https://auth0.com/docs/manage-users/user-accounts/metadata/manage-metadata-api)
- [limits](https://auth0.com/docs/troubleshoot/customer-support/operational-policies/rate-limit-policy/management-api-endpoint-rate-limits#self-service-subscription-limits)
