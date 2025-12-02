## [omni-infra-provider-bare-metal 0.7.1](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.7.1) (2025-12-02)

Welcome to the v0.7.1 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Utku Ozdemir

### Changes
<details><summary>1 commit</summary>
<p>

* [`821b331`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/821b331cec9fd69fe2ff848f2bda8af472df2a43) fix: always include the extra config docs in machine config
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.7.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.7.0)

## [omni-infra-provider-bare-metal 0.7.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.7.0) (2025-11-17)

Welcome to the v0.7.0 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Andrey Smirnov
* Mateusz Urbanek
* Noel Georgi
* Utku Ozdemir
* Justin Garrison
* Laura Brehm

### Changes
<details><summary>2 commits</summary>
<p>

* [`61f2a5d`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/61f2a5d55340ee8268901c98c568e61b1276dc83) chore: rekres, bump deps
* [`f303b3f`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/f303b3ff330d7de951f1e04ef7d4817113d9f578) feat: allow providing additional config documents from config endpoint
</p>
</details>

### Changes from siderolabs/gen
<details><summary>1 commit</summary>
<p>

* [`4c7388b`](https://github.com/siderolabs/gen/commit/4c7388b6a09d6a2ab6a380541df7a5b4bcc4b241) chore: update Go modules, replace YAML library
</p>
</details>

### Changes from siderolabs/talos
<details><summary>15 commits</summary>
<p>

* [`bc34de6e1`](https://github.com/siderolabs/talos/commit/bc34de6e1741969e873dd568054231acf4cb54fd) release(v1.11.5): prepare release
* [`3945c6c8f`](https://github.com/siderolabs/talos/commit/3945c6c8f029b20edcb3de0bf0a5e4c78023a403) feat: update containerd to 2.1.5
* [`8aec37684`](https://github.com/siderolabs/talos/commit/8aec376841aa910c960f2aea0ffd8a100cc2575b) release(v1.11.4): prepare release
* [`9c27f9e62`](https://github.com/siderolabs/talos/commit/9c27f9e62097db284961aa7014e0bef14401f97f) fix: race between VolumeConfigController and UserVolumeConfigController
* [`ac27129b1`](https://github.com/siderolabs/talos/commit/ac27129b19485142eb76a04eee4b372d1cabcdaf) fix: provide minimal platform metadata always
* [`19463323e`](https://github.com/siderolabs/talos/commit/19463323eb77b3b0ea51df2793853723185fbbbc) fix: image-signer commands
* [`62aa09644`](https://github.com/siderolabs/talos/commit/62aa09644196ae6a551168530f42884bc78e00f2) chore: update dependencies
* [`075f9ef22`](https://github.com/siderolabs/talos/commit/075f9ef22ffb61710165456313c4173d9765641d) fix: userspace wireguard handling
* [`35b97016c`](https://github.com/siderolabs/talos/commit/35b97016c02b08163bc230e1728e35e61e11418d) fix: log duplication on log senders
* [`d00754e35`](https://github.com/siderolabs/talos/commit/d00754e35b365ac45c40f62af45a74f38e5ccfd6) fix: add video kernel module to arm
* [`89bca7590`](https://github.com/siderolabs/talos/commit/89bca759000c11fa7c59e0c9045816c20858067b) fix: set a timeout for SideroLink provision API call
* [`23b21eb90`](https://github.com/siderolabs/talos/commit/23b21eb90b05d8ebb4adc71fb4a269c1b4049d8a) fix: imager build on arm64
* [`2a4f1771c`](https://github.com/siderolabs/talos/commit/2a4f1771c632476b1a6569e29bb1043c480ea349) feat: use image signer
* [`e043e1bc0`](https://github.com/siderolabs/talos/commit/e043e1bc004ed80a93809937096b5e5c59909704) chore: push `latest` tag only on main
* [`8edddafcd`](https://github.com/siderolabs/talos/commit/8edddafcd97b868df1c8e78cecf1eae70f0eaf83) fix: reserve the apid and trustd ports from the ephemeral port range
</p>
</details>

### Dependency Changes

* **github.com/cosi-project/runtime**                  v1.11.0 -> v1.12.0
* **github.com/grpc-ecosystem/go-grpc-middleware/v2**  v2.3.2 -> v2.3.3
* **github.com/insomniacslk/dhcp**                     da879a2c3546 -> 175e84fbb167
* **github.com/klauspost/compress**                    v1.18.0 -> v1.18.1
* **github.com/siderolabs/gen**                        v0.8.5 -> v0.8.6
* **github.com/siderolabs/omni/client**                v1.2.1 -> v1.3.2
* **github.com/siderolabs/talos**                      v1.11.3 -> v1.11.5
* **github.com/siderolabs/talos/pkg/machinery**        v1.11.3 -> v1.12.0-beta.0
* **golang.org/x/net**                                 v0.46.0 -> v0.47.0
* **golang.org/x/sync**                                v0.17.0 -> v0.18.0

Previous release can be found at [v0.6.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.6.0)

## [omni-infra-provider-bare-metal 0.6.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.6.0) (2025-11-07)

Welcome to the v0.6.0 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Utku Ozdemir

### Changes
<details><summary>1 commit</summary>
<p>

* [`9c50645`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/9c50645668285df586dcb0e916a02478a3865c03) feat: allow specifying a custom DHCP proxy port
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.5.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.5.0)

## [omni-infra-provider-bare-metal 0.5.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.5.0) (2025-10-17)

Welcome to the v0.5.0 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Andrey Smirnov
* Mateusz Urbanek
* Noel Georgi
* Dmitrii Sharshakov
* Oguz Kilcan
* Utku Ozdemir
* Alp Celik
* Amarachi Iheanacho
* Andrew Longwill
* Chris Sanders
* Grzegorz Rozniecki
* Markus Freitag
* Olivier Doucet
* Orzelius
* Serge Logvinov

### Changes
<details><summary>3 commits</summary>
<p>

* [`4e7f89b`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/4e7f89b560cb7d733053a7d870e00b6efadcb886) chore: bump Talos version to 1.11.3, make integration tests parallel
* [`e53367e`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/e53367edc0e9fe5ea11b1ca8c219aa0bd500bd90) chore: rekres, bump deps
* [`5bca2eb`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/5bca2eb8c19c5626e26edc9593fb3b5e3acd45fd) refactor: adapt to new QTransform controllers
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>2 commits</summary>
<p>

* [`4154a77`](https://github.com/siderolabs/crypto/commit/4154a771b09f0023e0d258bba6aecc29febabecb) feat: implement dynamic certificate reloader
* [`dae07fa`](https://github.com/siderolabs/crypto/commit/dae07fa14f963b34ea67abf0cbc50ba24f280524) chore: update to Go 1.25
</p>
</details>

### Changes from siderolabs/image-factory
<details><summary>20 commits</summary>
<p>

* [`a3a7661`](https://github.com/siderolabs/image-factory/commit/a3a7661df37083c3af0a929265a424f003c9db1a) release(v0.8.4): prepare release
* [`075aa3f`](https://github.com/siderolabs/image-factory/commit/075aa3fa0c10abc4e06d2be1d3f3a394e56d1947) fix: update Talos to 1.11.1
* [`02723cd`](https://github.com/siderolabs/image-factory/commit/02723cdf6b96b106b3a961f1eb88731366e0cecb) fix: translation ID
* [`94c6df1`](https://github.com/siderolabs/image-factory/commit/94c6df1f3497c5a4173fa3ddfd3169b65d70dc15) release(v0.8.3): prepare release
* [`7254abf`](https://github.com/siderolabs/image-factory/commit/7254abf251c3a1140a220969ac9bd684c55f8774) fix: disable redirects to PXE
* [`251aee0`](https://github.com/siderolabs/image-factory/commit/251aee03710e8c3603a9f4cf9677353a62e860ea) release(v0.8.2): prepare release
* [`418eebb`](https://github.com/siderolabs/image-factory/commit/418eebb19ff7a6948a8125db2461f257612fcd23) fix: don't filter out `rc` versions
* [`57ad419`](https://github.com/siderolabs/image-factory/commit/57ad419a199bcd9956ba8aa48db451e1ce3c61d5) release(v0.8.1): prepare release
* [`6392086`](https://github.com/siderolabs/image-factory/commit/63920865fa4bd1f4537880e5b491e685a88fd965) fix: prevent failure on cache.Get
* [`a1e3707`](https://github.com/siderolabs/image-factory/commit/a1e37078e10bae58d8ee3f117cdbc405de35e65c) feat: add fallback if S3 is missbehaving
* [`9760ab0`](https://github.com/siderolabs/image-factory/commit/9760ab0fee7196885f50a92abf872c5c94f3dd2c) release(v0.8.0): prepare release
* [`7c6d261`](https://github.com/siderolabs/image-factory/commit/7c6d26184cd3a6f903385230fcbddc92cf67d065) fix: set content-disposition on S3
* [`f3e97df`](https://github.com/siderolabs/image-factory/commit/f3e97df4e609aa1b6ffc39d6b4cb8c76e891669e) docs(image-factory): add info about S3 cache and CDN
* [`d25e7ac`](https://github.com/siderolabs/image-factory/commit/d25e7acdc3b9e0a1fb96a0013133fc8e89097d1b) fix: add extra context to logs from s3 cache
* [`a3a0dff`](https://github.com/siderolabs/image-factory/commit/a3a0dff1f8846a2373a63d428ea86717bbdc452f) fix: add optional region to S3 client
* [`a9e2d08`](https://github.com/siderolabs/image-factory/commit/a9e2d08b1162c0e470b87da8e6ad448b34426d7a) feat: add support for Object Storage and CDN cache
* [`b8bfc19`](https://github.com/siderolabs/image-factory/commit/b8bfc1985c4c93cd1aa12a251deaa1ecb6239d20) docs: add air-gapped documentation
* [`f8b4ef0`](https://github.com/siderolabs/image-factory/commit/f8b4ef0ea538b56238b9ea0a51daadf5d5999ae6) docs: add new translation
* [`0c83228`](https://github.com/siderolabs/image-factory/commit/0c83228ae5ad0349f376f56743a8d3b8e2858ec4) release(v0.7.6): prepare release
* [`6f409ec`](https://github.com/siderolabs/image-factory/commit/6f409ecd914094afe9293a23883806798a0cc5dd) fix: drop extractParams function
</p>
</details>

### Changes from siderolabs/talos
<details><summary>92 commits</summary>
<p>

* [`a0243ef77`](https://github.com/siderolabs/talos/commit/a0243ef77e6532ed2919689d305eeaf97458c0a1) release(v1.11.3): prepare release
* [`560241c00`](https://github.com/siderolabs/talos/commit/560241c00e0e9fdcd3ad614a28183f83407c07e5) fix: make Akamai platform usable
* [`1b23cad61`](https://github.com/siderolabs/talos/commit/1b23cad61cafcfa9130ef216e85df07716ca8a8a) fix: cherry-pick of commit `0fbb0b0` from #11959
* [`876719a92`](https://github.com/siderolabs/talos/commit/876719a92d4e4dfe8dfdd4d81c0671cf40e7bd45) fix: cherry-pick of commit `cd9fb27` from #11943
* [`9a30ab6f5`](https://github.com/siderolabs/talos/commit/9a30ab6f5cd418636258cc2812aecfe3e7bf9ee5) feat: bump go, kernel and runc
* [`0fbb0b028`](https://github.com/siderolabs/talos/commit/0fbb0b0280c1f8a4da954237e765c7682cea4402) fix: provide nocloud metadata with missing network config
* [`0dad32819`](https://github.com/siderolabs/talos/commit/0dad328195190b579ac33a6ce10c38847889469a) feat: update Flannel to v0.27.4
* [`49182b386`](https://github.com/siderolabs/talos/commit/49182b386b983814c6356dc21acd05a9a210bca3) fix: support secure HTTP proxy with gRPC dial
* [`a460f5726`](https://github.com/siderolabs/talos/commit/a460f572693726b5b13528759afd6c9a2f57f3fd) feat: update etcd 3.6.5, CoreDNS 1.12.4
* [`48ee8581b`](https://github.com/siderolabs/talos/commit/48ee8581bc5b0808bf70e7cdcdb38e5cf73695de) fix: don't set broadcast for /31 and /32 addresses
* [`7668c52dd`](https://github.com/siderolabs/talos/commit/7668c52dd4126e0637d42dbf54b005e170051c91) fix: provide refreshing CA pool (resolvers)
* [`511b4d2e8`](https://github.com/siderolabs/talos/commit/511b4d2e89320f79f66cd3f0f18db1a01e3f4aef) release(v1.11.2): prepare release
* [`ac452574e`](https://github.com/siderolabs/talos/commit/ac452574e79ef3564e622d44fd4516681740c8cf) fix: default gateway as string
* [`7cec0e042`](https://github.com/siderolabs/talos/commit/7cec0e0420d613910d0d90c542e8f00ff3cfc9b5) fix: uefi boot entry handling logic
* [`637154ed2`](https://github.com/siderolabs/talos/commit/637154ed2555a885a1de9dfdf14813b9b807fb38) docs: drop invalid v1.12 docs
* [`a6d2f65a6`](https://github.com/siderolabs/talos/commit/a6d2f65a61065285366dc3698a2b5d556dde8da0) chore(ci): rekres to use new runner groups
* [`cd82ee204`](https://github.com/siderolabs/talos/commit/cd82ee204eda75dd09cedd85b2414edebacfb5ca) refactor: efivarfs mock and tests
* [`996d97de6`](https://github.com/siderolabs/talos/commit/996d97de6e1fd5feea4e1052e0d1c6f6c0f3c6f9) chore: update pkgs
* [`bbf860c5c`](https://github.com/siderolabs/talos/commit/bbf860c5ccbdd2fdc877459d05b2f64b9c127a5d) docs: update component updates
* [`24c1bcecf`](https://github.com/siderolabs/talos/commit/24c1bcecf5d1fd82e24bf85a48ae3f966aedec2d) fix: bump trustd memory limit
* [`56d6d6f75`](https://github.com/siderolabs/talos/commit/56d6d6f755d35785f7be9665813e5847c7dfb14c) chore: pass in github token to imager
* [`682df89d7`](https://github.com/siderolabs/talos/commit/682df89d78312b7a56d017c953397d171aee4a37) fix: use correct order to determine SideroV1 keys directory path
* [`a838881fa`](https://github.com/siderolabs/talos/commit/a838881fafcdfe20b3ccb40b5535cc27946b19ea) fix: trim zero bytes in the DHCP host & domain response
* [`9c962ae9c`](https://github.com/siderolabs/talos/commit/9c962ae9c86168eb71677a7ce678a3a443d64f40) fix: re-create cgroups when restarting runners
* [`de243f9ae`](https://github.com/siderolabs/talos/commit/de243f9aede933336d7ca48937df40d168d5257e) test: fix flakiness in RawVolumes test
* [`ec8fde596`](https://github.com/siderolabs/talos/commit/ec8fde596fac2058b205fe84026355d6220e31dc) feat: update Kubernetes to 1.34.1
* [`797897dfb`](https://github.com/siderolabs/talos/commit/797897dfbf050b0b81a018ace9ac77de45b17410) test: improve test stability
* [`98273666e`](https://github.com/siderolabs/talos/commit/98273666e8ed9fd8a94b66bd3834bf78ecbc44c8) feat: update runc to 1.3.1
* [`8e85c8362`](https://github.com/siderolabs/talos/commit/8e85c83625502e08c058b865c123b0828a90fed6) release(v1.11.1): prepare release
* [`ff8644cd2`](https://github.com/siderolabs/talos/commit/ff8644cd2efefe00ef469f180392eb9fa63b8a52) fix: correctly handle status-code 204
* [`7d5fe2d0f`](https://github.com/siderolabs/talos/commit/7d5fe2d0f1d5761d5aba28c55999bd8222ef5e3f) feat: update Linux kernel (memcg_v1, ublk)
* [`9e310a9dd`](https://github.com/siderolabs/talos/commit/9e310a9dd9e70669c46900f6950c29929a308261) fix: enable support for VMWare arm64
* [`f7620f028`](https://github.com/siderolabs/talos/commit/f7620f02817b271686024799353b87f5f51c3cf7) feat: update CoreDNS to 1.12.3
* [`01bf2f6f9`](https://github.com/siderolabs/talos/commit/01bf2f6f9d203dad55910bdde3539e883b138f8e) feat: add SOCKS5 proxy support to dynamic proxy dialer
* [`8a578bc4a`](https://github.com/siderolabs/talos/commit/8a578bc4ac95fc543f0564281d1a6a54f3299061) feat: update Linux to 6.12.45
* [`d9d89a3a8`](https://github.com/siderolabs/talos/commit/d9d89a3a82be5e5a276b1a3328bc0daefbbff5d6) release(v1.11.0): prepare release
* [`364b48690`](https://github.com/siderolabs/talos/commit/364b4869004fde1ffed27e50b657be41c2127621) feat: update pkgs/tools for pcre2 10.46
* [`be70ea03f`](https://github.com/siderolabs/talos/commit/be70ea03fcf7aa8dd57eda966ed5445a8be91e37) feat: update pkgs for NVIDIA prod 570.172.08
* [`a5f80b4fe`](https://github.com/siderolabs/talos/commit/a5f80b4fe6dad879ed875cb6763a76223187259c) fix: bring back linux/armv7 build and update xz
* [`751dae432`](https://github.com/siderolabs/talos/commit/751dae432611b438b140aec5fc14c7f9734d4e87) fix: drop linux/armv7 build
* [`8cbd75320`](https://github.com/siderolabs/talos/commit/8cbd7532053d86cf71def0dab798401d4795aeb4) fix: update xz module (security)
* [`803ed1ef9`](https://github.com/siderolabs/talos/commit/803ed1ef96c0213352fac3d8c48a9f23cd0a9aa7) feat: update Kubernetes to 1.34.0
* [`a80898da9`](https://github.com/siderolabs/talos/commit/a80898da9d1219f6c8acc9f33f3d83e3856bd497) feat: update Linux to 6.12.43
* [`30c14aa71`](https://github.com/siderolabs/talos/commit/30c14aa71d33a5f70ddb35efc3840a3c5e23743a) feat: update Kubernetes default to v1.34.0-rc.2
* [`ed7d8cbac`](https://github.com/siderolabs/talos/commit/ed7d8cbac0aa388820adc217c5af647ada9d99d6) docs: link to kubernetes linux swap tuning
* [`1ee82120e`](https://github.com/siderolabs/talos/commit/1ee82120e96e1aa5bc6880ab77031a59a092ec6c) docs: apply fixes for what is new
* [`36102eae1`](https://github.com/siderolabs/talos/commit/36102eae179a9beed634c1faca1778de18b97ad1) release(v1.11.0-rc.0): prepare release
* [`0f22913d9`](https://github.com/siderolabs/talos/commit/0f22913d96e7088aaff697c7fd93cd7eb64240cb) fix: image cache lockup on a missing volume
* [`46cf25c7c`](https://github.com/siderolabs/talos/commit/46cf25c7c0b570faa307ee64ab46cf96db0e210d) feat: update Linux to 6.12.41
* [`62f6c97fe`](https://github.com/siderolabs/talos/commit/62f6c97fe6430a1c4b2dd78273a7b0718ea89462) fix: provide mitigation CVE-1999-0524
* [`350319063`](https://github.com/siderolabs/talos/commit/3503190637042083fff169a46bbdbe1cfd750c73) fix: actually use SIDEROV1_KEYS_DIR env var if it's provided
* [`430a27dc2`](https://github.com/siderolabs/talos/commit/430a27dc24b42c3dc7c8f6e04e128544bca39feb) fix: kubernetes upgrade options for kubelet
* [`e3a9097c4`](https://github.com/siderolabs/talos/commit/e3a9097c4fb99dceae69740fd43dcaeb4ac9da32) fix: set secs field in DHCPv4 packets
* [`babddd0e4`](https://github.com/siderolabs/talos/commit/babddd0e400386d7e8dbab806cb1724ca105dc4d) fix: dial with proxy
* [`23efda4db`](https://github.com/siderolabs/talos/commit/23efda4dbfbb135c81f538a433ee53ecc7c64a52) feat: use key provider with fallback option for auth type SideroV1
* [`e2a5a9b3f`](https://github.com/siderolabs/talos/commit/e2a5a9b3fe6f7eb2b44761c2bbedd2a9d183bcdc) chore: re-enable vulncheck
* [`f5d700a0c`](https://github.com/siderolabs/talos/commit/f5d700a0c6d5f99573a57cce871eb25a8c14b464) release(v1.11.0-beta.2): prepare release
* [`6186d1821`](https://github.com/siderolabs/talos/commit/6186d182189d229e3065631076f435d34bfc4f53) chore: disable vulncheck temporarily
* [`e4a2a8d9c`](https://github.com/siderolabs/talos/commit/e4a2a8d9c09f810e35923e4641db8921e6f85981) feat: update default Kubernetes to v1.34.0-rc.1
* [`4c4236d7e`](https://github.com/siderolabs/talos/commit/4c4236d7eb53185704f83667a27d191577a438e0) feat: update Go to 1.24.6
* [`a01a390f6`](https://github.com/siderolabs/talos/commit/a01a390f692bad314dacb84eaa06ac3b78034243) chore: add deadcode elimination linter
* [`49fad0ede`](https://github.com/siderolabs/talos/commit/49fad0ede4f8df9596fc3d6e4bff0a5fa89e2ea4) feat: add a pause function to dashboard
* [`21e8e9dc9`](https://github.com/siderolabs/talos/commit/21e8e9dc9ab1ec8c3550b6edd5c6c5b4e000e060) refactor: replace containerd/containerd/v2 module for proper DCE
* [`bbd01b6b7`](https://github.com/siderolabs/talos/commit/bbd01b6b7893d0d2004bdb9491d0f811f07c2ad3) refactor: fix deadcode elimination with godbus
* [`e8d9c81cc`](https://github.com/siderolabs/talos/commit/e8d9c81cc1b71827066442a9a26b387bb91202ba) refactor: stop using `text/template` in `machined` code paths
* [`85589662a`](https://github.com/siderolabs/talos/commit/85589662aadd34f1d3279b387bc3588adee21971) fix: unmarshal encryption STATE from META
* [`f10a626d2`](https://github.com/siderolabs/talos/commit/f10a626d2d5a8cfc612beabc1e74d87c35242bcc) docs: add what is new notes for 1.11
* [`5a15ce88b`](https://github.com/siderolabs/talos/commit/5a15ce88b62e0dd724954264f6ffd9f677463bae) release(v1.11.0-beta.1): prepare release
* [`614ca2e22`](https://github.com/siderolabs/talos/commit/614ca2e229c2e07ba664edbfd076a008eaebb894) fix: one more attempt to fix volume mount race on restart
* [`4b86dfe6f`](https://github.com/siderolabs/talos/commit/4b86dfe6fd0b7d55869c85816bb01b073817cc8f) feat: implement encryption locking to STATE
* [`8ae76c320`](https://github.com/siderolabs/talos/commit/8ae76c320c6115991c967ed946baaf9e8eb31d6d) feat: implement `talos.config.early` command line arg
* [`19f8c605e`](https://github.com/siderolabs/talos/commit/19f8c605ed0d0aecc80fdba646bac1d23539c1ca) docs: remove talos API flags from mgmt commands
* [`fa1d6fef8`](https://github.com/siderolabs/talos/commit/fa1d6fef8d664da263fe3b6dd2f59d83f2139ccc) feat: bootedentry resource
* [`7dee810d4`](https://github.com/siderolabs/talos/commit/7dee810d483155b9d9000eed30ec909efb441b90) fix: live reload of TLS client config for discovery client
* [`a5dc22466`](https://github.com/siderolabs/talos/commit/a5dc22466f2ab3fd9f32f0a4467c96ce075b3bec) fix: enforce minimum size on user volumes if not set explicitly
* [`7836e924d`](https://github.com/siderolabs/talos/commit/7836e924d4efc86fd6692915ebdc255d7d5545cc) feat: update containerd to 2.1.4
* [`5012550ec`](https://github.com/siderolabs/talos/commit/5012550ec7bbedf172dd7e8a6821c277f56fcb01) feat: add F71808E watchdog driver
* [`10ddc4cdd`](https://github.com/siderolabs/talos/commit/10ddc4cdd4aedc5101ea1f513ae72f2d5c752507) fix: grype scan
* [`d108e0a08`](https://github.com/siderolabs/talos/commit/d108e0a083720a5d3e059961afd0c2cb0a126d8a) fix(ci): use a random suffix for ami names
* [`504225546`](https://github.com/siderolabs/talos/commit/504225546252880af4506291b5ce6b4e9dac50f2) fix: issues with reading GPT
* [`bdaf08dd4`](https://github.com/siderolabs/talos/commit/bdaf08dd4fdb0a1c015685195e549c913c5fa824) feat: update PCI DB module to v0.3.2
* [`667dcebec`](https://github.com/siderolabs/talos/commit/667dcebec2b24f9bcb1bef1df4bb1a1c6219d78c) test: wait for service account test job longer
* [`ae176a4b7`](https://github.com/siderolabs/talos/commit/ae176a4b766f123a82c85a9418dfca70a8d09180) feat: update etcd to 3.6.4
* [`201b6801f`](https://github.com/siderolabs/talos/commit/201b6801f6651aa4bb43a6720109a2820d174714) fix: issue with volume remount on service restart
* [`2a911402b`](https://github.com/siderolabs/talos/commit/2a911402b5dd241b38a2dd7c2e3dc078acee7008) chore: tag aws snapshots created via ci with the image name
* [`d8bd84b56`](https://github.com/siderolabs/talos/commit/d8bd84b56cd0de0daab379ab9b9ee5ce3e99ac14) docs: add SBOM documentation
* [`7eec61993`](https://github.com/siderolabs/talos/commit/7eec61993296c33fa8d150e3ce6408313de3e912) feat: unify disk encryption configuration
* [`4ff2bf9e0`](https://github.com/siderolabs/talos/commit/4ff2bf9e06a5666fcd92257622699eec9b7a613d) feat: update etcd to v3.5.22
* [`31a67d379`](https://github.com/siderolabs/talos/commit/31a67d379627963b439d3705eacfe33424ba0d03) fix: do not download artifacts for cron Grype scan
* [`c6b6e0bb3`](https://github.com/siderolabs/talos/commit/c6b6e0bb3e258d1812a8f76ea488969862c6ea0c) docs: rewrite the getting started and prod docs for v1.10 and v1.11
* [`ca1c656e6`](https://github.com/siderolabs/talos/commit/ca1c656e6176546022b5a6a64370aad5d6c0c634) chore(ci): add more nvidia test matrix
* [`7a2e0f068`](https://github.com/siderolabs/talos/commit/7a2e0f068ea696aab21eec40d90b5f2ce3ebbe8b) feat: sync pkgs, update Linux to 6.12.40
</p>
</details>

### Dependency Changes

* **github.com/bougou/go-ipmi**                  v0.7.7 -> v0.7.8
* **github.com/cosi-project/runtime**            v1.10.7 -> v1.11.0
* **github.com/insomniacslk/dhcp**               5f8cf70e8c5f -> da879a2c3546
* **github.com/pin/tftp/v3**                     v3.1.0 -> 17016b3c2849
* **github.com/siderolabs/crypto**               v0.6.3 -> v0.6.4
* **github.com/siderolabs/image-factory**        v0.7.5 -> v0.8.4
* **github.com/siderolabs/omni/client**          da3f28f6b1f0 -> v1.2.1
* **github.com/siderolabs/talos**                v1.11.0-beta.0 -> v1.11.3
* **github.com/siderolabs/talos/pkg/machinery**  v1.11.0-beta.0 -> v1.11.3
* **github.com/spf13/cobra**                     v1.9.1 -> v1.10.1
* **github.com/stretchr/testify**                v1.10.0 -> v1.11.1
* **golang.org/x/net**                           v0.42.0 -> v0.46.0
* **golang.org/x/sync**                          v0.16.0 -> v0.17.0
* **google.golang.org/grpc**                     v1.74.2 -> v1.76.0
* **google.golang.org/protobuf**                 v1.36.6 -> v1.36.10

Previous release can be found at [v0.4.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.4.0)

