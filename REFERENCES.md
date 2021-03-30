# References
My references to write these codes. For anyone, who want to know more about my codes.


## OpenAPI Generation for go
- [OpenAPI3を使ってみよう！Go言語でクライアントとスタブの自動生成まで！](https://techblog.zozo.com/entry/openapi3/go)
- [Go の Open API 3.0 のジェネレータ oapi-codegen を試してみた](https://future-architect.github.io/articles/20200701/)
- [OpenAPI3 + OpenAPI generator でgolangサーバ・TypeScriptクライアントの実装を試す](https://qiita.com/doriven/items/7422f565d6ad2e8ff956)
### Note
I chose go net/http generator using [openapi-generator-cli](https://github.com/OpenAPITools/openapi-generator).
Since I'm new to golang, and didn't want do do myself many.
### Command
```openapi-generator-cli generate -i accounts.v1.yaml -g go-server --additional-properties=packageName=gen```

## REST API best practice
- [RESTful API設計のベストプラクティス　-Web API the Good Partsまとめ-](https://www.mushroom-blog.com/420/)
- [翻訳: WebAPI 設計のベストプラクティス](https://qiita.com/mserizawa/items/b833e407d89abd21ee72)
### Note
[My old api](https://github.com/nuxt-image-board/backend) was not unified, and took a long time to unify endpoints, and models.
If some models are still not unified, sorry about that.

## Golang net/http uses
- [Go using mux Router - How to pass my DB to my handlers](https://stackoverflow.com/questions/33646948/go-using-mux-router-how-to-pass-my-db-to-my-handlers)
- [HTTP Middleware の作り方と使い方](https://tutuz-tech.hatenablog.com/entry/2020/03/23/220326)

## Golang Config uses
- [Goでconfigファイルを読み込む](https://qiita.com/wooootack/items/c38f3bbd916843df1256)

## Golang go-jwt uses
- [【セキュリティ】jwt-goを使ってみる - その１](https://blog.motikan2010.com/entry/2017/05/12/jwt-go%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%A6%E3%81%BF%E3%82%8B)

## Golang mgm(odm) uses
- [mgm](https://github.com/Kamva/mgm)

## Golang elasticsearch uses
- [elasticが開発した公式のGo言語ElasticSearchクライアントについてまとめてみる](https://qiita.com/shiei_kawa/items/d992f7fdd4c75906ea0b)

## Golang net/httptest/httptest uses
- [net/http を httptest を使ってテストする方法](https://hawksnowlog.blogspot.com/2019/04/golang-net-http-test.html)
- [mopeneko/hacku2020-web/router(example)](https://github.com/mopeneko/hacku2020-web/tree/master/api/router)
- [GoのAPIのテストにおける共通処理](https://medium.com/@timakin/go-api-testing-173b97fb23ec)

## Golang specification
- [Goのパッケージのinit( )はいつどういう順番で呼ばれるか](https://qiita.com/YusukeIwaki/items/f1f92c23d7ee0ca8dc7a)
- [インタフェースの実装パターン #golang](https://qiita.com/tenntenn/items/eac962a49c56b2b15ee8)

## Sync elasticsearch and mongodb
- [MongoDB⇒ElasticSearchへの自動連携](https://qiita.com/chenglin/items/92a3ea29be7e72c66bb1)
- [Monstache](https://rwynn.github.io/monstache-site/)

## Editor problem
- [VS Codeのgo testに-vオプションを付ける方法](https://qiita.com/mako2kano/items/3923b9afac619bb781f7)