# Changelog

## [0.0.10](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.9...v0.0.10) (2023-12-07)


### Bug Fixes

* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.27 ([#95](https://github.com/spectrocloud-labs/validator-plugin-network/issues/95)) ([d3b7eca](https://github.com/spectrocloud-labs/validator-plugin-network/commit/d3b7ecaf072de035fec6eaed0a7a3d82016cf58a))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.28 ([#99](https://github.com/spectrocloud-labs/validator-plugin-network/issues/99)) ([c7a4560](https://github.com/spectrocloud-labs/validator-plugin-network/commit/c7a45606580021c2a8e7a3d4bc8b22606c9e039c))
* ensure unique rules to avoid overwriting conditions ([#102](https://github.com/spectrocloud-labs/validator-plugin-network/issues/102)) ([6cad704](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6cad70425f7165173350904ba5b67bbabd442eff))


### Other

* **deps:** update actions/setup-go action to v5 ([#101](https://github.com/spectrocloud-labs/validator-plugin-network/issues/101)) ([9e46eb3](https://github.com/spectrocloud-labs/validator-plugin-network/commit/9e46eb3deeb0377d51fdc88af54b3d572be0413c))
* **deps:** update actions/setup-python action to v5 ([#100](https://github.com/spectrocloud-labs/validator-plugin-network/issues/100)) ([d982bc8](https://github.com/spectrocloud-labs/validator-plugin-network/commit/d982bc891484a55b4640e094f26159ef1956a666))
* **deps:** update anchore/sbom-action action to v0.15.1 ([#98](https://github.com/spectrocloud-labs/validator-plugin-network/issues/98)) ([408dee2](https://github.com/spectrocloud-labs/validator-plugin-network/commit/408dee268d3be654998ea7341baa9141bba890fb))
* **deps:** update google-github-actions/release-please-action action to v4 ([#97](https://github.com/spectrocloud-labs/validator-plugin-network/issues/97)) ([6b4bcd2](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6b4bcd298830015e0a20f0b6c335c382973a0a03))

## [0.0.9](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.8...v0.0.9) (2023-11-29)


### Bug Fixes

* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.2 ([#93](https://github.com/spectrocloud-labs/validator-plugin-network/issues/93)) ([5fd248a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/5fd248ac4551a0bbb6cd2eaf0e3cda26da6e4c05))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.26 ([#92](https://github.com/spectrocloud-labs/validator-plugin-network/issues/92)) ([9d970b7](https://github.com/spectrocloud-labs/validator-plugin-network/commit/9d970b790a3906abd6051909d732da6713a1311f))
* refactor proxy config; use image w/ update-ca-certificates preinstalled ([#91](https://github.com/spectrocloud-labs/validator-plugin-network/issues/91)) ([8f0e61c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/8f0e61cf9027423947c8b36d8eabea09ce6b90b4))

## [0.0.8](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.7...v0.0.8) (2023-11-28)


### Features

* add support for proxy setup and enabling proxy certs ([#90](https://github.com/spectrocloud-labs/validator-plugin-network/issues/90)) ([02ac489](https://github.com/spectrocloud-labs/validator-plugin-network/commit/02ac4892b9700b9facfa45efb5bc30b6a576144e))


### Bug Fixes

* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.25 ([#88](https://github.com/spectrocloud-labs/validator-plugin-network/issues/88)) ([91f92c8](https://github.com/spectrocloud-labs/validator-plugin-network/commit/91f92c8eec2f2b9b5a3515778d401b0ca0f21678))


### Other

* **deps:** update anchore/sbom-action action to v0.15.0 ([#89](https://github.com/spectrocloud-labs/validator-plugin-network/issues/89)) ([6f5aa92](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6f5aa923d6189d7281b1282f21484239894b046d))
* **deps:** update docker/build-push-action digest to 4a13e50 ([#87](https://github.com/spectrocloud-labs/validator-plugin-network/issues/87)) ([8453b0a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/8453b0ac24d2e4c5eb65b56f8a6f5d26f978de36))


### Refactoring

* standardize get CR in Reconcile ([638bc76](https://github.com/spectrocloud-labs/validator-plugin-network/commit/638bc7667bd7b3bf0d241243f30efceb42abd83e))

## [0.0.7](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.6...v0.0.7) (2023-11-16)


### Bug Fixes

* **deps:** update kubernetes packages to v0.28.4 ([#82](https://github.com/spectrocloud-labs/validator-plugin-network/issues/82)) ([19525b6](https://github.com/spectrocloud-labs/validator-plugin-network/commit/19525b603ef53f7bd3681c556ef64ba929fbec17))
* set owner references on VR to ensure cleanup ([#84](https://github.com/spectrocloud-labs/validator-plugin-network/issues/84)) ([f9a4357](https://github.com/spectrocloud-labs/validator-plugin-network/commit/f9a43579ed05b576d2959d52592404b8e0827184))

## [0.0.6](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.5...v0.0.6) (2023-11-13)


### Bug Fixes

* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.1 ([#78](https://github.com/spectrocloud-labs/validator-plugin-network/issues/78)) ([97f0a94](https://github.com/spectrocloud-labs/validator-plugin-network/commit/97f0a94055390b2707d1e05cb883e5885e5c730d))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.17 ([#76](https://github.com/spectrocloud-labs/validator-plugin-network/issues/76)) ([88b278d](https://github.com/spectrocloud-labs/validator-plugin-network/commit/88b278d5c6e13a55ac59191c31df492687dd7c4b))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.18 ([#80](https://github.com/spectrocloud-labs/validator-plugin-network/issues/80)) ([f94ebc3](https://github.com/spectrocloud-labs/validator-plugin-network/commit/f94ebc3b5cb86ef94b3695e517549131ef488b07))
* typo in TCP conn rule failure ([6b51a54](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6b51a54b8fbbbd509a2b3124c325b04be5ff5b88))
* typo in ValidationResult names ([1b18967](https://github.com/spectrocloud-labs/validator-plugin-network/commit/1b18967633194c24a06eaa3e5f2477a06be78dc6))

## [0.0.5](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.4...v0.0.5) (2023-11-10)


### Bug Fixes

* **deps:** update module github.com/go-logr/logr to v1.3.0 ([#67](https://github.com/spectrocloud-labs/validator-plugin-network/issues/67)) ([4362ecd](https://github.com/spectrocloud-labs/validator-plugin-network/commit/4362ecdd33ef4427efdbd914d52b8b0840564de1))
* **deps:** update module github.com/onsi/gomega to v1.28.1 ([#64](https://github.com/spectrocloud-labs/validator-plugin-network/issues/64)) ([f3b43c7](https://github.com/spectrocloud-labs/validator-plugin-network/commit/f3b43c7af08941dba32c94bea505ab1ecbbc9e4e))
* **deps:** update module github.com/onsi/gomega to v1.29.0 ([#65](https://github.com/spectrocloud-labs/validator-plugin-network/issues/65)) ([75fde95](https://github.com/spectrocloud-labs/validator-plugin-network/commit/75fde955fdabbc659e2e8ae4148b4b21df72df25))
* **deps:** update module github.com/onsi/gomega to v1.30.0 ([#71](https://github.com/spectrocloud-labs/validator-plugin-network/issues/71)) ([806234a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/806234a462c83c83c3be0ee370e077444c2c8f72))


### Other

* add license ([d70ed62](https://github.com/spectrocloud-labs/validator-plugin-network/commit/d70ed62c20d16219cd55fdfc6dde002de91acb0f))
* better platform support for images ([#75](https://github.com/spectrocloud-labs/validator-plugin-network/issues/75)) ([037439a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/037439a33471fbec8635e815f8ccfa343c2690e2))
* **deps:** pin dependencies ([#73](https://github.com/spectrocloud-labs/validator-plugin-network/issues/73)) ([196783d](https://github.com/spectrocloud-labs/validator-plugin-network/commit/196783d1248b356f20a651cf8d4cdf0ee323693d))
* **deps:** update actions/checkout action to v4 ([#74](https://github.com/spectrocloud-labs/validator-plugin-network/issues/74)) ([214c3e8](https://github.com/spectrocloud-labs/validator-plugin-network/commit/214c3e8a29ee6d78c58f6f522ad71cbe759d6d97))
* **deps:** update gcr.io/kubebuilder/kube-rbac-proxy docker tag to v0.15.0 ([#62](https://github.com/spectrocloud-labs/validator-plugin-network/issues/62)) ([d4d00ff](https://github.com/spectrocloud-labs/validator-plugin-network/commit/d4d00ffee950688e4b1cb23682b4e83c11e7ae0e))
* **deps:** update google-github-actions/release-please-action digest to db8f2c6 ([#70](https://github.com/spectrocloud-labs/validator-plugin-network/issues/70)) ([2617096](https://github.com/spectrocloud-labs/validator-plugin-network/commit/2617096a9100c4dc327d448b29c2b6933080b7fe))
* **deps:** update helm/chart-testing-action action to v2.6.0 ([#68](https://github.com/spectrocloud-labs/validator-plugin-network/issues/68)) ([2ba587e](https://github.com/spectrocloud-labs/validator-plugin-network/commit/2ba587ea82312f82de4203233ca4a8a95bfe6b7d))
* **deps:** update helm/chart-testing-action action to v2.6.1 ([#69](https://github.com/spectrocloud-labs/validator-plugin-network/issues/69)) ([8613f0d](https://github.com/spectrocloud-labs/validator-plugin-network/commit/8613f0d838bed1b45753f64f39508d8a8130b7cb))
* tidy imports ([#66](https://github.com/spectrocloud-labs/validator-plugin-network/issues/66)) ([65a4d9b](https://github.com/spectrocloud-labs/validator-plugin-network/commit/65a4d9b771bfd1fa10c34ee8df6981a2914962b1))
* update validator core & align w/ VR changes ([#77](https://github.com/spectrocloud-labs/validator-plugin-network/issues/77)) ([929a56f](https://github.com/spectrocloud-labs/validator-plugin-network/commit/929a56f693140b2ec1a341ca853d03ad295fff70))

## [0.0.4](https://github.com/spectrocloud-labs/valid8or-plugin-network/compare/v0.0.3...v0.0.4) (2023-10-20)


### Bug Fixes

* ct lints ([7431656](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/74316561923f058f0c6013f60adf8543f97f349d))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.3 ([#59](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/59)) ([8e66c5c](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/8e66c5ca6a02808585b90df8eb5d1bbaa42dc15e))


### Other

* **deps:** bump golang.org/x/net from 0.16.0 to 0.17.0 ([#51](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/51)) ([24cc765](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/24cc7658ceb3d4acf9678b7eb58c246b89137d41))
* **deps:** update actions/checkout digest to b4ffde6 ([#57](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/57)) ([f0d1197](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/f0d1197713be76b934e534d8b9efc7c4168b552a))
* **deps:** update actions/setup-python digest to 65d7f2d ([#55](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/55)) ([dd9316b](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/dd9316baca87fe1345f567f5b992aca41bbddac8))
* **deps:** update gcr.io/kubebuilder/kube-rbac-proxy docker tag to v0.14.4 ([#49](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/49)) ([558bc60](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/558bc60847a23ce437c3c2b9595d056511db1bb3))
* **deps:** update google-github-actions/release-please-action digest to 4c5670f ([#56](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/56)) ([fb00c56](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/fb00c56cbe0c2028859374ef3ccc00e1b65db1dc))
* enable renovate automerges ([5723d6e](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/5723d6e01ac7b611260bb209b0e1d70117a4ac65))
* release 0.0.4 ([afdb2f0](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/afdb2f0ef424441ed8dd9628045ca842efdf674f))


### Refactoring

* valid8or -&gt; validator ([#60](https://github.com/spectrocloud-labs/valid8or-plugin-network/issues/60)) ([97f0c6e](https://github.com/spectrocloud-labs/valid8or-plugin-network/commit/97f0c6eae41dfe7d8760597de6c0fe79617a885c))

## [0.0.3](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.2...v0.0.3) (2023-10-10)


### Bug Fixes

* **deps:** update kubernetes packages to v0.28.2 ([#42](https://github.com/spectrocloud-labs/validator-plugin-network/issues/42)) ([8a4ff90](https://github.com/spectrocloud-labs/validator-plugin-network/commit/8a4ff905613d4e0d969718c396abe422aa9d8ef1))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.12.1 ([#43](https://github.com/spectrocloud-labs/validator-plugin-network/issues/43)) ([4796e94](https://github.com/spectrocloud-labs/validator-plugin-network/commit/4796e9437b070bdecb02033f7adecda0d0d92e52))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.0 ([#46](https://github.com/spectrocloud-labs/validator-plugin-network/issues/46)) ([6ee30d7](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6ee30d702ed4bf673dc8173d83a5ae4f92752753))
* **deps:** update module github.com/onsi/gomega to v1.28.0 ([#45](https://github.com/spectrocloud-labs/validator-plugin-network/issues/45)) ([99b204f](https://github.com/spectrocloud-labs/validator-plugin-network/commit/99b204f14d89980ac95ab518efde99d6c0c446c9))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.7 ([#32](https://github.com/spectrocloud-labs/validator-plugin-network/issues/32)) ([bf9c1e9](https://github.com/spectrocloud-labs/validator-plugin-network/commit/bf9c1e907c1f5e25a9897a40af98fbf010826477))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.8 ([#35](https://github.com/spectrocloud-labs/validator-plugin-network/issues/35)) ([6819847](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6819847b3d8ea080639a9be314460569782342de))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.9 ([#47](https://github.com/spectrocloud-labs/validator-plugin-network/issues/47)) ([18fe137](https://github.com/spectrocloud-labs/validator-plugin-network/commit/18fe13759068cb3ec504a4d138b34a0ec0be2280))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.2 ([#41](https://github.com/spectrocloud-labs/validator-plugin-network/issues/41)) ([81bc09e](https://github.com/spectrocloud-labs/validator-plugin-network/commit/81bc09e128a9fc5755e7b0c9e763dd85bc42e1d2))


### Other

* add YAML tags to Network validator spec ([1c2f57c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/1c2f57c3b628d99e00ce5883a56c0e17e0c782fd))
* **deps:** update actions/checkout digest to 8ade135 ([#44](https://github.com/spectrocloud-labs/validator-plugin-network/issues/44)) ([1723b5a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/1723b5a84b44f557958fc3271330159abbd99a19))
* **deps:** update actions/upload-artifact digest to a8a3f3a ([#31](https://github.com/spectrocloud-labs/validator-plugin-network/issues/31)) ([fdbca6c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/fdbca6cc4ac02320f82110134c60b0bf095fbda2))
* **deps:** update docker/build-push-action action to v5 ([#38](https://github.com/spectrocloud-labs/validator-plugin-network/issues/38)) ([ec02adb](https://github.com/spectrocloud-labs/validator-plugin-network/commit/ec02adbd568f0b828a27d2ef8b59eeb22ce50fad))
* **deps:** update docker/build-push-action digest to 0a97817 ([#36](https://github.com/spectrocloud-labs/validator-plugin-network/issues/36)) ([0f51409](https://github.com/spectrocloud-labs/validator-plugin-network/commit/0f5140971db8372e9d83f936183e3f6acf4e08ec))
* **deps:** update docker/login-action action to v3 ([#39](https://github.com/spectrocloud-labs/validator-plugin-network/issues/39)) ([6d256d6](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6d256d6c3c756f86c38cac056703563b15e6c458))
* **deps:** update docker/setup-buildx-action action to v3 ([#40](https://github.com/spectrocloud-labs/validator-plugin-network/issues/40)) ([cadafc3](https://github.com/spectrocloud-labs/validator-plugin-network/commit/cadafc33c69671a2bc662f7dab215df1bf6b0d41))
* release 0.0.3 ([860e9fc](https://github.com/spectrocloud-labs/validator-plugin-network/commit/860e9fcab1a1acbee4d0fef41f6a6b7a578689ce))

## [0.0.2](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.1...v0.0.2) (2023-09-06)


### Bug Fixes

* **deps:** update kubernetes packages to v0.28.1 ([#21](https://github.com/spectrocloud-labs/validator-plugin-network/issues/21)) ([bf92497](https://github.com/spectrocloud-labs/validator-plugin-network/commit/bf92497bb01aeb674af92943154050c88b366b61))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.12.0 ([#22](https://github.com/spectrocloud-labs/validator-plugin-network/issues/22)) ([f327be7](https://github.com/spectrocloud-labs/validator-plugin-network/commit/f327be73080ee2f51597473e6b138bdb1a8f5e58))
* **deps:** update module github.com/onsi/gomega to v1.27.10 ([#19](https://github.com/spectrocloud-labs/validator-plugin-network/issues/19)) ([9ab822c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/9ab822c231aa63c538cc95369bfffeb2231a21b0))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.5 ([#20](https://github.com/spectrocloud-labs/validator-plugin-network/issues/20)) ([b4ca38a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/b4ca38a579bb2e46819268c3e8a84c9557c1edd8))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.1 - abandoned ([#23](https://github.com/spectrocloud-labs/validator-plugin-network/issues/23)) ([59bb191](https://github.com/spectrocloud-labs/validator-plugin-network/commit/59bb191f3a812a6efdc392db6424448f957c95bd))
* remove metrics bind address from chart values ([a832a1c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/a832a1cf1b2bc9c357c87c3a15797e6f6b137aa7))


### Other

* **deps:** update actions/checkout action to v4 ([#28](https://github.com/spectrocloud-labs/validator-plugin-network/issues/28)) ([e40247c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/e40247c45192ac25c6a0fd1c1d2d282821eee3e3))
* **deps:** update actions/checkout digest to f43a0e5 ([#13](https://github.com/spectrocloud-labs/validator-plugin-network/issues/13)) ([cdc8139](https://github.com/spectrocloud-labs/validator-plugin-network/commit/cdc81390456334988680ca43cf831255af4f71cb))
* **deps:** update actions/setup-go digest to 93397be ([#14](https://github.com/spectrocloud-labs/validator-plugin-network/issues/14)) ([e5f29bd](https://github.com/spectrocloud-labs/validator-plugin-network/commit/e5f29bd94cd5cf93166c92161974ed42bda2303e))
* **deps:** update docker/setup-buildx-action digest to 885d146 ([#18](https://github.com/spectrocloud-labs/validator-plugin-network/issues/18)) ([26be287](https://github.com/spectrocloud-labs/validator-plugin-network/commit/26be28723eae14ac0ff0a4a2b95c97c17c051c8e))
* update validator ([c4de6ff](https://github.com/spectrocloud-labs/validator-plugin-network/commit/c4de6ffab4458eebbc61adeb713fcba0e164ad76))


### Docs

* update README.md ([9cdb991](https://github.com/spectrocloud-labs/validator-plugin-network/commit/9cdb9919345312ca9e52d751ff41cdc102957486))

## [0.0.1](https://github.com/spectrocloud-labs/validator-plugin-network/compare/v0.0.1...v0.0.1) (2023-08-31)


### Features

* add dns, icmp, iprange, mtu, tcp checks ([cab7c7a](https://github.com/spectrocloud-labs/validator-plugin-network/commit/cab7c7a34d6815572c3c37eeb799fca887ed850b))
* inital commit ([456a1fa](https://github.com/spectrocloud-labs/validator-plugin-network/commit/456a1faf45afb45c2604efbd4bae9872e8aa8e1b))


### Bug Fixes

* bump max len for chart.name, chart.chart ([6ecd682](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6ecd682053f99f6a6ec5a3a5caee55678163f66a))
* chart helpers generating invalid DNS names ([835402f](https://github.com/spectrocloud-labs/validator-plugin-network/commit/835402fc427e623fad8df635cb6fd0c5e0d4045d))
* securityContext blocking MTU check w/ ping ([beb2cf1](https://github.com/spectrocloud-labs/validator-plugin-network/commit/beb2cf1d940aeb6d3d07b022eb81e4c284e01da1))
* unique leader election id ([aa7920f](https://github.com/spectrocloud-labs/validator-plugin-network/commit/aa7920f45de86c85f9f05c96ecdada68a8f02780))
* update NetworkValidator CRD ([74e33e7](https://github.com/spectrocloud-labs/validator-plugin-network/commit/74e33e70d71441a4c0eaa2f28c1668bcc0e8a4fd))
* update RBAC template ([9460480](https://github.com/spectrocloud-labs/validator-plugin-network/commit/94604808b1a0c89fcaed5ef4c2c6c1dfe8fea250))


### Other

* release 0.0.1 ([78d491c](https://github.com/spectrocloud-labs/validator-plugin-network/commit/78d491cda744e2048673c912169539cc31b27d2f))
* update README ([6ba7fed](https://github.com/spectrocloud-labs/validator-plugin-network/commit/6ba7fed3c0e0c18ccfafe6a836f868a11f69b228))
