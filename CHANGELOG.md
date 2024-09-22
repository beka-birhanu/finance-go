# [1.2.0](https://github.com/beka-birhanu/finance-go/compare/v1.1.0...v1.2.0) (2024-09-22)


### Bug Fixes

* add input validation, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([cb995e9](https://github.com/beka-birhanu/finance-go/commit/cb995e9492ef13988804b7075c4f46ef2ec0cf57))
* change sortfield option createdAt -> date ([27dd756](https://github.com/beka-birhanu/finance-go/commit/27dd756fa19f0ccdaf188b47c293379a803b97eb))
* collect error maping from ierr.IErr to errapi struct, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([5bbbebf](https://github.com/beka-birhanu/finance-go/commit/5bbbebf1a92a75bc10b63cc3524203ba8782d7ba))
* custom marshaling float32 avoid casting to float64, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([f665f55](https://github.com/beka-birhanu/finance-go/commit/f665f55f70463c5be96eb86747a0f36fa6b368eb))
* move amount rounding to domain layer ([6e86335](https://github.com/beka-birhanu/finance-go/commit/6e8633543bc7f638a529c40a0f5c7fd04fdac094))
* **repo:** fix query builder logic for pagination ([e100bb8](https://github.com/beka-birhanu/finance-go/commit/e100bb87f24edbf62ae340ff4b7d25d5cd8de1d3))
* update confirmUserID to return better error, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([1dd523c](https://github.com/beka-birhanu/finance-go/commit/1dd523c7a7599cdeafe4c4e815bfd65a35d836e0))


### Features

* add error message filtering in utils for graph api's, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([288cced](https://github.com/beka-birhanu/finance-go/commit/288cced1d54d66642c82d79d07cd85ee570ca599))
* add utils to hlep constructing error for graph, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([f85e5f7](https://github.com/beka-birhanu/finance-go/commit/f85e5f7362e63bc4738ab8b13db5d34fee588d49))
* **api:** define expense model ([e6b280e](https://github.com/beka-birhanu/finance-go/commit/e6b280ee69f1435313d12279d6de875d4f0a5016))
* **api:** implement get expense query, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([de2ec68](https://github.com/beka-birhanu/finance-go/commit/de2ec68e7c147da2b1b87ba5a7b8128d47ed9023))
* implement create expense graph endpoint, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([ca9632a](https://github.com/beka-birhanu/finance-go/commit/ca9632a028f78f43032369da7edf995519544bf4))
* implement get multiple expenses, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([e16f3de](https://github.com/beka-birhanu/finance-go/commit/e16f3ded99f658d38d0d17642f1b3229ed0beb32))
* implement patch expense graph api, [#24](https://github.com/beka-birhanu/finance-go/issues/24) ([df0622a](https://github.com/beka-birhanu/finance-go/commit/df0622a57e308233996d0defcfc94110ded8c102))
* udpate router to register graphql endpoint ([5d53872](https://github.com/beka-birhanu/finance-go/commit/5d53872c844c06d4d8e0379abd3e0e2c431e5b1e))

# [1.1.0](https://github.com/beka-birhanu/finance-go/compare/v1.0.0...v1.1.0) (2024-09-19)


### Bug Fixes

* fix config for semantic-release ([696f742](https://github.com/beka-birhanu/finance-go/commit/696f7426114bcf872e28eb8029bf1efce97a90a9))


### Features

* **api:** consume rate limiter servince in router ([8e15967](https://github.com/beka-birhanu/finance-go/commit/8e159677527487cccc414c7cfd83dc1f910df374))
* **api:** define IRateLimiter interface, [#23](https://github.com/beka-birhanu/finance-go/issues/23) ([78fa682](https://github.com/beka-birhanu/finance-go/commit/78fa68279042a195a6c868a0fb41750930cc9770))
* **api:** implement reate limiter per IP ([4446150](https://github.com/beka-birhanu/finance-go/commit/44461503e75c71ade8bad196b71b5de9545ddc0d))
* implement middleware that consums IRateLimter interface ([85b8712](https://github.com/beka-birhanu/finance-go/commit/85b8712179f6b4d18a29bf9e0534c43b6c40e1ee))

## [1.0.1](https://github.com/beka-birhanu/finance-go/compare/v1.0.0...v1.0.1) (2024-09-19)


### Bug Fixes

* fix config for semantic-release ([696f742](https://github.com/beka-birhanu/finance-go/commit/696f7426114bcf872e28eb8029bf1efce97a90a9))
