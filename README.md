# portal-role-sync

- portal を参照して、メンバーとそのロール情報を取得し、auth0 に渡す。

# 仕様
1. [portal](https://github.com/MISW/Portal)のapiを叩いてユーザー(とそのロール)一覧の情報を取得する。
2. auth0のmanagement-apiを叩いてuserのapp_metadataに
    ```app_metadata.json
    {
        "app_metadata": {
            "portal_role": ${ロール}
        }
    }
    ```
    をセットする。 
    (Auth0ではこの```portal_role```を認証に用いている。例えばロールが```member```と```admin```のユーザーのみ許可するなど。)
    なお、初めから全ユーザのmeta_data更新をした場合apiを叩くlimitを超えるため、一旦getしてから、updateの必要があるものだけupdateする。


# 環境変数

- [.env.template](./.env.template)
    - auth0の変数: auth0ログイン後、```Applications > APIs > Auth0 Management API > Machine to Machine Applications > MISW Portal``` 
    - portal_token: [みすportal](https://github.com/MISW/Portal)に設定した```EXTERNAL_INTEGRATION_TOKENS``` 
- ```.env```ファイルを作って[.envrc](./.envrc)を利用すると便利。

# 参考
- [Auth0: Manage Metadata Using the Management API](https://auth0.com/docs/manage-users/user-accounts/metadata/manage-metadata-api) 
- [limits](https://auth0.com/docs/troubleshoot/customer-support/operational-policies/rate-limit-policy/management-api-endpoint-rate-limits#self-service-subscription-limits) 
