## [omni-infra-provider-bare-metal 0.8.1](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.8.1) (2026-03-03)

Welcome to the v0.8.1 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Andrey Smirnov
* Mateusz Urbanek
* Noel Georgi
* Kevin Tijssen
* Dmitrii Sharshakov
* Laura Brehm
* Orzelius
* Utku Ozdemir
* Artem Chernyshev
* Edward Sammut Alessi
* Tim Jones
* Bryan Lee
* Max Makarov
* Pranav Patil
* Alexis La Goutte
* Andreas Freund
* Andrei Kvapil
* Christopher Puschmann
* Daddie0
* Daniil Kivenko
* Florian Ströger
* Fritz Schaal
* Jan Paul
* Jonas Lammler
* Justin Garrison
* Lennard Klein
* Matthew Sanabria
* Mickaël Canévet
* Mikolaj Pawlikowski
* Nico Berlee
* Olav Thoresen
* Skye Soss
* Spencer Smith
* Sébastien Masset
* dataprolet
* drew

### Changes
<details><summary>3 commits</summary>
<p>

* [`9ed9832`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/9ed98329b94f69138635508425c48f30b884cf56) chore: bump deps including new redfish changes
* [`89de24a`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/89de24a5545c9e5fc21e09e3cca2791993b1f9f5) fix: correct /tftp/ HTTP path, debug build tags, and comment typos
* [`529806a`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/529806a8616bcc78cca7ce3d963b4123e34975ff) fix: use a private IP range in tests
</p>
</details>

### Changes from siderolabs/image-factory
<details><summary>50 commits</summary>
<p>

* [`f0c7a7b`](https://github.com/siderolabs/image-factory/commit/f0c7a7b53ce86c49c1531e1f1fd15c5bf3f00a70) release(v1.0.3): prepare release
* [`dd92631`](https://github.com/siderolabs/image-factory/commit/dd926314f61e4bbd8797c0302e4a8a14b9d693fb) docs: correct path to hack/copy-artifacts.sh
* [`ddc1a83`](https://github.com/siderolabs/image-factory/commit/ddc1a8389189e77e3f3679927ce550e9549a3e48) fix: update Talos to fix rpi_5 build
* [`b3d07e5`](https://github.com/siderolabs/image-factory/commit/b3d07e5e38da475018493adb5106eca9348de517) docs: remove redundant Kubernetes version prerequisite
* [`9666795`](https://github.com/siderolabs/image-factory/commit/96667959f60f6c4b6b670affdedc6ea898f6cfb2) fix: values.schema.json
* [`8a8da46`](https://github.com/siderolabs/image-factory/commit/8a8da46331b9dcd6353e93879f63c6a422d8d035) feat: adjust security context for user namespace mode
* [`bc631dc`](https://github.com/siderolabs/image-factory/commit/bc631dc3f9515bdbeabcb903190291805625ed9c) fix: values.schema.json
* [`8ea6fe9`](https://github.com/siderolabs/image-factory/commit/8ea6fe9eccba498f761061a4842616f58566e68e) feat: add user namespace support with Kubernetes version validation
* [`324c464`](https://github.com/siderolabs/image-factory/commit/324c464e22fff6ae13f4a199d7664229f628f07a) fix: skip initializing TUF if keyless signing is disabled
* [`a42b9d9`](https://github.com/siderolabs/image-factory/commit/a42b9d91c35d28817e43c8abc794f4cb3e7ae429) release(v1.0.2): prepare release
* [`80d1ba3`](https://github.com/siderolabs/image-factory/commit/80d1ba3e0e2a94f86cd37f80786743d928eb2b24) fix: pass nameoptions to verify bundle too
* [`eec01d1`](https://github.com/siderolabs/image-factory/commit/eec01d1d0351b31faa7357389933589d16d3dc04) release(v1.0.1): prepare release
* [`ec1c0a7`](https://github.com/siderolabs/image-factory/commit/ec1c0a790c99c55fcc3c315429f96e266bee7343) fix: pass insecure to the cosign new bundle verifier
* [`14d0f2a`](https://github.com/siderolabs/image-factory/commit/14d0f2a1fa2d40448c28c51a8d3c57caa44d5bbf) release(v1.0.0): prepare release
* [`a90529c`](https://github.com/siderolabs/image-factory/commit/a90529cc0066cd4c401b6e97bb69becfdcffc4f7) feat: add more security contexts
* [`ec69fe2`](https://github.com/siderolabs/image-factory/commit/ec69fe25da648422ef6a414bcddacab36e275579) fix: extra kernel args for overlays
* [`aa325ee`](https://github.com/siderolabs/image-factory/commit/aa325ee4ffe3f5cc8ed818027e0b31e055d7fcdf) feat: add Helm docs and schema
* [`3c18e05`](https://github.com/siderolabs/image-factory/commit/3c18e053c118131006b86affa7d0bb2af754cf95) feat: add Sidero google service account email also to verfiers
* [`151feb5`](https://github.com/siderolabs/image-factory/commit/151feb5589624ad2a8365a9a056e08c4d6780b2c) fix: docs url
* [`42a1c45`](https://github.com/siderolabs/image-factory/commit/42a1c45849be02ca572ee6f66b875135d49b4805) feat: add helm to kres
* [`ac4718a`](https://github.com/siderolabs/image-factory/commit/ac4718a617f88bfcfbe28edb2ee03a997fd19f7a) feat: update Talos and pkgs
* [`1d6468e`](https://github.com/siderolabs/image-factory/commit/1d6468ee6daac0eddd0ae01cdd47d83190f6b9d0) feat: add helm e2e to CI
* [`2f0499c`](https://github.com/siderolabs/image-factory/commit/2f0499cc73b5c20ad806a26bc5392eafbf03a87c) feat: added e2e tests
* [`2eccf98`](https://github.com/siderolabs/image-factory/commit/2eccf98ad5eb3fc79f13bfc4a712743d09fc65ca) fix: made changes on the recommendation of copilot
* [`e27ea36`](https://github.com/siderolabs/image-factory/commit/e27ea3647da994a34394635de869dfbaf7070f3a) feat: Added E2E with KUTTL
* [`9f6b9e7`](https://github.com/siderolabs/image-factory/commit/9f6b9e79665192daaa54efab4bdbfe09424569db) feat: Added additional tests
* [`4939747`](https://github.com/siderolabs/image-factory/commit/49397476eac0f33a6fc489355d05b80004953c1f) feat: Added helm unittests
* [`dcaa1db`](https://github.com/siderolabs/image-factory/commit/dcaa1db583160b605f89cfbf2f1a8ef36c59618b) feat: added helmchart
* [`1f85622`](https://github.com/siderolabs/image-factory/commit/1f85622c69e8a5a6401b16a5e41d9f04fc6a8267) feat: add cloudflare credentials helper
* [`852856d`](https://github.com/siderolabs/image-factory/commit/852856dc9d8e3db6a7b167626b8144e890e75f20) fix: installer internal config
* [`c8c6576`](https://github.com/siderolabs/image-factory/commit/c8c657680b2b55a630352ac5a1764342d608fd9c) release(v1.0.0-beta.0): prepare release
* [`56bd21b`](https://github.com/siderolabs/image-factory/commit/56bd21baa70dfadac318d409bc8ecf74a2b1a3c6) fix: allow `Cache-Control` header in CORS
* [`83f4d91`](https://github.com/siderolabs/image-factory/commit/83f4d91a063c56dd5d45c8674d5c67d51f14388e) fix: clarify bootloader selection
* [`c8c5faa`](https://github.com/siderolabs/image-factory/commit/c8c5faa6153dded74bb6734ae0812b7cee5ed201) feat: allow using image GET/HEAD API by the JS code on any domains
* [`e732d90`](https://github.com/siderolabs/image-factory/commit/e732d90618033f734fc5a5f9571537b7a0779e92) feat: support acm for secureboot
* [`5f103c1`](https://github.com/siderolabs/image-factory/commit/5f103c16c4854b0fbb89b60eb8fd1c1e6418197c) feat: support copying to clipboard
* [`c3532c4`](https://github.com/siderolabs/image-factory/commit/c3532c48692d2ab61b1b6af9b396dd95c312ea20) feat: update Talos with GRUB and other fixes
* [`b5ba663`](https://github.com/siderolabs/image-factory/commit/b5ba6630ed93021b6a4b820e23200aba3858c60f) fix: avoid pulling Talos core in schematic pkg
* [`b2b0cc8`](https://github.com/siderolabs/image-factory/commit/b2b0cc8561957b2b9a7af936c2c40217b43cae6b) fix: update cosign to v3.0.4
* [`fca99d0`](https://github.com/siderolabs/image-factory/commit/fca99d01a5be765ad0c853ac37a3e02e581dc824) chore: update `docs/developing.md`
* [`49f4226`](https://github.com/siderolabs/image-factory/commit/49f42261588d07a2728a1393022835d49426a609) chore: separate kres integration-test variables
* [`190aa22`](https://github.com/siderolabs/image-factory/commit/190aa22d6265ba318696e424a15ade8da55aa87b) fix: add missing libarchive dependency
* [`37bd795`](https://github.com/siderolabs/image-factory/commit/37bd7954478cccba2664add5e228540969f6aea3) fix: image-factory rootless
* [`99cbfd7`](https://github.com/siderolabs/image-factory/commit/99cbfd73d4ed07f2d5919ef02a090e9615845e9a) fix: don't enforce bundle verified
* [`cf3e56a`](https://github.com/siderolabs/image-factory/commit/cf3e56a9faf0384d2bf9878754dc437a1ba76106) chore: bump talos
* [`8723b02`](https://github.com/siderolabs/image-factory/commit/8723b0274e1b146dcba294dbd8b32686e3959654) fix: drop sbc board support
* [`f0150c4`](https://github.com/siderolabs/image-factory/commit/f0150c419ddc611401146b46ae1c7779a7358255) feat: use rootless Image Factory
* [`f57218f`](https://github.com/siderolabs/image-factory/commit/f57218fbf014441bf36d5571b488f21dcce16ce8) feat: refactor configuration of image factory
* [`e440ce7`](https://github.com/siderolabs/image-factory/commit/e440ce7a1c63f643f42c5a9ccbe9efba1bffa9c5) fix: support new cosign bundle format
* [`5eb1775`](https://github.com/siderolabs/image-factory/commit/5eb17756a1ba3b8ba9c8df72683e3bfa6aa94247) feat: introduce Enterprise Image Factory
</p>
</details>

### Changes from siderolabs/talos
<details><summary>222 commits</summary>
<p>

* [`59311a792`](https://github.com/siderolabs/talos/commit/59311a7924b908aaa2761e82e03f6fa473a4c3ee) release(v1.13.0-alpha.2): prepare release
* [`009f0d6ca`](https://github.com/siderolabs/talos/commit/009f0d6ca0cf13e5778a7c46587ac0dc9d30d5e9) chore: update pkgs
* [`ba56b0295`](https://github.com/siderolabs/talos/commit/ba56b02954fb275f8ff2ed20e38b51a75c3a8371) feat: include hid-multitouch.ko kernel module in rootfs
* [`ae29a0dcc`](https://github.com/siderolabs/talos/commit/ae29a0dcce527b90553b25230abbb5a8d4bd504c) feat: update Linux to 6.18.13
* [`7cf1de279`](https://github.com/siderolabs/talos/commit/7cf1de2794a1d4838efca378aff433fad5e1823c) fix: bring in new version of go-cmd and go-blockdevice
* [`c8800b41e`](https://github.com/siderolabs/talos/commit/c8800b41e511ce6bb4dda3e28b69c4d091177435) fix: update path handling on talosctl cgroups
* [`0a7b6eb2c`](https://github.com/siderolabs/talos/commit/0a7b6eb2c98979aa8a604f677c4dd1d54f1285e5) chore: test extensions
* [`8b1c974a2`](https://github.com/siderolabs/talos/commit/8b1c974a2a733c870f371ccb7a86ccc616dbc7ea) refactor: drop termui-widgets library
* [`5baa0028e`](https://github.com/siderolabs/talos/commit/5baa0028e65765fc0fd1179f72377bf2a2085deb) fix: add owning inventory annotation to talos manifests
* [`d3e793d14`](https://github.com/siderolabs/talos/commit/d3e793d14117891103ca4df8507124b18913a56c) fix: stop Kubernetes client from dynamically reloading the certs
* [`6a5a0e3bd`](https://github.com/siderolabs/talos/commit/6a5a0e3bd4197a4fadfcfe094876e46d4b878a0a) feat: support pattern link aliases
* [`9758bd4fe`](https://github.com/siderolabs/talos/commit/9758bd4fe0e28803acf11f3b9c9da744883aa9dc) feat: update Go to 1.26
* [`e00aed0f6`](https://github.com/siderolabs/talos/commit/e00aed0f6694bb3c8e14a0ef413ef0e62ae02981) feat: update Kubernetes v1.36.0-alpha.1
* [`f20445ad0`](https://github.com/siderolabs/talos/commit/f20445ad0981175d6444340325af5fc747993559) chore: improve logging of disk encryption handling
* [`f018fbe7b`](https://github.com/siderolabs/talos/commit/f018fbe7ba145ff86ebe0d4d09b323b9715ef1a9) fix: handle raw encryption keys with `\n` properly
* [`e5b0eb017`](https://github.com/siderolabs/talos/commit/e5b0eb017ff989e812d6444f668bf17723bb7ec4) fix: hold user volumes root mountpoint
* [`8a0e79774`](https://github.com/siderolabs/talos/commit/8a0e79774409ce7605f9cd21d769f47e5db656db) refactor: split locate and provision
* [`a59db0e92`](https://github.com/siderolabs/talos/commit/a59db0e92213296c4c9599fb0d230908caabdf30) fix: improve OpenStack bare metal network configuration reliability
* [`659009ad8`](https://github.com/siderolabs/talos/commit/659009ad875c0625ac24094dc44020b015ab8b50) fix: remove stale endpoints
* [`dab0d4783`](https://github.com/siderolabs/talos/commit/dab0d478378dfc6c2862c38633ca4494a41e7ecd) fix: allow static hosts in `/etc/hosts` without hostname
* [`45f214154`](https://github.com/siderolabs/talos/commit/45f214154cea364d86bfbba81a5ad4f272a4c8fd) feat: update go-kubernetes to use new Myers diff
* [`35ad0448c`](https://github.com/siderolabs/talos/commit/35ad0448c9ae93cd642d80ebb7d95b768ba0ab9b) fix: switch to better Myers algorithm implementation
* [`0048464be`](https://github.com/siderolabs/talos/commit/0048464be854d94fb607e38daa83e00767fe8cbc) feat: update etcd to v3.6.8
* [`5df10f260`](https://github.com/siderolabs/talos/commit/5df10f2604b537504f76b14e028f88a946aacbd7) fix: use mcopy instead of diskfs to populate VFAT
* [`ce53ffa90`](https://github.com/siderolabs/talos/commit/ce53ffa900a438f6669460a2ce9af874c1f87708) fix: disks flag parsing and handling in create qemu command
* [`3bd3dd7ca`](https://github.com/siderolabs/talos/commit/3bd3dd7ca92401312079e37584bfbf7942eab93a) fix: memory overuse in imager VFAT
* [`f118ee47e`](https://github.com/siderolabs/talos/commit/f118ee47eaba662dc161d37fae5ae8f2b3de9819) fix: read multi-doc machine config with newer talosctl
* [`70c6c2154`](https://github.com/siderolabs/talos/commit/70c6c2154e87d4a6748aebdfa2c50cbc97a0dd89) feat: add filter for KubeSpan advertised networks
* [`daf18abf4`](https://github.com/siderolabs/talos/commit/daf18abf419b21a6e70dcca0b5b83d33cfee6188) fix: fix talosctl debug in enforcing mode
* [`33b5b2565`](https://github.com/siderolabs/talos/commit/33b5b25652360a114d0b2cea412bf018cbf84df3) fix: ignore volumes in wave calculation without provisioning
* [`a16392559`](https://github.com/siderolabs/talos/commit/a16392559a488993c3e26810df57da3cae5c24c5) feat: add explicit service account support to Talos client
* [`4d531884e`](https://github.com/siderolabs/talos/commit/4d531884e9c28d480f24b61a83f140df0ffbe4b3) chore: update dependencies
* [`406b8c83c`](https://github.com/siderolabs/talos/commit/406b8c83c9b33b1917b9dd16aa1efeb2df189f0f) feat: update doc links to docs.siderolabs.com
* [`87615f551`](https://github.com/siderolabs/talos/commit/87615f551183cd322dafebf368a347d928a14442) feat: implement network policies with Flannel CNI
* [`6995bc1b1`](https://github.com/siderolabs/talos/commit/6995bc1b1ea54e1a8fd6426fef11293f35106ac7) chore: update homebrew formula on release
* [`7942d5a98`](https://github.com/siderolabs/talos/commit/7942d5a98c1d689a94e78219be09a0fc69d07b08) fix: image gc controller config
* [`52e8727d0`](https://github.com/siderolabs/talos/commit/52e8727d0112967a62a3d9ae6bf26d713db242e1) feat: add IPv6 GRE support
* [`9690dbad0`](https://github.com/siderolabs/talos/commit/9690dbad02cfc8682d697679b655e753039c5254) chore: bump tools (including linter)
* [`2628eb2ec`](https://github.com/siderolabs/talos/commit/2628eb2ece05d7f817fc42e12b979d3f8ca9710c) fix: typo with rpi_5 profile name
* [`d5ebcd7ca`](https://github.com/siderolabs/talos/commit/d5ebcd7cae1a20c8000e2f4d5a02c81e4dbe5186) fix: stop building talosctl debug on Windows
* [`8b85c7c63`](https://github.com/siderolabs/talos/commit/8b85c7c637cc08d35bbf6968abebb8c4cdfb82ad) chore: update deps
* [`d905035b5`](https://github.com/siderolabs/talos/commit/d905035b5e5c7787a5171ba2e0127c89755e8774) fix: swap volume configuration for min/max size
* [`d43a01ccb`](https://github.com/siderolabs/talos/commit/d43a01ccbdd318080b54e52d2f2fbec93042c458) feat: implement `talosctl debug`
* [`34a31c979`](https://github.com/siderolabs/talos/commit/34a31c9797d5a7e1700c3d945a21367b81c79385) feat: add mount options support for existing volumes
* [`1bf95eed1`](https://github.com/siderolabs/talos/commit/1bf95eed185152c38397cd3b43b6ff9d421678c5) feat: improve dashboard uptime display
* [`055add7ae`](https://github.com/siderolabs/talos/commit/055add7aeb158b6f4e09ef06966de7622d1b3940) release(v1.13.0-alpha.1): prepare release
* [`900516e68`](https://github.com/siderolabs/talos/commit/900516e68950e4b94696f6a9b481cefee44b3360) chore: update image signer
* [`938de566e`](https://github.com/siderolabs/talos/commit/938de566eca30af3cc4355a94931186f19b682f2) feat: bump kernel
* [`388cec727`](https://github.com/siderolabs/talos/commit/388cec72796d0ecd0c7103efcaab9066e9b62509) feat(overlays): add new overlays
* [`9f2dd6312`](https://github.com/siderolabs/talos/commit/9f2dd6312f9d49e4d03347c98b100119f94cf807) refactor: api tests
* [`a90783146`](https://github.com/siderolabs/talos/commit/a90783146fc2d475055bfce0f8b5120969f74dc7) feat: add a helper module to generate standard patches
* [`1fec5b23d`](https://github.com/siderolabs/talos/commit/1fec5b23d0c10e53863a7c0f89f862708a7f4069) fix: implement merger for PercentageSize
* [`8b245b8f2`](https://github.com/siderolabs/talos/commit/8b245b8f269b6c8cb463f2cf537d2ed2ab6924ec) feat: implement new image service APIs
* [`d90c775b8`](https://github.com/siderolabs/talos/commit/d90c775b8441705003de3427b2e6831dcbfb449f) chore: rename internal `talosctl debug air-gapped`
* [`2165280d0`](https://github.com/siderolabs/talos/commit/2165280d0eedf59899ad44e2f3289d81b3dab466) refactor: change the way one2many proxying is picked
* [`b1b703dbe`](https://github.com/siderolabs/talos/commit/b1b703dbe2b25785ded0c77f23d674d9b9934975) chore: move sync logging code to go-kubernetes package
* [`e48c6d7ab`](https://github.com/siderolabs/talos/commit/e48c6d7ab9c8a2e28ebe2115ac09f1557bbcca33) fix: allow to expose a port multiple times in Docker
* [`410d8cb57`](https://github.com/siderolabs/talos/commit/410d8cb5727ccf054c9097f33bc916d87076a599) fix: undo CRLF on Windows (talosctl edit)
* [`859d3f03c`](https://github.com/siderolabs/talos/commit/859d3f03c444d98b94a06adac3648562e3b1228b) feat: add RPi5 to the list of supported SBCs
* [`0bd48bbc6`](https://github.com/siderolabs/talos/commit/0bd48bbc6f365770167ee753be563eb4179fcadb) fix(talosctl): pass --k8s-endpoint flag to rotate-ca kubernetes rotation
* [`b9e27ebe7`](https://github.com/siderolabs/talos/commit/b9e27ebe72c4302c416fd8efb007c3966004ddd6) feat: update Linux kernel with dm-integrity
* [`6aa9b0677`](https://github.com/siderolabs/talos/commit/6aa9b0677ed7ca4955fead474e36a533b3250ad9) fix: skip empty documents on config decoding
* [`494492489`](https://github.com/siderolabs/talos/commit/494492489b29b615a8a874c0648690ed3b9adb58) fix: always set advertised peer URLs
* [`782cc507d`](https://github.com/siderolabs/talos/commit/782cc507dc33c87caa5ff985eea5f4439c3e1012) fix: open the filesystem as read-only
* [`28e61a740`](https://github.com/siderolabs/talos/commit/28e61a740a906fadfea098f38a9c9f4e8c32773e) fix: set GRUB prefix correctly on arm64
* [`a4f1c5239`](https://github.com/siderolabs/talos/commit/a4f1c5239ef7227856640c230e0d0364d9eedbd2) feat: update GRUB to 2.14
* [`562920701`](https://github.com/siderolabs/talos/commit/562920701e2999cbb6687e55de96719aba4064fd) fix: use node podCIDRs for kubespan advertiseKubernetesNetworks
* [`39460365c`](https://github.com/siderolabs/talos/commit/39460365c1726095e20cf3cc7c079c234b8022d6) feat: implement layering for ProbeSpec
* [`b5c760f70`](https://github.com/siderolabs/talos/commit/b5c760f7076570bc04be02af0ea493f95d8338d0) feat: add ProbeConfig for network connectivity probes
* [`4b274f761`](https://github.com/siderolabs/talos/commit/4b274f76159495cc6c2977ec3bbade71e35aade8) feat: support aws cert manager in imager
* [`417209512`](https://github.com/siderolabs/talos/commit/41720951251102f1c174e501a3103e55720a1d8b) fix: fallback to /proc/meminfo for memory modules
* [`7f1147bed`](https://github.com/siderolabs/talos/commit/7f1147bed495a06d336f5be1da6073921b5e52dc) fix: add warnings to 802.3ad bond
* [`ddd6b186e`](https://github.com/siderolabs/talos/commit/ddd6b186eb8f527324736576182dafbce3423da5) refactor: generate GRUB images
* [`c7aa266ea`](https://github.com/siderolabs/talos/commit/c7aa266ea5c9d3fbd465dc651f2ebfec622612e7) fix: overwrite resolver config with machine config
* [`cf70f05fa`](https://github.com/siderolabs/talos/commit/cf70f05fa40312c30d8345c2fb15ce8eda86a7a7) fix: oracle platform file format
* [`8c7b8f5b7`](https://github.com/siderolabs/talos/commit/8c7b8f5b7d6dec144f7985a7c8a8a582c38f3154) feat: add support for negative max size
* [`77bc3d21f`](https://github.com/siderolabs/talos/commit/77bc3d21fa40e188af4b5dd93e1cda289e858d56) fix: marshal of FailOverMac property
* [`38e280c93`](https://github.com/siderolabs/talos/commit/38e280c9319ef1ecb1455b3cc8b8d0d1d7426ccd) fix: make OOM expression a bit less sensitive
* [`3d1301640`](https://github.com/siderolabs/talos/commit/3d1301640d44d58303160400e4954c36f53341f9) fix: wipe the first/last 1MiB in addition to wiping by signatures
* [`1aa6528ad`](https://github.com/siderolabs/talos/commit/1aa6528adcddfb6a5ed66cc26cac1a0fcdb37516) fix: make OOM controller more precise by considering separate cgroup PSI
* [`f7072c050`](https://github.com/siderolabs/talos/commit/f7072c050e607de16781a65eb97ab2a1828b05fb) fix: check if the device is not mounted when wiping
* [`743c3b94b`](https://github.com/siderolabs/talos/commit/743c3b94b958e4abcbf70d4064f2ae0e0bbb0712) fix: use correct containerd import path
* [`f2dd08594`](https://github.com/siderolabs/talos/commit/f2dd08594e8e474c7b3891dc46c64f27c724dbc0) feat: report image pull progress in the console
* [`72fe98a06`](https://github.com/siderolabs/talos/commit/72fe98a06f31536454f201d703f8ae6a071235b5) fix: boot with GRUB
* [`d4ed13d93`](https://github.com/siderolabs/talos/commit/d4ed13d9394b087e8877eba25950f344894803a1) fix: add talos version to Hetzner Cloud client user agent
* [`150c41c30`](https://github.com/siderolabs/talos/commit/150c41c30ed3f066f10bd2bdc2afa9b2c5a97597) feat: update Linux to 6.18.5
* [`01a367891`](https://github.com/siderolabs/talos/commit/01a3678913de0fa4d309a361428c117d24ce0d1e) fix: use append instead of prepend in service-account-issuer
* [`d1954278a`](https://github.com/siderolabs/talos/commit/d1954278a1ba3470b2e5ccae90762078c18d69e9) feat: add extraArgs from service-account-issuer
* [`91b88f7f9`](https://github.com/siderolabs/talos/commit/91b88f7f994cccad15cbec1aa8019bd19b84ae91) feat: support multiple values for extraArgs
* [`96e604874`](https://github.com/siderolabs/talos/commit/96e604874b17e7aa8b62bfb25737f349e539bc5a) fix: add hostname to endpoints
* [`7033275a7`](https://github.com/siderolabs/talos/commit/7033275a7a22d51e83c9e760ba37d2ad6ab22f28) refactor: move BootloaderKind into machinery
* [`71adaf0ea`](https://github.com/siderolabs/talos/commit/71adaf0ea5b558c8a16e2acfdec3671611455985) fix: sort mirrors and tls configs when generating the machine config
* [`34f09a300`](https://github.com/siderolabs/talos/commit/34f09a3004fe1b77c16dd33b04adca95fb6876a5) feat: add VLAN support to OpenStack platform
* [`5127ef7c2`](https://github.com/siderolabs/talos/commit/5127ef7c28b360f9c7c033f77c58cef729e5278d) fix: wipe disk by signatures
* [`415bfaedb`](https://github.com/siderolabs/talos/commit/415bfaedb6ae8d42b5927fdc5b7cfe8aa781a791) fix: panic in configpatcher when the whole section is missing
* [`e5aca71cd`](https://github.com/siderolabs/talos/commit/e5aca71cd0557557e50c39d82eda2c938f627d62) fix: fix healthcheck timeout
* [`634b71e2d`](https://github.com/siderolabs/talos/commit/634b71e2d028bf13d838acad8809c95384b6eed9) docs: move talosctl pcap example to Example Block
* [`818492731`](https://github.com/siderolabs/talos/commit/8184927316c5de7d9b04f21474a60cc791c3d26d) feat: implement KubeSpan multi-document configuration
* [`4d0604b9d`](https://github.com/siderolabs/talos/commit/4d0604b9d93851f444a00dbd84fcac76d21d35c2) chore: remove unrelated machineconfig
* [`e36863470`](https://github.com/siderolabs/talos/commit/e36863470b14496c3d84417e63fef45e6060603b) feat: add it87 hwmon module
* [`308c75090`](https://github.com/siderolabs/talos/commit/308c75090774d2510c2ec08e63e179a5c0fa6987) fix: resolve SideroLink Wireguard endpoint on reconnect
* [`e4ef494de`](https://github.com/siderolabs/talos/commit/e4ef494decdf97664c4803aa3861015fce49760e) fix: drop the persist config flag from gen config
* [`c3176adcf`](https://github.com/siderolabs/talos/commit/c3176adcf981811a326c971c81c4b591f54e116a) feat: add EnvironmentConfig document
* [`c839b3880`](https://github.com/siderolabs/talos/commit/c839b38809b3a0029061d43477555ec31e283aa5) feat: expose more SSA options in the upgrade-k8s command
* [`b8ff9677e`](https://github.com/siderolabs/talos/commit/b8ff9677e4f9a64908ae00bb1d80aa2442a00a60) fix: handle correctly incomplete RegistryTLSConfig
* [`99f2ddada`](https://github.com/siderolabs/talos/commit/99f2ddada895011036af1435dd10bac3be0a9171) fix: bond config via platform
* [`2449ffea4`](https://github.com/siderolabs/talos/commit/2449ffea45304459ea8895b535b6f070a9249172) fix: allow HostnameConfig to be used with incomplete machine config
* [`35fc52087`](https://github.com/siderolabs/talos/commit/35fc5208728dbc3e0b139aff4c06f25208445637) fix: lock down etcd listen address to IPv4 localhost
* [`27253d731`](https://github.com/siderolabs/talos/commit/27253d7317a473cbbc0f5c0eee634173bdd2eda7) feat: use new xfs config file
* [`c9d84ae21`](https://github.com/siderolabs/talos/commit/c9d84ae21e203529a6952c165ff04d602a2a6ad6) fix: generate OCI-compliant image config
* [`7a4b2b33a`](https://github.com/siderolabs/talos/commit/7a4b2b33abe8a3011f37f0a8f4848dd846d0396f) fix: update VIP config example
* [`080efcbda`](https://github.com/siderolabs/talos/commit/080efcbda2c4334f9d8c70804a5a37f0cdb2df2d) feat: add k8s-version parameter to k8s-bundle
* [`b764f5f72`](https://github.com/siderolabs/talos/commit/b764f5f724bf8af3acaac74942ea91a86e593322) fix: skip sync test when kube-proxy is disabled
* [`70e67787d`](https://github.com/siderolabs/talos/commit/70e67787d6d34d93a34871b2d25d64f6a7575d76) feat: imager: populate filesystems with root owned files
* [`7416dca59`](https://github.com/siderolabs/talos/commit/7416dca59378dc282e42ea30107cf40326cc593c) fix: print talosctl images to release notes
* [`dc2009e47`](https://github.com/siderolabs/talos/commit/dc2009e4779684a6a4252d4dfd2aa02d1b60c2da) chore: use context when creating filesystems
* [`85f7be6e3`](https://github.com/siderolabs/talos/commit/85f7be6e3f14bf160cf32bccf7418b31968d474f) chore: update slack links
* [`154952175`](https://github.com/siderolabs/talos/commit/154952175ab73ac65722732b146a0ee1c56b2f4d) fix: disable swap for system services
* [`d98b415af`](https://github.com/siderolabs/talos/commit/d98b415afea7b1820153151c0273df24a101742e) fix: drop more non-overlay SBC stuff
* [`226cd6bc1`](https://github.com/siderolabs/talos/commit/226cd6bc1d70662cb7f7736ac6fad117170a36fb) fix: do not allocate for the actual disk image file
* [`53f5bf8d2`](https://github.com/siderolabs/talos/commit/53f5bf8d2c97e91bee06bcb5948170015486ea77) fix: overlay installers
* [`10d0cfd93`](https://github.com/siderolabs/talos/commit/10d0cfd93a083fb8b71b7c0297df52feb55e044b) fix: overlay install in image mode
* [`77086694d`](https://github.com/siderolabs/talos/commit/77086694d18b69802e542156fc12cd7cf066efc2) fix: partition data population
* [`4d5657b1a`](https://github.com/siderolabs/talos/commit/4d5657b1a34c939b63b2cc3ee11ed45ad1bf23c3) fix: drop SBC board code
* [`c4f3f6d3e`](https://github.com/siderolabs/talos/commit/c4f3f6d3e59b58016ba8546c5bd3e8e465fbbf52) feat: implement kubernetes server-side apply
* [`f12fd2b0a`](https://github.com/siderolabs/talos/commit/f12fd2b0a9fdf8f53ec5714d3ad18b695973e0b0) test: bump Image Factory tests
* [`c76484e58`](https://github.com/siderolabs/talos/commit/c76484e5879a7e48197e442cf22044d3d0363846) release(v1.13.0-alpha.0): prepare release
* [`f0d8a6851`](https://github.com/siderolabs/talos/commit/f0d8a685173354e5fd148786872062a342c4282a) test: skip the source bundle on exact tag
* [`c57701d65`](https://github.com/siderolabs/talos/commit/c57701d6590388e7d6418af67e8237c7d60ccf54) fix: remove interactive installer
* [`43937c1cd`](https://github.com/siderolabs/talos/commit/43937c1cd42758a15026261fe8f0e06daaebdcbd) feat: update Linux and systemd
* [`72a194df8`](https://github.com/siderolabs/talos/commit/72a194df88f2800cee3372241fbad419b07f7bbf) feat: add VM CPU hot-add rules
* [`f09ae1e0d`](https://github.com/siderolabs/talos/commit/f09ae1e0d2e1b7842d504b594b71a325af7733e5) fix: probe small images correctly
* [`8f2b33799`](https://github.com/siderolabs/talos/commit/8f2b337994fdeff76a0ae9e1730b4b9f596ff1bb) feat: imager support rootless builds
* [`c7525a97e`](https://github.com/siderolabs/talos/commit/c7525a97ef8615e903be183d7938b6d2a3b89464) feat: support creating filesystems from folder
* [`e2bffb5ce`](https://github.com/siderolabs/talos/commit/e2bffb5cebaaf28f9dfff24f41ecbb2809fc60e5) chore: refactor imager code so it's more clear
* [`0fb50dbd0`](https://github.com/siderolabs/talos/commit/0fb50dbd0a5b7b80187e50d501cec4b3fe434dc2) fix: invalid versions check in talos-bundle
* [`b5dd56032`](https://github.com/siderolabs/talos/commit/b5dd5603207a46d8eed240173f06aeffd6a9c0e7) test: upgrade versions in upgrade tests
* [`3dfa4d6e4`](https://github.com/siderolabs/talos/commit/3dfa4d6e40dcae2db47e89443568be3ae48b3ae1) fix: make upgrade work with SELinux enforcing=1
* [`786c8e2ee`](https://github.com/siderolabs/talos/commit/786c8e2ee757c2d7b30d5bded954e584af3a058e) feat: ship pigz/igzip in rootfs to speed up image decompression
* [`48d242918`](https://github.com/siderolabs/talos/commit/48d242918bc97e6a01434bee6fcdcfa735fd1f5a) feat: update containerd to 2.2.1
* [`536541afe`](https://github.com/siderolabs/talos/commit/536541afe497d5f61cfcd0c01cf580ab5b3be164) fix: mount volume mount/unmount race
* [`39117d457`](https://github.com/siderolabs/talos/commit/39117d45766b139ed6a0c1290f757e4b26d31d92) feat: update dependencies
* [`f0f420725`](https://github.com/siderolabs/talos/commit/f0f420725c6a4f628cdc1b80d59713c375beb9b7) fix: bond setting change detection
* [`8d6a7a867`](https://github.com/siderolabs/talos/commit/8d6a7a8677a5d1d61432fa94ca030351fd9852f2) feat: update Kubernetes to 1.35.0
* [`845a0d09c`](https://github.com/siderolabs/talos/commit/845a0d09cd770a15db762ddda4d3d27f58656cfe) feat: update etcd 3.6.7, CoreDNS 1.13.2
* [`b95912e04`](https://github.com/siderolabs/talos/commit/b95912e04907b78bd06987c6d3948f8f1804d844) feat: enforce `proc_mem.force_override=never` by default
* [`681f3e84c`](https://github.com/siderolabs/talos/commit/681f3e84c85677f49ddbcd4a47e325d4a85af692) test: run virtiofs tests only when virtiofsd is running
* [`0592ff0cd`](https://github.com/siderolabs/talos/commit/0592ff0cdbf54475dc91bfb7c9b9c3047bbe13da) fix: drop the Omni API URL check on IP address
* [`a4879a5fa`](https://github.com/siderolabs/talos/commit/a4879a5fa2ded9b7b52ff7506b5493ae12939bba) feat: update Linux to 6.18.1
* [`43b43ff18`](https://github.com/siderolabs/talos/commit/43b43ff189b7e5f37eaa75f4926c26ee21ffa5cb) docs: split talosctl commands into groups
* [`6d17c18bf`](https://github.com/siderolabs/talos/commit/6d17c18bf908d3cd69ff920d0cff67b653a385f3) feat: enable Powercap and Intel RAPL
* [`884e76662`](https://github.com/siderolabs/talos/commit/884e76662af34448d9904372f1256f59ce161f99) docs: fix the talosctl cluster create help output
* [`6dc31be4f`](https://github.com/siderolabs/talos/commit/6dc31be4f982f62ba4aeb1b3b4e65ce022447eb4) fix: exclude new Virtual IPs configured with new config
* [`94905c73e`](https://github.com/siderolabs/talos/commit/94905c73e93fd7dac38d911dc4264e4d0fe0081d) feat(talosctl): support running qemu x86 on Mac
* [`f871ab241`](https://github.com/siderolabs/talos/commit/f871ab241c0f034401fbf61e32e7201cced49441) fix: provide json support in `nft` binary
* [`694f45413`](https://github.com/siderolabs/talos/commit/694f45413fec8cc4f58a79e76034bd4bcec2bbdf) feat: external volumes
* [`39feb16d2`](https://github.com/siderolabs/talos/commit/39feb16d2ed3bcb65d66483c0729bcec29f7b93e) fix: update containerd 2.2.0 with cgroups patch
* [`82027eb9b`](https://github.com/siderolabs/talos/commit/82027eb9b30aa128099b27f638098d78857ecb4b) fix: bond configuration with new settings
* [`121b13b8f`](https://github.com/siderolabs/talos/commit/121b13b8f8d6e5a487971f727c6e028c7ffa20f3) fix: disable kexec on arm64
* [`7eaa725d0`](https://github.com/siderolabs/talos/commit/7eaa725d0dba18392279f5b43d167aaf18f43b99) fix: selection of boot entry
* [`949bdb90a`](https://github.com/siderolabs/talos/commit/949bdb90ab2fd711c47583d96bd29a1ca90bbf41) feat: add Secure Boot to CloudStack platform config
* [`798143a88`](https://github.com/siderolabs/talos/commit/798143a886e4055e764a9ad17cefe8ad4db0572e) fix: discard better klog message from Kubernetes client
* [`008cd0986`](https://github.com/siderolabs/talos/commit/008cd0986cbbbd5527d91c01b951e311ba014b97) fix: disable kexec in talosctl cluster create on arm64
* [`bb62b29ed`](https://github.com/siderolabs/talos/commit/bb62b29edb2fb704846ceeed2019f0ebaced30be) chore: prepare talos for 1.13
* [`c0935030a`](https://github.com/siderolabs/talos/commit/c0935030ac3d966149591a3aaa8e430da768d678) chore: fork reference docs for 1.13.x
* [`e387e48b3`](https://github.com/siderolabs/talos/commit/e387e48b30b3a3b991f1f611099f48fddefa851b) fix: do not override DNS on MacOS
* [`1e7e87fb1`](https://github.com/siderolabs/talos/commit/1e7e87fb192521937b581ecd94a0aa0c861f2a5f) fix: rework NFT rules for KubeSpan
* [`51bcfb567`](https://github.com/siderolabs/talos/commit/51bcfb567915d2b27e4b5321e080220bc618086b) feat: rename image default and source bundle
* [`585abe944`](https://github.com/siderolabs/talos/commit/585abe94431f06b3ebf4b6a64ad1b5918708f866) feat: update Kubernetes to v1.35.0-rc.1
* [`f301e3e9b`](https://github.com/siderolabs/talos/commit/f301e3e9ba47d5f46f1990a9bd21fd4e671c38f3) fix: update KubeSpan MSS clamping
* [`74c1df6f4`](https://github.com/siderolabs/talos/commit/74c1df6f4b2ac8d989d1e42d6c7c0016411638ee) test: propagate MTU size to QEMU in `talosctl cluster create`
* [`d347ca1af`](https://github.com/siderolabs/talos/commit/d347ca1af162c8d948899d58fc3f76dd0a94f138) fix: update CNI plugins to 1.9.0
* [`e3f8196b4`](https://github.com/siderolabs/talos/commit/e3f8196b4c767ca68df9f6c85ed25c7e12fb4d87) chore: update Grype and Syft
* [`e1b8ab323`](https://github.com/siderolabs/talos/commit/e1b8ab3236e956bc4b37e227423aea0f97612a5c) docs: add misssing period
* [`cd04c3dde`](https://github.com/siderolabs/talos/commit/cd04c3dde70f604603fd7996c62adf5a17cfbd41) docs: update release notes
* [`fc8ae3249`](https://github.com/siderolabs/talos/commit/fc8ae3249fac82cbdb5521ca8797a8451bdaa9fd) docs: add omni join token example to create qemu command
* [`9fa00773c`](https://github.com/siderolabs/talos/commit/9fa00773caf2d092d953ff58d04cf94803039b94) chore: update go-blockdevice
* [`ba13b6786`](https://github.com/siderolabs/talos/commit/ba13b678654e2896e1a99b1af8b51a9239b0a559) fix: correct condition to use UKI cmdline in GRUB
* [`d2ce3f47f`](https://github.com/siderolabs/talos/commit/d2ce3f47f8515231f27983abaaf269a059e2e90d) docs: drop machine.network example
* [`cf087c1e0`](https://github.com/siderolabs/talos/commit/cf087c1e01bc1226049a57186f48b2e6b5739c5c) test: bird2 extension
* [`13df94388`](https://github.com/siderolabs/talos/commit/13df943884a59bd1d42721ba42bcb36349d40624) fix: adapt SELinuxSuite.TestNoPtrace to new strace version
* [`861787c38`](https://github.com/siderolabs/talos/commit/861787c380bff3ba2fa29f49837bc173a2719578) fix: mark secureboot as supported for metal
* [`04e3e87ad`](https://github.com/siderolabs/talos/commit/04e3e87adcbd24ee0d82dce4cc27121d34d316f4) fix: clean up kubelet mounts
* [`21057903a`](https://github.com/siderolabs/talos/commit/21057903a2ca01d88cc5f97c084567d1981f73c5) fix: clear provisioning data on SideroLink config change
* [`0f9f4c05f`](https://github.com/siderolabs/talos/commit/0f9f4c05ffad9413e1f1533c68eae38dc91c9716) feat: update Kubernetes to 1.35.0-rc.0
* [`d4309d7b1`](https://github.com/siderolabs/talos/commit/d4309d7b1aec9d2852173fd704b09dfabe2cf217) fix: add a timeout for DNS resolving for NTP
* [`dd6c1089c`](https://github.com/siderolabs/talos/commit/dd6c1089c8f30d815c80ab10544a0fef27ddd14c) feat: update Linux to 6.18.0
* [`e9a30bf9a`](https://github.com/siderolabs/talos/commit/e9a30bf9a8ee55ab9ae5d9c9a18362434b0202ad) test: revert add direct connectivity CA rotation test
* [`cc95562bc`](https://github.com/siderolabs/talos/commit/cc95562bc830496986a395cdde352d48d4a1d146) fix: don't disable LACP by default
* [`c9fe4679b`](https://github.com/siderolabs/talos/commit/c9fe4679bf9c1dcdf175b95a02f1eaacab4ff085) test: add platform acquire/not valid config unit-test
* [`5a03a7a20`](https://github.com/siderolabs/talos/commit/5a03a7a20acffa8eedf40524f8d070e37e41f24e) chore: fix longhorn test
* [`a0cfc3527`](https://github.com/siderolabs/talos/commit/a0cfc3527481c4784edf87c3d7823b10a21d1e4d) feat: implement logs persistence
* [`51b732bea`](https://github.com/siderolabs/talos/commit/51b732beabc9948e58f9aa4d81b79afb9bd61243) fix: selection of boot entry
* [`18f8ac369`](https://github.com/siderolabs/talos/commit/18f8ac369ba52f2640508134d3983f006f698129) feat: update Kubernetes to 1.35.0-beta.0
* [`92fa7c5e4`](https://github.com/siderolabs/talos/commit/92fa7c5e43da96a492003a2c9184cf818fbbb9f0) chore: update pkgs for NVIDIA 580.105.08
* [`f489299b6`](https://github.com/siderolabs/talos/commit/f489299b603a2aff0f292fa941ae8925fdda3492) chore: correct condition for running k8s integration tests
* [`ab149750d`](https://github.com/siderolabs/talos/commit/ab149750d475ef059debfc3730e9e0a32ad6e601) chore: update tools/pkgs to 1.13.0-alpha.0
* [`87ff9f860`](https://github.com/siderolabs/talos/commit/87ff9f8606e04fe99e23261418a762372647b077) test: fix the image-factory test to pass IF endpoint
* [`2ffe538e7`](https://github.com/siderolabs/talos/commit/2ffe538e7307f0ac3dbac2eba4b36ea98162ec78) test: add direct connectivity CA rotation test
* [`70f6b80e0`](https://github.com/siderolabs/talos/commit/70f6b80e03acd507580211724cc51b7867bf8a76) chore(ci): skip multipath extension tests
* [`561cfb60c`](https://github.com/siderolabs/talos/commit/561cfb60c313a9bdc70ed2ff2729549bc8c50fcb) chore: update pkgs and tools version
* [`2f42202a7`](https://github.com/siderolabs/talos/commit/2f42202a7ccee0e33e43b2081929b5510db5d713) fix: simplify OOM expression
* [`7b06ae8c2`](https://github.com/siderolabs/talos/commit/7b06ae8c2cf1069cb77cddee0986afc5af837bcc) test: fix flaky LinkSpec/Wireguard test
* [`e715f3871`](https://github.com/siderolabs/talos/commit/e715f387137fa566a4824c051b624e013a93c49f) feat: present kernel log as `talosctl logs kernel`
* [`e2ee39b8a`](https://github.com/siderolabs/talos/commit/e2ee39b8ac54ada49dd0a7ffaab4b0ae5d684792) fix: support specifying patch file without '@' symbol
* [`e202b1f9e`](https://github.com/siderolabs/talos/commit/e202b1f9e82823aa5b31625024bce65bcc53b29f) fix: trim trailing dots from certificate SANs
* [`7f7079f9c`](https://github.com/siderolabs/talos/commit/7f7079f9c0fbb30ce781aa1223d7df1a175a6206) fix: assign value of multicast setting properly
* [`eba96141e`](https://github.com/siderolabs/talos/commit/eba96141e0afc147af9a8f1969e207501232b1de) feat: update etcd to 3.6.6
* [`9945ceef3`](https://github.com/siderolabs/talos/commit/9945ceef37b13bc6e93637dcf395a8c9019e60ed) docs: add API Server Cipher Suites changelog
* [`9ed488d09`](https://github.com/siderolabs/talos/commit/9ed488d09648c09a9a5c1ed6a5cd245b84cd415d) feat: update TLS cipher suites for API server
* [`f1c04e4d6`](https://github.com/siderolabs/talos/commit/f1c04e4d6af14243a328d22bf810f27b13d83898) feat: generate mirrors patch
* [`a89108995`](https://github.com/siderolabs/talos/commit/a89108995ff13fbbef0bf5cbf429cede5ff81078) fix: add CA subject to generated certificate
* [`35dd612a5`](https://github.com/siderolabs/talos/commit/35dd612a5e59d8781e147fc36eb14f3e8bc66811) fix: add more resilient move
* [`83675838f`](https://github.com/siderolabs/talos/commit/83675838f3655b44cbd850fd82b4d17acfb00c33) feat: extend flags of cache-cert-gen
* [`80ab7a064`](https://github.com/siderolabs/talos/commit/80ab7a0643fc8057283a8ba3eb912d0ee453c143) chore: remove spammy 'clean up unused volumes' logs
* [`74d35900a`](https://github.com/siderolabs/talos/commit/74d35900af0f6451426b70eec3b6db4b72eb993c) chore: disable k8s integration tests for 1GiB worker nodes
* [`4f6218674`](https://github.com/siderolabs/talos/commit/4f621867407ec8f568f67833172ebaf2ff400346) feat: support TALOS_HOME env var
* [`0c59b3ea3`](https://github.com/siderolabs/talos/commit/0c59b3ea3f6bc49cef409a1456b4ffa3bf1d28df) feat: add multicast to linkconfig
* [`6db06f4d5`](https://github.com/siderolabs/talos/commit/6db06f4d5d51abd9e80ead6e4417f0f68856c569) feat: implement multicast setting
* [`eeded98f5`](https://github.com/siderolabs/talos/commit/eeded98f527a230c65cb041a29fefc5f693d9879) fix: add riscv64 talosctl to release artifacts
* [`a6bbae91b`](https://github.com/siderolabs/talos/commit/a6bbae91bad56328851fa91e01c17b8af7340b3c) fix: fix typos across the project
* [`83f2bdb9c`](https://github.com/siderolabs/talos/commit/83f2bdb9ce6c9466716a6ac9c94dc2222e569ee8) feat: support relative voume size
</p>
</details>

### Changes from siderolabs/talos-metal-agent
<details><summary>4 commits</summary>
<p>

* [`3bcd6af`](https://github.com/siderolabs/talos-metal-agent/commit/3bcd6afe20451d4bc99615b5a0a38d7c7ec69869) release(v0.1.4): prepare release
* [`f2f51f9`](https://github.com/siderolabs/talos-metal-agent/commit/f2f51f98b903202a98236ecaaeb3337ef7a57f0f) fix: default to IPMI port 623 when it is unsupported
* [`b475ccc`](https://github.com/siderolabs/talos-metal-agent/commit/b475ccc8d14ce28e03323f9b5ab8eafb581c3207) chore: bump deps, rekres
* [`8e92d6e`](https://github.com/siderolabs/talos-metal-agent/commit/8e92d6eeedd1cefb8e0473f1051d274807df2292) chore: bump extensions ref in boot assets image
</p>
</details>

### Dependency Changes

* **github.com/cosi-project/runtime**            v1.13.0 -> v1.14.0
* **github.com/insomniacslk/dhcp**               175e84fbb167 -> 5adc3eb26f91
* **github.com/klauspost/compress**              v1.18.3 -> v1.18.4
* **github.com/pin/tftp/v3**                     17016b3c2849 -> v3.2.0
* **github.com/siderolabs/image-factory**        v0.9.0 -> v1.0.3
* **github.com/siderolabs/omni/client**          v1.4.7 -> v1.5.8
* **github.com/siderolabs/talos**                v1.12.2 -> v1.13.0-alpha.2
* **github.com/siderolabs/talos-metal-agent**    v0.1.3 -> v0.1.4
* **github.com/siderolabs/talos/pkg/machinery**  v1.13.0-alpha.0 -> 58e006461d30
* **github.com/stmcginnis/gofish**               v0.20.0 -> v0.21.4
* **golang.org/x/net**                           v0.49.0 -> v0.51.0
* **google.golang.org/grpc**                     v1.78.0 -> v1.79.1
* **google.golang.org/protobuf**                 v1.36.11 -> f2248ac996af

Previous release can be found at [v0.8.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.8.0)

## [omni-infra-provider-bare-metal 0.8.0](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.8.0) (2026-02-05)

Welcome to the v0.8.0 release of omni-infra-provider-bare-metal!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/omni-infra-provider-bare-metal/issues.

### Contributors

* Andrey Smirnov
* Noel Georgi
* Mateusz Urbanek
* Amarachi Iheanacho
* Dmitrii Sharshakov
* Orzelius
* Laura Brehm
* Oguz Kilcan
* Justin Garrison
* Utku Ozdemir
* Bryan Lee
* George Gaál
* 459below
* Adrian L Lange
* Aleksandr Gamzin
* Alp Celik
* Andrew Longwill
* Artem Chernyshev
* Chris Sanders
* Christopher Puschmann
* Dmitry
* Edward Sammut Alessi
* Febrian
* Florian Grignon
* Giau. Tran Minh
* Grzegorz Rozniecki
* Jonas Lammler
* Lennard Klein
* Markus Freitag
* Max Makarov
* Michael Smith
* Mike Beaumont
* Misha Aksenov
* MrMrRubic
* Olivier Doucet
* Pranav
* Serge Logvinov
* Skye Soss
* Skyler Mäntysaari
* SuitDeer
* Tom
* aurh1l
* frozenprocess
* frozensprocess
* kassad
* leppeK
* samoreno
* theschles
* winnie

### Changes
<details><summary>1 commit</summary>
<p>

* [`3cf79cc`](https://github.com/siderolabs/omni-infra-provider-bare-metal/commit/3cf79ccd07a740e03d452489b8c81bed82aebe80) test: set required Omni SQLite storage path flag to integration test
</p>
</details>

### Changes from siderolabs/image-factory
<details><summary>16 commits</summary>
<p>

* [`fa266e0`](https://github.com/siderolabs/image-factory/commit/fa266e0b201a1e7f564dafb31692dda905ddb319) release(v0.9.0): prepare release
* [`6799661`](https://github.com/siderolabs/image-factory/commit/67996611c90872bbea58ea3298d3dc33994791a1) feat: show booter command in final wizard
* [`fb22bce`](https://github.com/siderolabs/image-factory/commit/fb22bcea42c92cbee1a7fe8e67c39e63b5081b57) feat: support selecting bootloader
* [`e881e4b`](https://github.com/siderolabs/image-factory/commit/e881e4b03141bff1999848e7f43a3c8d285bf049) feat: bump deps
* [`d1bec57`](https://github.com/siderolabs/image-factory/commit/d1bec579736f08a79e335bddad055ff620aa22f1) feat: implement schematic GET API
* [`f1dad9d`](https://github.com/siderolabs/image-factory/commit/f1dad9da10024c2c2a5bf529f4f9d1e9e06b0dc6) feat: better test matrix
* [`bc4f959`](https://github.com/siderolabs/image-factory/commit/bc4f9590b2ab6c7241549bee5babb2b6b721fad1) fix: remove secureboot talosctl preset
* [`db5e4dc`](https://github.com/siderolabs/image-factory/commit/db5e4dc3508b4d9e6d0f1e68e93b8c5bba607b8f) feat: add a prompt about using `talosctl cluster create qemu`
* [`2c5037c`](https://github.com/siderolabs/image-factory/commit/2c5037cf1db80a42289f1d96c9737271bad7f9a3) chore: bump deps
* [`1559666`](https://github.com/siderolabs/image-factory/commit/15596662c79c0d9a1f0cc8d06951bf74d2457390) feat: replace hardcoded artifact image constants with CLI-configurable values
* [`c27ee27`](https://github.com/siderolabs/image-factory/commit/c27ee27755d55ad5161eb2e26ed462fbe1c5d4c0) fix: return 400 when an invalid image name is requested
* [`58125d4`](https://github.com/siderolabs/image-factory/commit/58125d4d3574753d8478e7878363639cb588d8a9) feat: support proxying external installer registry
* [`d782950`](https://github.com/siderolabs/image-factory/commit/d782950320a676204c36c2a9992ab7e76ff4215e) feat: support serving TLS froom Image Factory
* [`743fe7f`](https://github.com/siderolabs/image-factory/commit/743fe7f7404defa7a1019b0dd491716c146be053) feat: support disable cosign signature verification
* [`3a20123`](https://github.com/siderolabs/image-factory/commit/3a20123740181e744c2be808c1398720abab2c4c) chore: rekres with parallel jobs
* [`241963f`](https://github.com/siderolabs/image-factory/commit/241963fbf19a47479a1b29d42bc9fa513f5f1728) chore(ci): use runner groups
</p>
</details>

### Changes from siderolabs/talos
<details><summary>388 commits</summary>
<p>

* [`54e5b438d`](https://github.com/siderolabs/talos/commit/54e5b438d8dcf6395e6424808d1155d02abf3bc0) release(v1.12.2): prepare release
* [`30da0bc19`](https://github.com/siderolabs/talos/commit/30da0bc19eb699dabf966cce38ef4477add193d4) fix: oracle platform file format
* [`7ddb37b1f`](https://github.com/siderolabs/talos/commit/7ddb37b1f3e2abf6c3406d35be92093fe4512eff) fix: make OOM expression a bit less sensitive
* [`e438ec23e`](https://github.com/siderolabs/talos/commit/e438ec23eefef97bbaa160dd6bb133b48a267ac7) fix: marshal of FailOverMac property
* [`717ed7265`](https://github.com/siderolabs/talos/commit/717ed726569d1270e2fb48df60e5fd7f43d1885b) fix: check if the device is not mounted when wiping
* [`c95c9fd06`](https://github.com/siderolabs/talos/commit/c95c9fd06508f02a770100f87da754a6fd3b9fa8) fix: wipe the first/last 1MiB in addition to wiping by signatures
* [`52bed358d`](https://github.com/siderolabs/talos/commit/52bed358d3606d04e6b4acded5dfe26cdb5f0ec9) fix: add talos version to Hetzner Cloud client user agent
* [`0e447a431`](https://github.com/siderolabs/talos/commit/0e447a4318ff2b7a398a719144690b22dce1e3f7) fix: make OOM controller more precise by considering separate cgroup PSI
* [`3b974b99e`](https://github.com/siderolabs/talos/commit/3b974b99e583c3a5bdd80e239517ef1ebc19de9c) fix: sort mirrors and tls configs when generating the machine config
* [`8b16fe50b`](https://github.com/siderolabs/talos/commit/8b16fe50bb44c7cb4bd3f50580a3ea18cdc3a727) feat: add VLAN support to OpenStack platform
* [`eb8480c4c`](https://github.com/siderolabs/talos/commit/eb8480c4ce088bd9fe705302c7e588aa01da207b) fix: panic in configpatcher when the whole section is missing
* [`4d44306dd`](https://github.com/siderolabs/talos/commit/4d44306dd148c872803578dc3880bbab307612b9) fix: wipe disk by signatures
* [`cca4cd269`](https://github.com/siderolabs/talos/commit/cca4cd269b0a4ac24627d195fad4bd9fa00c3f85) feat: add it87 hwmon module
* [`d9480eef2`](https://github.com/siderolabs/talos/commit/d9480eef2ed45b35d5f1782b651c1499451536c5) fix: resolve SideroLink Wireguard endpoint on reconnect
* [`e16c2d5bb`](https://github.com/siderolabs/talos/commit/e16c2d5bba1b6dce241905dc9e4846d45a774f78) fix: handle correctly incomplete RegistryTLSConfig
* [`dedd273df`](https://github.com/siderolabs/talos/commit/dedd273dfcd5d721e63cbe0124623ce2b5e50df4) fix: bond config via platform
* [`f527cff23`](https://github.com/siderolabs/talos/commit/f527cff239cf246891ef6e053d0aec5ce8900e22) fix: allow HostnameConfig to be used with incomplete machine config
* [`10918136c`](https://github.com/siderolabs/talos/commit/10918136c6338506d08dd86b57d82b880ea50348) fix: lock down etcd listen address to IPv4 localhost
* [`9f8d938db`](https://github.com/siderolabs/talos/commit/9f8d938db68f4c872ccf65573339e4761b4a09d4) fix: print talosctl images to release notes
* [`95433c167`](https://github.com/siderolabs/talos/commit/95433c167493a7650513379866e544bdb0adbc2e) fix: update VIP config example
* [`919394fee`](https://github.com/siderolabs/talos/commit/919394fee8122bd583ac1f0cfc55d8a0d3e3d3cb) feat: update Go to 1.25.6
* [`7ea2ef7cf`](https://github.com/siderolabs/talos/commit/7ea2ef7cf4d0d48ac9b30eca9b7ec17aa83fde50) release(v1.12.1): prepare release
* [`78a785604`](https://github.com/siderolabs/talos/commit/78a785604ad40eb9f1634c9db5477bd6ce99428c) chore: run rekres and update dependencies
* [`c31067173`](https://github.com/siderolabs/talos/commit/c3106717392a34fcca959b414f5064d6c799eaa3) fix: disable swap for system services
* [`a7e8426cf`](https://github.com/siderolabs/talos/commit/a7e8426cfb46f4c46476243032e2f4ade1fe9dfc) test: skip the source bundle on exact tag
* [`943984167`](https://github.com/siderolabs/talos/commit/943984167c22af0853d2c956677a241acece807f) fix: probe small images correctly
* [`42df71637`](https://github.com/siderolabs/talos/commit/42df71637763b1bf10bdf0fe89f650c367605b8c) fix: invalid versions check in talos-bundle
* [`a3e90e445`](https://github.com/siderolabs/talos/commit/a3e90e445f0f99b050eb98fcd9565b2b5e3397bf) fix: make upgrade work with SELinux enforcing=1
* [`ac91ade2c`](https://github.com/siderolabs/talos/commit/ac91ade2c7e435e63ed2546244d428a81abd22ad) release(v1.12.0): prepare release
* [`82553b2a1`](https://github.com/siderolabs/talos/commit/82553b2a1a713836f496b822e86e5e6788c5ebd1) fix: mount volume mount/unmount race
* [`33f6e22ec`](https://github.com/siderolabs/talos/commit/33f6e22ecb3b393d1488730c67d6f973a46b0b39) fix: bond setting change detection
* [`d5be50ac5`](https://github.com/siderolabs/talos/commit/d5be50ac55cac1c1c1deff4971fd991f364696a1) docs: split talosctl commands into groups
* [`70d3ab9ac`](https://github.com/siderolabs/talos/commit/70d3ab9ac090095c2fc8cbbfaa9c5c472d76c794) feat: update Kubernetes to 1.35.0
* [`101814d88`](https://github.com/siderolabs/talos/commit/101814d889924afe7c049106c638a32ae107a139) feat: update etcd 3.6.7, CoreDNS 1.13.2
* [`ce286825a`](https://github.com/siderolabs/talos/commit/ce286825a7f969f847ea7ad17bd2a31fa85d301c) fix: drop the Omni API URL check on IP address
* [`96f724adc`](https://github.com/siderolabs/talos/commit/96f724adccbc6fac844f9a341e36eede331b3947) feat: enable Powercap and Intel RAPL
* [`e195427c1`](https://github.com/siderolabs/talos/commit/e195427c17a004b5bcaa6f1870ce6c855ae61f1d) docs: fix the talosctl cluster create help output
* [`e025355b7`](https://github.com/siderolabs/talos/commit/e025355b759bb110925631f5f84230e99b9069df) feat(talosctl): support running qemu x86 on Mac
* [`21a914a1d`](https://github.com/siderolabs/talos/commit/21a914a1d1ca48d6bb4d47ddc8be0d0fdf74800d) fix: exclude new Virtual IPs configured with new config
* [`ca645777d`](https://github.com/siderolabs/talos/commit/ca645777dae5ad07501501dafc4717e7383045b0) fix: provide json support in `nft` binary
* [`6dd0558a3`](https://github.com/siderolabs/talos/commit/6dd0558a314af9a0dfda77b4f58a7115ef86b6fc) feat: sync pkgs
* [`c931847cc`](https://github.com/siderolabs/talos/commit/c931847ccaadf84f84e5f2befadaffb55740b592) feat: update containerd to v2.1.6
* [`a2a77004d`](https://github.com/siderolabs/talos/commit/a2a77004deac3efe6ac14f906a8bd0a3b0f926ca) release(v1.12.0-rc.1): prepare release
* [`47198780b`](https://github.com/siderolabs/talos/commit/47198780bfc084347b9ae675aaeb27a1c1d58d38) fix: bond configuration with new settings
* [`03a424bdf`](https://github.com/siderolabs/talos/commit/03a424bdf1b8a270dd694fc2738d81a3261d80cf) fix: disable kexec on arm64
* [`688fb789b`](https://github.com/siderolabs/talos/commit/688fb789beb979544e16447e512419629ea61b21) feat: add Secure Boot to CloudStack platform config
* [`66e67fd13`](https://github.com/siderolabs/talos/commit/66e67fd1394946b3425543a1aac52d4a8338e375) fix: discard better klog message from Kubernetes client
* [`d8403498c`](https://github.com/siderolabs/talos/commit/d8403498c92e0f9c37b04ad6786b2c84df5e7c95) fix: disable kexec in talosctl cluster create on arm64
* [`5ced4258c`](https://github.com/siderolabs/talos/commit/5ced4258c18f5649590a50c2927ab8e16db298ec) fix: do not override DNS on MacOS
* [`fabf3f0e7`](https://github.com/siderolabs/talos/commit/fabf3f0e73918b650b33ef0f009cacb9a15ecbc0) fix: selection of boot entry
* [`93cec4b9d`](https://github.com/siderolabs/talos/commit/93cec4b9dfdef0566152ef80c28439a7dbb0c320) fix: update CNI plugins to 1.9.0
* [`964098d96`](https://github.com/siderolabs/talos/commit/964098d9696a804de5d27284cd79dccffa7c81b9) fix: update KubeSpan MSS clamping
* [`bce04084d`](https://github.com/siderolabs/talos/commit/bce04084d6f5a9c703c7d63d1558d7d43c54dfbf) feat: rename image default and source bundle
* [`d1abc0f84`](https://github.com/siderolabs/talos/commit/d1abc0f8473c1a562e37a712624f803ce0f60fec) chore: update pkgs
* [`061307687`](https://github.com/siderolabs/talos/commit/0613076873bbd2d763da30ae2e9e1903486f7cb8) release(v1.12.0-rc.0): prepare release
* [`bc4de5b79`](https://github.com/siderolabs/talos/commit/bc4de5b7926a9a2e7a7af9da4763effb5c33693e) fix: constants file
* [`4a15763a9`](https://github.com/siderolabs/talos/commit/4a15763a962cad0c020e01f66948ba1f326c9201) docs: update release notes
* [`297336549`](https://github.com/siderolabs/talos/commit/29733654902be5cb72b71a9a64ea0ed3c0a0f011) fix: correct condition to use UKI cmdline in GRUB
* [`0ac58929d`](https://github.com/siderolabs/talos/commit/0ac58929db6960ef91c1bcfbc891264e18e1e930) docs: drop machine.network example
* [`184a45c40`](https://github.com/siderolabs/talos/commit/184a45c405530c73c31d5b6c642cda4ddd1772ca) test: bird2 extension
* [`8eac9f37d`](https://github.com/siderolabs/talos/commit/8eac9f37d9dddc507c988cfb187b939a5624f563) docs: add omni join token example to create qemu command
* [`e79a94d57`](https://github.com/siderolabs/talos/commit/e79a94d57781d6ede61e6205f6f5d0f0708a8ddb) fix: adapt SELinuxSuite.TestNoPtrace to new strace version
* [`7a1bb4c26`](https://github.com/siderolabs/talos/commit/7a1bb4c26a99c7f4e37196b40aced6334eeda731) fix: mark secureboot as supported for metal
* [`5c6ee6ace`](https://github.com/siderolabs/talos/commit/5c6ee6aceeb87785c08a05f2ddc6b7cbcad0bc9a) fix: clear provisioning data on SideroLink config change
* [`2e6fe4684`](https://github.com/siderolabs/talos/commit/2e6fe4684b98ca4432284b7b51dfcd1a8b91a03c) feat: update Linux to 6.18.0
* [`473bc17c1`](https://github.com/siderolabs/talos/commit/473bc17c199165dd0f925981753dec431cc5613b) feat: update Kubernetes to 1.35.0-rc.0
* [`6dc8e82b3`](https://github.com/siderolabs/talos/commit/6dc8e82b31d095a357b9f6d99420bb860e51261c) fix: add a timeout for DNS resolving for NTP
* [`a7dbbbd4d`](https://github.com/siderolabs/talos/commit/a7dbbbd4d87feeace427e4c63f67880c72f7cd22) fix: don't disable LACP by default
* [`3ca342c09`](https://github.com/siderolabs/talos/commit/3ca342c0979ffcfe7bee95a4e56c98ddece8abb5) chore: fix longhorn test
* [`364ebb6ba`](https://github.com/siderolabs/talos/commit/364ebb6baf3c77a1e2dd28d83b6af7cfe821e1e8) fix: selection of boot entry
* [`aa286d3f6`](https://github.com/siderolabs/talos/commit/aa286d3f6eb28a813c982a9cc1230c138e56b33a) feat: update Kubernetes to 1.35.0-beta.0
* [`f4891eebb`](https://github.com/siderolabs/talos/commit/f4891eebb192d2895f27f85502fd223290217d90) feat: implement logs persistence
* [`c9a4f95b4`](https://github.com/siderolabs/talos/commit/c9a4f95b42c3347266f60215558f6bde77d4f8a5) release(v1.12.0-beta.1): prepare release
* [`d321d7da0`](https://github.com/siderolabs/talos/commit/d321d7da04fa87e0622f6ec7b5311d5578c534ba) chore: correct condition for running k8s integration tests
* [`736f32a80`](https://github.com/siderolabs/talos/commit/736f32a8077aea0f4a72f3545571882b9d79207c) chore: disable k8s integration tests for 1GiB worker nodes
* [`d9de616c4`](https://github.com/siderolabs/talos/commit/d9de616c48056fc079e693439d4c91a85e154222) chore(ci): skip multipath extension tests
* [`57d6683cd`](https://github.com/siderolabs/talos/commit/57d6683cde0195194acf6880ee85c406216fecc1) chore: update pkgs and tools version
* [`949323ab5`](https://github.com/siderolabs/talos/commit/949323ab51bf5cb95922af7169b698d333c5c9ab) feat: present kernel log as `talosctl logs kernel`
* [`7531fcbc7`](https://github.com/siderolabs/talos/commit/7531fcbc76f3e59e2e8af823d72ffad2cfcaa40a) test: fix flaky LinkSpec/Wireguard test
* [`1dbc64d69`](https://github.com/siderolabs/talos/commit/1dbc64d698f6654e8f8ca5baa13ae9d56745fe6a) fix: simplify OOM expression
* [`0ffb1d857`](https://github.com/siderolabs/talos/commit/0ffb1d8577c9b4da0850a36e80708122b93de303) fix: trim trailing dots from certificate SANs
* [`9a2f6d9c9`](https://github.com/siderolabs/talos/commit/9a2f6d9c9ec5670a12fb033935661f70a80da503) fix: support specifying patch file without '@' symbol
* [`582b0feab`](https://github.com/siderolabs/talos/commit/582b0feab2845d3265cdc852adac78a723953408) fix: assign value of multicast setting properly
* [`16aa6ac47`](https://github.com/siderolabs/talos/commit/16aa6ac471d98b5cdea11d7a4d22ea1048cbd2ce) feat: update etcd to 3.6.6
* [`4396f09c8`](https://github.com/siderolabs/talos/commit/4396f09c8c82ca15b7c09dde8ff1c69a1fe32b08) docs: add API Server Cipher Suites changelog
* [`fdf6fe8e6`](https://github.com/siderolabs/talos/commit/fdf6fe8e6299d620abb3f5c23dcab3cb38fb9367) feat: update TLS cipher suites for API server
* [`139cce3b4`](https://github.com/siderolabs/talos/commit/139cce3b45a7643144aac3042d2bf291e097199d) fix: add CA subject to generated certificate
* [`9b294af22`](https://github.com/siderolabs/talos/commit/9b294af225677a87524491ebd2f21106931dead1) feat: generate mirrors patch
* [`15465f0c5`](https://github.com/siderolabs/talos/commit/15465f0c513ed46886c9f4179c996368843a2daf) fix: add more resilient move
* [`b4147e3a1`](https://github.com/siderolabs/talos/commit/b4147e3a17eebc775cc8ae6087ded6fced11a261) feat: extend flags of cache-cert-gen
* [`72d3d1c9f`](https://github.com/siderolabs/talos/commit/72d3d1c9f53e9b62c189a6369a3060aee4c98d9c) chore: remove spammy 'clean up unused volumes' logs
* [`d6c78de84`](https://github.com/siderolabs/talos/commit/d6c78de84745f27f3051c971451339e760c71397) feat: support TALOS_HOME env var
* [`4040e0814`](https://github.com/siderolabs/talos/commit/4040e0814fc186b2f4e1a2c25520ac08c4d07633) feat: implement multicast setting
* [`eb636dc1f`](https://github.com/siderolabs/talos/commit/eb636dc1f96d1739f1858c4bf825cedc3e0d11e2) feat: add multicast to linkconfig
* [`e34e458c4`](https://github.com/siderolabs/talos/commit/e34e458c4b141ace9604a49b890b2714a59a614e) feat: update dependencies
* [`36152d278`](https://github.com/siderolabs/talos/commit/36152d2787f0cbf3b2efda9c30596f991a811022) fix: add riscv64 talosctl to release artifacts
* [`aebbbaf27`](https://github.com/siderolabs/talos/commit/aebbbaf2746956dc5f88cce6a95061ba447bb36a) feat: support relative voume size
* [`3d997d742`](https://github.com/siderolabs/talos/commit/3d997d7421f3d1b3fda55c92d0e11d75d16daf26) release(v1.12.0-beta.0): prepare release
* [`e62384ba3`](https://github.com/siderolabs/talos/commit/e62384ba34031d43fadebdc84a7d31dd41bf0678) fix: re-creating STATE after partition drop
* [`6919d232a`](https://github.com/siderolabs/talos/commit/6919d232abbaaf44120b9c882e2bc27e4b95deee) docs: update kernel args size
* [`887b296dc`](https://github.com/siderolabs/talos/commit/887b296dc5b111cf54961c1346c4dca4744ccdf9) test: randomize MAC addresses used in the unit-tests
* [`6063fbf91`](https://github.com/siderolabs/talos/commit/6063fbf9124d1953d3bd933bed7f70d42ede2afb) feat: update dependencies
* [`542a67a06`](https://github.com/siderolabs/talos/commit/542a67a066a842a5673755323a3936894b0825ef) feat: add riscv64 build of talosctl
* [`68560b53a`](https://github.com/siderolabs/talos/commit/68560b53ab81335057c0c5524af6f6d2b6882bcf) fix: split volume/disk locators
* [`2c3d30e94`](https://github.com/siderolabs/talos/commit/2c3d30e94f426f2567e9cb97cc3ca9499f53cc7f) docs: fix image-cache-path flag description
* [`93f2e87c2`](https://github.com/siderolabs/talos/commit/93f2e87c2d00c69aacc5f4422182db01b9e617fd) feat: shorthand for generating secrets to stdout
* [`5e1de0035`](https://github.com/siderolabs/talos/commit/5e1de003596837ffe4cf9dd90df4ea121fa2eacc) feat: implement time and resolvers multi-doc configuration
* [`399240be3`](https://github.com/siderolabs/talos/commit/399240be3a51c7053afb9ac60b9e19bd05857615) feat: drop partitions on reset with system partitions wipe
* [`5cca96655`](https://github.com/siderolabs/talos/commit/5cca966557651bb3018ba15d01e0b87146e508fe) feat: add new rockchip sbcs
* [`00fe50d86`](https://github.com/siderolabs/talos/commit/00fe50d868b0463fa32f56ec154bd92bae732f11) fix: uefi bootorder setting
* [`3a881184b`](https://github.com/siderolabs/talos/commit/3a881184bf149410b93657e885796ecf5005b547) chore: improve error handling for system disk reset
* [`859194e67`](https://github.com/siderolabs/talos/commit/859194e6780018ec8e637e87884aa16d3a14cfa6) chore: extract system+user volume config transformers, test
* [`308c6bc41`](https://github.com/siderolabs/talos/commit/308c6bc414d5c6c207bc021ca2949df602725e52) feat: add full disk volumes
* [`82ac1119e`](https://github.com/siderolabs/talos/commit/82ac1119ec102cc591935bbf0afb73431832b775) feat: implement new registry configuration
* [`106f45799`](https://github.com/siderolabs/talos/commit/106f45799d29c7436592b9f1194f6beeed5e394a) feat: update Linux kernel with userfaultfd/VDPA
* [`721a1e0d7`](https://github.com/siderolabs/talos/commit/721a1e0d7cc0cb3eb4d957510accff7762ff366c) chore: rename+improve `client.ErrEventNotSupported`
* [`43f4e317f`](https://github.com/siderolabs/talos/commit/43f4e317f1976762f2999e71ccd6761248a85f12) fix: race between VolumeConfigController and UserVolumeConfigController
* [`66c01a706`](https://github.com/siderolabs/talos/commit/66c01a706f0b1dba88e30dbc1781d7fb7ef57756) chore: deprecate interactive installer mode
* [`957770f65`](https://github.com/siderolabs/talos/commit/957770f65af0d50670b7bbe3758246ced37e9a3e) feat(machined): add panic/force mode reboot
* [`60be0daf8`](https://github.com/siderolabs/talos/commit/60be0daf8414a69b1a60970b14aceb872b31e415) feat: implement multi-doc Wireguard config
* [`cf014cb5d`](https://github.com/siderolabs/talos/commit/cf014cb5d3294ecdcf769315f4795fb8f82a239f) fix: only set default bootloader if none is set
* [`e9b016f80`](https://github.com/siderolabs/talos/commit/e9b016f809d83da33e57492df4a96d68a270ed8c) fix: use strict platform match when pulling images
* [`fafab391b`](https://github.com/siderolabs/talos/commit/fafab391b4d3947daad014438a833ae67b8995fe) feat: update Kubernetes to 1.35.0-alpha.3
* [`7bf3aaca9`](https://github.com/siderolabs/talos/commit/7bf3aaca9129ad40d49f9eadf7ad9be23cf99b32) feat: allow glibc aarch64 so files in extensions
* [`c8561ee2d`](https://github.com/siderolabs/talos/commit/c8561ee2d04c7f9f06c9ec1b3be34ef2a7057efc) feat: implement bridge multi-document config
* [`f4ad3077b`](https://github.com/siderolabs/talos/commit/f4ad3077b0c56b200a37e97abd1a51c63a04c648) feat: implement bond multi-doc configuration
* [`75fe47582`](https://github.com/siderolabs/talos/commit/75fe475828580d9b9a18a2fde0e59f7a9f047ca3) fix: stop attaching to tearing down mount parents
* [`c93a9c6b4`](https://github.com/siderolabs/talos/commit/c93a9c6b41396fe8f8f3f49f475d622e4a45b689) fix: improve OOM controller stability and make test strict on false positives
* [`021bbfefb`](https://github.com/siderolabs/talos/commit/021bbfefbecc688fc4c61876c264416f72c7a7a2) feat: update Go 1.25.4, containerd 2.1.5
* [`e25db484f`](https://github.com/siderolabs/talos/commit/e25db484f54414dcd7b8f08c1a741b58435e52f5) test: disable parallelism in Longhorn tests
* [`54b93aff0`](https://github.com/siderolabs/talos/commit/54b93aff0c372761dfe9621a782a347b6877c2e9) feat: update Linux 6.17.7, runc 1.3.3
* [`2af69ff35`](https://github.com/siderolabs/talos/commit/2af69ff35712ac843c66e30fdf6a380aae2ed499) fix: provide minimal platform metadata always
* [`92eeaa482`](https://github.com/siderolabs/talos/commit/92eeaa4826cf71a5962da8ea055a11732fbc851e) fix: update YAML library
* [`aa24da9aa`](https://github.com/siderolabs/talos/commit/aa24da9aab9c5dc2f51401ae8ba0161e63c09924) fix: bump kubelet credendial provider config to v1
* [`335f91761`](https://github.com/siderolabs/talos/commit/335f9176151f7d45c0f847abecb20184483a6cd3) feat: add short -c flag for --cluster
* [`4c095281b`](https://github.com/siderolabs/talos/commit/4c095281be93cb11290eb43f60b4cc1a168bef17) fix: set a timeout for SideroLink provision API call
* [`75e4c4a59`](https://github.com/siderolabs/talos/commit/75e4c4a598181a18638aadcb77c89fbe762c6b9f) fix: log duplication on log senders
* [`e3cbc92c0`](https://github.com/siderolabs/talos/commit/e3cbc92c0579beb0262d2d2d6a0d00d56bbbdc17) fix: add video kernel module to arm
* [`d69305a67`](https://github.com/siderolabs/talos/commit/d69305a670ac982ba7dd00cfc8e7cf736cbfb385) fix: userspace wireguard handling
* [`ee5fee7c8`](https://github.com/siderolabs/talos/commit/ee5fee7c8a0f482894534bd2f8e5b0c2b2076854) fix: image-signer commands
* [`be028b67a`](https://github.com/siderolabs/talos/commit/be028b67a068c0d0d4465725c96b28ad9b276e8a) feat: add support for multi-doc VLAN config
* [`f3df0f80b`](https://github.com/siderolabs/talos/commit/f3df0f80b9d64e282bf163ba04ed9363e40865a3) feat: add directory backed UserVolumes
* [`0327e7790`](https://github.com/siderolabs/talos/commit/0327e77902a05978c79a9efb92bc50a792e4e0be) feat: add support for dashboard custom console parameter
* [`fed948b8a`](https://github.com/siderolabs/talos/commit/fed948b8ae416db886df6ed783bde60aae2a25c8) release(v1.12.0-alpha.2): prepare release
* [`fb4bfe851`](https://github.com/siderolabs/talos/commit/fb4bfe851c7c308eeaf4a11e0ac5c944f66dc0c4) chore: fix LVM test
* [`f4ee0d112`](https://github.com/siderolabs/talos/commit/f4ee0d1128ba2f35d54ec3d35a83fc62fd222f2e) chore: disable VIP operator test
* [`288f63872`](https://github.com/siderolabs/talos/commit/288f6387260843570d53d28a4d77e564b3182979) feat: bump deps
* [`b66482c52`](https://github.com/siderolabs/talos/commit/b66482c529beda8b1abf9ed6b71ece354c1540be) feat: allow disabling injection of extra cmdline in cluster create
* [`704b5f99e`](https://github.com/siderolabs/talos/commit/704b5f99e6bef4410629427ac65fd2742ddb335d) feat: update Kubernetes to 1.35.0-alpha.2
* [`1dffa5d99`](https://github.com/siderolabs/talos/commit/1dffa5d9965a6c7d872f052bfb1750ea550671c2) feat: implement virtual IP operator config
* [`43b1d7537`](https://github.com/siderolabs/talos/commit/43b1d7537507a916629cc2d6db7440a99ffcb748) fix: validate provisioner when destroying local clusters
* [`b494c54c8`](https://github.com/siderolabs/talos/commit/b494c54c81e6ca81cef8ce26da772c1fc336ea8d) fix: talos import on non-linux
* [`61e95cb4b`](https://github.com/siderolabs/talos/commit/61e95cb4b7b354d175d1dfce3d0fa43deefad187) feat: support bootloader option for ISO
* [`d11072726`](https://github.com/siderolabs/talos/commit/d110727263c57c02392f201938d2b71976b8c4d6) fix: provide offset for partitions in discovered volumes
* [`39eeae963`](https://github.com/siderolabs/talos/commit/39eeae96311be2b8e2d3660d878f852ba92ca064) feat: update dependencies
* [`9890a9a31`](https://github.com/siderolabs/talos/commit/9890a9a31deb11ab170b94c667143314db08f76f) test: fix OOM test
* [`c0772b8ed`](https://github.com/siderolabs/talos/commit/c0772b8eda429675a06899b9c4a4d1dd7d5f6a5f) feat: add airgapped mode to QEMU backed talos
* [`ac60a9e27`](https://github.com/siderolabs/talos/commit/ac60a9e27deed63db0e4e61ffa30d46f4cab590a) fix: update test for PCI driver rebind/IOMMU
* [`6c98f4cdb`](https://github.com/siderolabs/talos/commit/6c98f4cdb049c58ef4f6e8193ef66c2338a2877d) feat: implement new DHCP network configuration
* [`da92a756d`](https://github.com/siderolabs/talos/commit/da92a756d9668fa043b4794db45d5c985d8ea4a6) fix: drop 'ro' falg from defaults
* [`28fd2390c`](https://github.com/siderolabs/talos/commit/28fd2390cb6e02f400bb237dd674c7d0d40f8ed3) fix: imager build on arm64
* [`4e12df8c5`](https://github.com/siderolabs/talos/commit/4e12df8c5c27ae115c4eac70a7e2fceb03dac5f5) test: integration test for OOM controller
* [`7e498faba`](https://github.com/siderolabs/talos/commit/7e498faba93f972ba82edf41550d3b94256e83e9) feat: use image signer
* [`eccb21dd3`](https://github.com/siderolabs/talos/commit/eccb21dd3ba03eb4ab03c4da87a51a4e3d8da49a) feat: add presets to the 'cluster create qemu' command
* [`ec0a813fa`](https://github.com/siderolabs/talos/commit/ec0a813facf5be5ca3e9ba65924ae18b2b05a7d9) feat: unify cmdline handling GRUB/systemd-boot
* [`37e4c40c6`](https://github.com/siderolabs/talos/commit/37e4c40c6a2477e45bbf067effc4389d4639c905) fix: skip module signature tests on docker provisioner only
* [`8124efb42`](https://github.com/siderolabs/talos/commit/8124efb42fd5a3eb81f41e84974e4242246ca7c4) fix: cache e2e
* [`4adcda0f5`](https://github.com/siderolabs/talos/commit/4adcda0f5427e1bae49f6dda58318324a3b24ac5) fix: reserve the apid and trustd ports from the ephemeral port range
* [`ced57b047`](https://github.com/siderolabs/talos/commit/ced57b047a389e26f7e5bfa3efab5b64f3fced87) feat: support optionally disabling module sig verification
* [`1e5c4ed64`](https://github.com/siderolabs/talos/commit/1e5c4ed644cbc60d8518fe4298e63a5cf5dc8cf5) fix: build talosctl image cache-serve non-linux
* [`dbdd2b237`](https://github.com/siderolabs/talos/commit/dbdd2b237e0aefbba439b90472abf9ec7eea6aa6) feat: add static registry to talosctl
* [`77d8cc7c5`](https://github.com/siderolabs/talos/commit/77d8cc7c589a190c8cb86e6e1684233129b648a1) chore: push `latest` tag only on main
* [`59d9b1c75`](https://github.com/siderolabs/talos/commit/59d9b1c75dbff09e405906ebcfb3ad1a69cb8f4b) feat: update dependencies
* [`bf6ad5171`](https://github.com/siderolabs/talos/commit/bf6ad51710c367764e582ccc1fb77b4d989c874d) feat: add back install script
* [`da451c5ba`](https://github.com/siderolabs/talos/commit/da451c5ba4ee97e7ef108bb6d73d5aa8bc7c72fd) chore: drop documentation except for fresh reference
* [`2f23fedeb`](https://github.com/siderolabs/talos/commit/2f23fedeb725a5786b6ffac2aef8125eecd6cb6e) fix: file leak in reading cgroups
* [`b412ffdbc`](https://github.com/siderolabs/talos/commit/b412ffdbc29d77a81aed88be62f21bc2999afcde) docs: update README.md for docs link
* [`8dc51bae7`](https://github.com/siderolabs/talos/commit/8dc51bae79a37b56c058d40787dbda6e828fd0d3) feat: add drm_gpuvm and drm_gpusvm_helper modules
* [`4ca58aeb8`](https://github.com/siderolabs/talos/commit/4ca58aeb81145cb7ebef071865b3d853a4712729) fix: make Akamai platform usable
* [`061f8e76f`](https://github.com/siderolabs/talos/commit/061f8e76fd58906ff823a0e467d6efcf5161ed9f) feat: bump pkgs
* [`a9fa852da`](https://github.com/siderolabs/talos/commit/a9fa852dadd75740d73588fd2156f6f1ad782fdd) feat: update uefi image to talos linux logo
* [`04753ba69`](https://github.com/siderolabs/talos/commit/04753ba6983b6ff2754cf62b8d60cc6065921dbd) feat: update go to 1.25.2
* [`9a42b05bd`](https://github.com/siderolabs/talos/commit/9a42b05bdac2bf0cbbc97d040be7860f48c69386) feat: implement link aliasing
* [`d732bd0be`](https://github.com/siderolabs/talos/commit/d732bd0be73c3d17d140c00be0e9d27ea621909b) chore(ci): run only nvidia tests for NVIDIA workflows
* [`8d1468209`](https://github.com/siderolabs/talos/commit/8d1468209aa28f59df9dc52466c506defa8c3cc3) fix: stop populating apiserver cert SANs
* [`02473244c`](https://github.com/siderolabs/talos/commit/02473244c17ef0149515f300bcd201f9347acabc) fix: wait for mount status to be proper mode
* [`825622d90`](https://github.com/siderolabs/talos/commit/825622d90a7716f7b6027651a5b9389173432393) fix: resource proto definitions
* [`2c6003e79`](https://github.com/siderolabs/talos/commit/2c6003e790003f6ef1a03b8d2af8030fb57c5d02) docs: add Project Calico installation in two mode
* [`4fb4c8678`](https://github.com/siderolabs/talos/commit/4fb4c86780def54eed4d999b1f0ce93042269076) feat: add disk.EnableUUID to generated ova
* [`33fb48f8f`](https://github.com/siderolabs/talos/commit/33fb48f8f90ccf44e95c93ac7ec1adcd1b4e0373) fix: add dashboard spinner
* [`053fd0bd4`](https://github.com/siderolabs/talos/commit/053fd0bd4d324bc21e076b3a30466ed61c7684e1) feat: update Linux to 6.17
* [`34e107e1b`](https://github.com/siderolabs/talos/commit/34e107e1bd14b0a56ebfa0c65e0c7da715976d99) docs: fix broken link
* [`dfbece56b`](https://github.com/siderolabs/talos/commit/dfbece56bd45e95c9ec477af4b53ffcefdfec66c) docs: update the kubespan docs
* [`8b041a72c`](https://github.com/siderolabs/talos/commit/8b041a72ca9c07985c024c1136c85c85df92beda) docs: update scaleway.md
* [`435dcbf82`](https://github.com/siderolabs/talos/commit/435dcbf820cd9f8cc9fecc0f7d42819acef36106) fix: provide nocloud metadata with missing network config
* [`ec3bd878f`](https://github.com/siderolabs/talos/commit/ec3bd878f9770ceb932b654aabad1711880da829) refactor: remove the go-blockdevice v1 completely
* [`33544bde9`](https://github.com/siderolabs/talos/commit/33544bde9c15745f4ae692c7647d661b32d4bed4) fix: minor improvements to fs
* [`fd2eebf7f`](https://github.com/siderolabs/talos/commit/fd2eebf7fa4831d33383a53d6d058c74789553e4) feat: create merge patch from diff of two machine configs
* [`eadbdda94`](https://github.com/siderolabs/talos/commit/eadbdda9471289fae5159c8cc024a735a1547807) fix: uefi boot order setting
* [`cd9fb2743`](https://github.com/siderolabs/talos/commit/cd9fb274342c5a973b3d087b991a7eea5df4142a) fix: support secure HTTP proxy with gRPC dial
* [`adf87b4b9`](https://github.com/siderolabs/talos/commit/adf87b4b931ded1edeb64217b0e9d5edfd046004) feat: update Flannel to v0.27.4
* [`5dfb7e1fe`](https://github.com/siderolabs/talos/commit/5dfb7e1fe7d9cc6db3e4c2b6f587e641b4a0842b) feat: serve etcd image from registry.k8s.io
* [`5ca841804`](https://github.com/siderolabs/talos/commit/5ca8418049e3b878585014a3764021f2d30a0df7) fix: nftables flaky test
* [`a940e45a7`](https://github.com/siderolabs/talos/commit/a940e45a7fe041b17437f774eb52b9f3a42e3633) feat: generate list of images required to build talos
* [`3472d6e79`](https://github.com/siderolabs/talos/commit/3472d6e79caa13fd42df7774101397b0a30f62f5) fix: revert "chore: use new mount/v3 package in efivarfs"
* [`42c0bdbf3`](https://github.com/siderolabs/talos/commit/42c0bdbf320bf24311b2d56b2e0f7155e86b3713) feat: add provisioner flag to images default command
* [`6bc0b1bcf`](https://github.com/siderolabs/talos/commit/6bc0b1bcf7d9dc9f2417a7db63d1e76e7ddc6aa3) feat: drop and lock deprecated features
* [`362a8e63b`](https://github.com/siderolabs/talos/commit/362a8e63b798c4a4fc31fe5e728d2429fc953166) fix: change the compression format
* [`6e58f58aa`](https://github.com/siderolabs/talos/commit/6e58f58aaeb6e16883d8dc8757ad92b6b6da7e84) fix: mkdir artifacts path
* [`3165a2b84`](https://github.com/siderolabs/talos/commit/3165a2b84cb80dd5fd09bf496fdccaf1628593d0) release(v1.12.0-alpha.1): prepare release
* [`e455c7ea9`](https://github.com/siderolabs/talos/commit/e455c7ea9c919a2f70ddecceaa8f3b4e25566048) chore: use testing/synctest in tests
* [`7f048e962`](https://github.com/siderolabs/talos/commit/7f048e962e217687ab67ed7027c5228e8ccb7d16) feat: update dependencies
* [`fe36b3d32`](https://github.com/siderolabs/talos/commit/fe36b3d3200db57f3e21017ff7a4808b330a1d55) fix: stop returning EINVAL on remount of detached mounts
* [`c6279e04c`](https://github.com/siderolabs/talos/commit/c6279e04c45504af243c0aef9f255317426b4ca0) chore: use new mount/v3 package in efivarfs
* [`d5197effb`](https://github.com/siderolabs/talos/commit/d5197effb0b48290d613140b68796cb8f30b9a70) feat: update etcd 3.6.5, CoreDNS 1.12.4
* [`33714b715`](https://github.com/siderolabs/talos/commit/33714b7158a0d569be1d0b1d7b012280856db484) feat: release cloud image using factory
* [`d10a2747e`](https://github.com/siderolabs/talos/commit/d10a2747e0e835876aff158e6b6f7882cef9fa44) docs: deprecate JSON6902 patches and interactive installer
* [`1e604cbf5`](https://github.com/siderolabs/talos/commit/1e604cbf514bece1e112d8afd5d1cd6ccb1045c3) fix: don't set broadcast for /31 and /32 addresses
* [`65a66097a`](https://github.com/siderolabs/talos/commit/65a66097a05e5c0e2334d5eff494a0e71534716f) refactor: split cluster create logic into smaller parts
* [`ab847310e`](https://github.com/siderolabs/talos/commit/ab847310efde540b5bfe17570b99af1bb705832b) fix: provide refreshing CA pool (resolvers)
* [`d63c3ed7d`](https://github.com/siderolabs/talos/commit/d63c3ed7db2b22f7e394fc45d101d03cba463177) docs: update secureboot docs
* [`493f7ed9d`](https://github.com/siderolabs/talos/commit/493f7ed9d2710eb240eab6b6ab532f41abc818c1) feat: support embedded config
* [`251df70f6`](https://github.com/siderolabs/talos/commit/251df70f6d33f1d5a3b1b9e4c0c249d8bc85c4b3) feat: add a userspace OOM controller
* [`7bae5b40b`](https://github.com/siderolabs/talos/commit/7bae5b40b4f22f0f07a586ebd9cda9436086a5f8) feat: implement link configuration
* [`724857dec`](https://github.com/siderolabs/talos/commit/724857decb95ddeebb2ac5d33c38a71bf7512805) fix(ci): skip netbird extension for tests
* [`e06a08698`](https://github.com/siderolabs/talos/commit/e06a086989331f28406e8d4234e02d9a6b83f87d) fix: default gateway as string
* [`7ed07412e`](https://github.com/siderolabs/talos/commit/7ed07412e963e6ee91615adbea095944aa6a56e5) fix: uefi boot entry handling logic
* [`ea4ed165a`](https://github.com/siderolabs/talos/commit/ea4ed165ad860a5beea17ca2d404bdaa6e5ad933) refactor: efivarfs mock and tests
* [`1fca111e2`](https://github.com/siderolabs/talos/commit/1fca111e24bcae81b78f007e67b71c9155c0169f) feat: support setting wake-on-lan for Ethernet
* [`94f78dbe7`](https://github.com/siderolabs/talos/commit/94f78dbe798cb227a0c38b70a1d6840803989290) docs: add a documentation for running Talos in KVM
* [`46902f8fd`](https://github.com/siderolabs/talos/commit/46902f8fdee257a09be4bc1753c6b3f845ef8089) docs: add TrueFullstaq to adopters
* [`a28e5cbd5`](https://github.com/siderolabs/talos/commit/a28e5cbd50d11aa6c253a6a9ce1999b9d45effad) chore: update pkgs and tools
* [`7cf403db8`](https://github.com/siderolabs/talos/commit/7cf403db8ca0e1719195001895cfbc12835b0fdd) docs: step-by-step scaleway documentation to get an image
* [`687285fa2`](https://github.com/siderolabs/talos/commit/687285fa26ec42dadbfb72580099f6e20bbaf85e) docs: remove 'curl' in wget command
* [`9db6dc06c`](https://github.com/siderolabs/talos/commit/9db6dc06c3010cd89ce4cb0ec0bde178db0447a4) feat: stop mounting state partition
* [`53ce93aae`](https://github.com/siderolabs/talos/commit/53ce93aaed3bd5bfcbe926fa69ca3b4b8b45c74f) test: try to clear connection refused more aggressively
* [`51db5279c`](https://github.com/siderolabs/talos/commit/51db5279c423e4b8637a05e52b26dfc5aa719cbc) fix: bump trustd memory limit
* [`25204dc8a`](https://github.com/siderolabs/talos/commit/25204dc8a8df79bc876a0bec2492e1147a81d954) fix(machined): change `constants.MinimumGOAMD64Level` using build tag
* [`9cd2d794d`](https://github.com/siderolabs/talos/commit/9cd2d794d060b637dbac5263ae417a4e83d54efe) feat: ship nft binary with Talos rootfs
* [`b1416c9fe`](https://github.com/siderolabs/talos/commit/b1416c9fe1d5ea9cd68f9b6b766a288a267cee61) feat: record last log the failed service
* [`0b129f9ef`](https://github.com/siderolabs/talos/commit/0b129f9efdf57dd9692f7cece6b97719a7ccf80e) feat: enforce more KSPP and hardening sysctls
* [`11872643c`](https://github.com/siderolabs/talos/commit/11872643c310212c52b4fd7e13b6cc7d6ec7e4fc) chore: drop docs folder
* [`d30fdcd88`](https://github.com/siderolabs/talos/commit/d30fdcd88f421824cf17b9ecec25be7c8044e857) chore: pass in github token to imager
* [`b88f27d80`](https://github.com/siderolabs/talos/commit/b88f27d804d60a706f598b50676dad5dd2a9726a) chore: make reset test code a bit better
* [`1cde53d01`](https://github.com/siderolabs/talos/commit/1cde53d0173fd1ae637855e15fe34bb74bb027a0) test: fix several issues with tests
* [`16cd127a0`](https://github.com/siderolabs/talos/commit/16cd127a04bb5fc907b7ca04f1c81d4c7150eab2) docs: add docs on updating image cache
* [`c3ae92b14`](https://github.com/siderolabs/talos/commit/c3ae92b1424d4a2c9bc18cfa394b10eda6c9a20f) fix: build kernel checks only on linux
* [`2120904ec`](https://github.com/siderolabs/talos/commit/2120904ec534a91f66dcea419b5a29e36a16f6e4) feat: create detached tmpfs
* [`6bbee6de5`](https://github.com/siderolabs/talos/commit/6bbee6de5b18b25deb4e6f515251187e259aa424) docs: remove 'ceph-data' from volume examples/docs
* [`07acb3bd2`](https://github.com/siderolabs/talos/commit/07acb3bd2d4f92e80706d1835130bbe6e944d096) fix: use correct order to determine SideroV1 keys directory path
* [`2d57fa002`](https://github.com/siderolabs/talos/commit/2d57fa00281f8090b85097c66df634101b0cde79) fix: trim zero bytes in the DHCP host & domain response
* [`451cb5f78`](https://github.com/siderolabs/talos/commit/451cb5f78fac3b2ddfec7d545629fe8c88ea2367) docs: clarify disk partition confusion
* [`a2122ee5c`](https://github.com/siderolabs/talos/commit/a2122ee5cb9c84f33e0c4b30e9223bb239621d55) feat: implement HostConfig multi-doc
* [`69ab076b4`](https://github.com/siderolabs/talos/commit/69ab076b4d6e52484677ee7f68a853dc4edfe2bc) fix: re-create cgroups when restarting runners
* [`297b5cc28`](https://github.com/siderolabs/talos/commit/297b5cc2856710b74b4e0e46b00ae33aea4c1bf7) docs: add docs on node labels
* [`e168512dd`](https://github.com/siderolabs/talos/commit/e168512dd020da9eac654dae2ba891cf33415c44) fix: apply 'ro' flag to iso9660 filesystems
* [`7f7acfbb9`](https://github.com/siderolabs/talos/commit/7f7acfbb9f10c243d0b132c1ef079cb77d2727e0) docs: fix typo in doc
* [`d57882b18`](https://github.com/siderolabs/talos/commit/d57882b1830504fe4bfd5344edae613168db7f0e) feat: update Kubernetes to 1.34.1
* [`f85f82f32`](https://github.com/siderolabs/talos/commit/f85f82f32f098f97588f404550f72d64786fe329) test: fix flakiness in RawVolumes test
* [`82569e319`](https://github.com/siderolabs/talos/commit/82569e319eb57b1199db6bfd3e612fb771c8c7cd) feat: update Linux 6.16.6
* [`2fd2ab4e4`](https://github.com/siderolabs/talos/commit/2fd2ab4e43e06910154705d6ef1d0576a7c04a2b) fix: remove CoreDNS cpu limit
* [`ce9bc32a0`](https://github.com/siderolabs/talos/commit/ce9bc32a08695873d9054afe2608a76cf7c6088a) chore(ci): rekres to use new runner groups
* [`8b64f68f6`](https://github.com/siderolabs/talos/commit/8b64f68f6946c2979f6fe2bf617f31639a927bf8) test: improve test stability
* [`272cb860d`](https://github.com/siderolabs/talos/commit/272cb860d4cfb8464b29ff31567e25fe6c275849) chore: drop the --input-dir flag from the cluster create command
* [`1b6533675`](https://github.com/siderolabs/talos/commit/1b65336752933acdcbf681767785157714866f88) docs: add note about ca-signed certs for secureboot
* [`d3f88f50c`](https://github.com/siderolabs/talos/commit/d3f88f50c5394536ee80d19464359408a37d81ff) docs: document talos vip failover behavior
* [`005fc8bd5`](https://github.com/siderolabs/talos/commit/005fc8bd50fbc4b15b26032b43d1d32c1da22f11) docs: add docs on syncing configs after a kube upgrade
* [`4d876d9af`](https://github.com/siderolabs/talos/commit/4d876d9af9fcc9828f09d05db124fbdce9c17785) feat: update Go to 1.25.1
* [`2b556cd22`](https://github.com/siderolabs/talos/commit/2b556cd22a3563f1d86a648ea6c69a4d45edad76) feat: implement multi-doc StaticHostConfig
* [`a7b776842`](https://github.com/siderolabs/talos/commit/a7b7768420566b6840fc52bb2152e9bf165f8cd3) docs: replace Raspberry Pi 5 links with Talos builder
* [`a349b20ed`](https://github.com/siderolabs/talos/commit/a349b20ed4b3c05dcd0175541b795331f0f7c64d) docs: clarify that talos does not support intermediate ca
* [`895133de9`](https://github.com/siderolabs/talos/commit/895133de99158ce3f50b557b77c81d4f0f9d6b40) feat: support configuring PCR states to bind disk encryption
* [`c1360103b`](https://github.com/siderolabs/talos/commit/c1360103b5e037cf713b7d787436f01e7182821c) docs: fix command for uploading image on Hetzner
* [`43b5b9d89`](https://github.com/siderolabs/talos/commit/43b5b9d8992ad6df37619b3719b57948e4bd9671) fix: correctly handle status-code 204
* [`feeb0d312`](https://github.com/siderolabs/talos/commit/feeb0d312ecacb451e5313390939c7c9349d2ba6) feat: update runc to 1.3.1
* [`421634a14`](https://github.com/siderolabs/talos/commit/421634a1417f529551a75d0bb9be08b73f1120b1) docs: add docs on multihoming
* [`41af2d230`](https://github.com/siderolabs/talos/commit/41af2d230c2dd5dce5bc931f76a2eb69405dc554) refactor: clean up internal cluster creation code
* [`3000d9e43`](https://github.com/siderolabs/talos/commit/3000d9e431deaf952d08da724da40789cd743f2c) fix: don't bootstrap talos cluster if there's no config present
* [`79cb871d0`](https://github.com/siderolabs/talos/commit/79cb871d088e5b1c3a3488610ded14e7a28cec29) feat: use the id of the volume in the mapped luks2 name
* [`6c322710d`](https://github.com/siderolabs/talos/commit/6c322710d64786f19e2e0e39d65596c8dce71952) chore: refactor mount package
* [`ced7186e2`](https://github.com/siderolabs/talos/commit/ced7186e2a5f0634d9441b12a5340f5ca4c451ff) refactor: update COSI to 1.11.0
* [`de2e24fcd`](https://github.com/siderolabs/talos/commit/de2e24fcda590a1ef3f80a5372bb70865a2f47c3) docs: clarify that install-cni image is deprecated
* [`bef8ef509`](https://github.com/siderolabs/talos/commit/bef8ef509380aba259efcc2f5d1f6632e034160b) docs: add docs on cilium's compatibility with kubespan
* [`e5acb10fc`](https://github.com/siderolabs/talos/commit/e5acb10fcceba69060507a35caea21281bdc71cc) feat: update pkgs
* [`c4c1daf0e`](https://github.com/siderolabs/talos/commit/c4c1daf0e2e6675626b974b0c008e101d919c8b5) docs: add info about br_netfilter
* [`5c52ecac3`](https://github.com/siderolabs/talos/commit/5c52ecac364f917e5f45859f680494a08f85cb90) docs: clarify interactive dashboard resolution control
* [`15ecb02a4`](https://github.com/siderolabs/talos/commit/15ecb02a4545639ffb8ba5c6e5a413e53129b619) feat: update Linux kernel (memcg_v1, ublk)
* [`53f18c2f6`](https://github.com/siderolabs/talos/commit/53f18c2f60c84c4b0f944cc343ae1f538e8d1236) fix: enable support for VMWare arm64
* [`3bbe1c0da`](https://github.com/siderolabs/talos/commit/3bbe1c0da5485b6cd3e7fadd8f020e0d0aca406a) docs: add docs on grow flag
* [`b9fb09dcd`](https://github.com/siderolabs/talos/commit/b9fb09dcdbcca60f695ac317c45e18fa092541a8) release(v1.12.0-alpha.0): prepare release
* [`6a389cad3`](https://github.com/siderolabs/talos/commit/6a389cad35f80b27fe9c43db9e701ee9f6f6142a) chore: update dependencies
* [`9d98c2e89`](https://github.com/siderolabs/talos/commit/9d98c2e891258dcf2ef90519d38d0aefb77cd0db) feat: add a cgroup preset for PSI and --skip-cri-resolve
* [`072f77b16`](https://github.com/siderolabs/talos/commit/072f77b1623cdc838093465b7266b26e20a248ea) chore: prepare for future Talos 1.12-alpha.0 release
* [`96f41ce88`](https://github.com/siderolabs/talos/commit/96f41ce8840783f783fcc8e0fd6b43302b9bfe43) docs: update qemu and docker docs
* [`a751cd6b7`](https://github.com/siderolabs/talos/commit/a751cd6b7474a4dc20137e917dbb2229fe9cc8bd) docs: activate Talos v1.11 docs by default
* [`e8f1ec1c5`](https://github.com/siderolabs/talos/commit/e8f1ec1c5bbd8a6cfb68886e6283e7caaf5fb063) docs: fix broken create qemu command v1.11 docs
* [`639f0dfdd`](https://github.com/siderolabs/talos/commit/639f0dfdd88c5596439601f3f9600b3aafb24227) feat: update Linux to 6.16.4
* [`8aa7b3933`](https://github.com/siderolabs/talos/commit/8aa7b3933d07ea45a96844b9c91347a08950e243) fix: bring back linux/armv7 build and update xz
* [`9cae7ba6b`](https://github.com/siderolabs/talos/commit/9cae7ba6b97a67a5d282c6f667ccb4c3e2111447) feat: update CoreDNS to 1.12.3
* [`cfef3ad45`](https://github.com/siderolabs/talos/commit/cfef3ad4544498a47de17f6b05fb8374c35e3dd8) fix: drop linux/armv7 build
* [`42ea2ac50`](https://github.com/siderolabs/talos/commit/42ea2ac5058457dafe666f8d79f08d3c8ee60cfb) fix: update xz module (security)
* [`4fcfd35b9`](https://github.com/siderolabs/talos/commit/4fcfd35b9510f45d0ef7ae3657eb0916d549d2dd) docs: fix module name example
* [`50824599a`](https://github.com/siderolabs/talos/commit/50824599a4fa7b72d563a35a4746ca063becf672) chore: update some tools
* [`bcd297490`](https://github.com/siderolabs/talos/commit/bcd297490c608f593b6dd274945aa2b73c3fd3ee) feat: allow Ed25119 in FIPS mode
* [`5992138bb`](https://github.com/siderolabs/talos/commit/5992138bb981e84dae917f0f0fdafee4049bc5ec) test: ignore one leaking goroutine
* [`d155326c1`](https://github.com/siderolabs/talos/commit/d155326c1206979f30a5355f7bdb23cb051e9b78) docs: add sbc unofficial ports docs
* [`285fa7d22`](https://github.com/siderolabs/talos/commit/285fa7d222be1f5e63c0bb725b206966e2722a3b) docs: add the deploy application docs
* [`527791f09`](https://github.com/siderolabs/talos/commit/527791f0974afe9c8558b82fa19f4354487693ed) feat: update Kubernetes to 1.34.0
* [`a1c0e237d`](https://github.com/siderolabs/talos/commit/a1c0e237d6e047bb59c4fbd48e2c2b9e36dd4808) feat: update Linux to 6.15.11, Go to 1.25
* [`4d7fc25f8`](https://github.com/siderolabs/talos/commit/4d7fc25f8bf20d4489080795a3d0ce0dfb1bc6b8) docs: switch order of wipe disk command
* [`7368a994d`](https://github.com/siderolabs/talos/commit/7368a994df07cc4e50e3709ac766d8062db070a0) feat: add SOCKS5 proxy support to dynamic proxy dialer
* [`d63591069`](https://github.com/siderolabs/talos/commit/d635910697b221aee3e9afa6d9e5b398236b6a21) chore: silence linter warnings
* [`07eb4d7ec`](https://github.com/siderolabs/talos/commit/07eb4d7ec148a7e3c4c6dde080469c1a2fb410fb) fix: set default ram unit to MiB instead of MB
* [`6b732adc4`](https://github.com/siderolabs/talos/commit/6b732adc43684facfd329f424a34a7e4df36d77b) feat: update Linux to 6.12.43
* [`b6410914f`](https://github.com/siderolabs/talos/commit/b6410914f74ce01672fdef7e912e37970909281c) feat: add human readable byte size cli flags
* [`ec70cef99`](https://github.com/siderolabs/talos/commit/ec70cef99005fd7e383fea63b5c23774882fcf28) feat: update NVIDIA drivers and kernel
* [`0879efa69`](https://github.com/siderolabs/talos/commit/0879efa690ad657e4aed251fcbeba8f5645d73ce) feat: update Kubernetes default to v1.34.0-rc.2
* [`f504639df`](https://github.com/siderolabs/talos/commit/f504639df4388619f731196ed8e79a6818b6ed5f) feat: add a user-facing create qemu command
* [`558e0b09a`](https://github.com/siderolabs/talos/commit/558e0b09ab65b353e83b98c9ddf6cb2b67fd060e) test: fix the Image Factory PXE boot test
* [`d73f0a2e5`](https://github.com/siderolabs/talos/commit/d73f0a2e5b788c7b69c2fb827f7111d5f9c8e706) docs: make readme badges consistent
* [`f1369af98`](https://github.com/siderolabs/talos/commit/f1369af98e1f6d48fed137e31237956abbd28b0f) chore: use new filesystem api on STATE partition
* [`366cedbe7`](https://github.com/siderolabs/talos/commit/366cedbe7495ce15bcd0e6c6f7f0add65a41a861) docs: link to kubernetes linux swap tuning
* [`2f5a16f5e`](https://github.com/siderolabs/talos/commit/2f5a16f5e4ae186a309aef5e3d285897d0fe2df1) fix: make --with-uuid-hostnames functionality available to qemu provider
* [`70612c1f9`](https://github.com/siderolabs/talos/commit/70612c1f9fc9056e8a3669ff10a385c4e8e03350) refactor: split the PlatformConfigController
* [`511748339`](https://github.com/siderolabs/talos/commit/51174833997fd9a0a599ab1dde947834b682ab14) docs: add system extension tier documentation
* [`009fb1540`](https://github.com/siderolabs/talos/commit/009fb1540e0b9f5daac6302f42e8813e596fc87c) test: don't run nvidia tests on integration/aws
* [`99674ef20`](https://github.com/siderolabs/talos/commit/99674ef20d34166d60563d4bf46fbbfc57399509) docs: apply fixes for what is new
* [`92db677b5`](https://github.com/siderolabs/talos/commit/92db677b5d32de32ec7e785531b32202e03283b4) fix: image cache lockup on a missing volume
* [`9c97ed886`](https://github.com/siderolabs/talos/commit/9c97ed886b89b2fb84f47866abdf1000839143c4) fix: version contract parsing in encryption keys handling
* [`1fc670a08`](https://github.com/siderolabs/talos/commit/1fc670a08dc7af8eaeabdc7134eb77a5c939df40) fix: dial with proxy
* [`18447d0af`](https://github.com/siderolabs/talos/commit/18447d0afdbcc8fa7db6ae008e4bc4d5b0a0b00a) feat: update Linux to 6.12.41
* [`f65f39b78`](https://github.com/siderolabs/talos/commit/f65f39b78b0c7881e5f51c66ad022c17c2cd4960) fix: provide mitigation CVE-1999-0524
* [`8817cc60c`](https://github.com/siderolabs/talos/commit/8817cc60cfaf4b50f11c38d3b25df7df48382033) fix: actually use SIDEROV1_KEYS_DIR env var if it's provided
* [`b08b20a10`](https://github.com/siderolabs/talos/commit/b08b20a1005256a9e3fc7cae8bcf8eea87f6ac09) feat: use key provider with fallback option for auth type SideroV1
* [`7a52d7489`](https://github.com/siderolabs/talos/commit/7a52d7489c9709708d55f8f001d70700addc7e1e) fix: kubernetes upgrade options for kubelet
* [`ea8289f55`](https://github.com/siderolabs/talos/commit/ea8289f550787593b1cd35f2d8da59aa5311880e) feat: add a user facing docker command
* [`54ad64765`](https://github.com/siderolabs/talos/commit/54ad64765090d90013e4917d1bf494592069beec) chore: re-enable vulncheck
* [`26bbddea9`](https://github.com/siderolabs/talos/commit/26bbddea95669278363c604316ed85986f312d71) fix: darwin build
* [`b5d5ef79e`](https://github.com/siderolabs/talos/commit/b5d5ef79e7a2d76e29a7c872c1c418fffc63b0df) fix: set secs field in DHCPv4 packets
* [`c07911933`](https://github.com/siderolabs/talos/commit/c0791193373e36c35f29c70318432331b4c6ab2a) chore: refactor how tools are being installed
* [`34f25815c`](https://github.com/siderolabs/talos/commit/34f25815c036d2c91bdfddc9c7d40ca2edf677bd) docs: fork docs for v1.12
* [`b66b995d3`](https://github.com/siderolabs/talos/commit/b66b995d34306192cbaa4ef68fe39f821b37d1f0) feat: update default Kubernetes to v1.34.0-rc.1
* [`b967c587d`](https://github.com/siderolabs/talos/commit/b967c587d9f217f25798e0bee0c90393e55dc085) docs: fix clone URL to include `.git`
* [`b72c68398`](https://github.com/siderolabs/talos/commit/b72c6839806103ac0a76acd46f30eabea0375790) docs: edit the insecure, etcd-metrics, inline and extramanifests
* [`e5b9c1fff`](https://github.com/siderolabs/talos/commit/e5b9c1ffffec9fd49ffb84a36c918e75eaa8f1ef) docs: remov RAS Syndrome
* [`701fe774b`](https://github.com/siderolabs/talos/commit/701fe774bd19de7c9f21e043e1520161a8c5fff7) docs: fix cilium links and bump to 1.18.0
* [`d306713a1`](https://github.com/siderolabs/talos/commit/d306713a13a18d7af6caffd5890d54d91d22cad7) feat: update Go to 1.24.6
* [`721595a00`](https://github.com/siderolabs/talos/commit/721595a0009f78a2722802ab665957fd767c4d1e) chore: add deadcode elimination linter
* [`dc4865915`](https://github.com/siderolabs/talos/commit/dc4865915d567942adea3efa66f8ad360f9c4cce) refactor: stop using `text/template` in `machined` code paths
* [`545be55ed`](https://github.com/siderolabs/talos/commit/545be55edc863245638d4387cb9ee7e7b068f2ba) feat: add a pause function to dashboard
* [`06a6c0fe3`](https://github.com/siderolabs/talos/commit/06a6c0fe332940b7a70ea2652bc2a5e7bc51bbf3) refactor: fix deadcode elimination with godbus
* [`2dce8f8d4`](https://github.com/siderolabs/talos/commit/2dce8f8d4693a85d2f3bf46169af8cf502d49f9d) refactor: replace containerd/containerd/v2 module for proper DCE
* [`9b11d8608`](https://github.com/siderolabs/talos/commit/9b11d86081df8cf77860d2d27eed5d8001ff721e) chore: rekres to configure slack notify workflow for CI failures
* [`5ce6a660f`](https://github.com/siderolabs/talos/commit/5ce6a660f67f4e2776550a1e621179beb8a6788c) docs: augment the pod security docs
* [`ada51ff69`](https://github.com/siderolabs/talos/commit/ada51ff696011e15dcd9c661da1d839bdc341745) fix: unmarshal encryption STATE from META
* [`b9e9b2e07`](https://github.com/siderolabs/talos/commit/b9e9b2e07a645f53ca23355810d485a2622870c9) docs: add what is new notes for 1.11
* [`53055bdf4`](https://github.com/siderolabs/talos/commit/53055bdf49ce4c81f63c159cdbaa8ea85d9ca2b8) docs: fix typo in kubevirt page
* [`8d12db480`](https://github.com/siderolabs/talos/commit/8d12db480c38ec37aee5ae7721b2e5ca55ad733e) fix: one more attempt to fix volume mount race on restart
* [`34d37a268`](https://github.com/siderolabs/talos/commit/34d37a268a9e0098179369af128261dbfc956d1d) chore: rekres to use correct slack channel for slack-notify
* [`326a00538`](https://github.com/siderolabs/talos/commit/326a00538210bf98b01795d314c1e154a74d2d58) feat: implement `talos.config.early` command line arg
* [`a5f3000f2`](https://github.com/siderolabs/talos/commit/a5f3000f2e8a79d4e9a5be95fbcac91a2d78675b) feat: implement encryption locking to STATE
* [`c1e65a342`](https://github.com/siderolabs/talos/commit/c1e65a34256944743e768613b119c0caa517b54d) docs: remove talos API flags from mgmt commands
* [`181d0bbf5`](https://github.com/siderolabs/talos/commit/181d0bbf5381343d35a01190da45e3442320d7c5) feat: bootedentry resource
* [`7ad439ac3`](https://github.com/siderolabs/talos/commit/7ad439ac35859695074d3a3efdcdb5c0cab1a5c6) fix: enforce minimum size on user volumes if not set explicitly
* [`50e37aefd`](https://github.com/siderolabs/talos/commit/50e37aefdbde973bcc8aa352639946490fbe7d94) fix: live reload of TLS client config for discovery client
* [`87efd75ef`](https://github.com/siderolabs/talos/commit/87efd75efb3e62b88b4f65a221f9fbdd4b4d6ef9) feat: update containerd to 2.1.4
* [`724b9de6d`](https://github.com/siderolabs/talos/commit/724b9de6d5195bcccc5f484c696429b2f09ab16e) feat: add F71808E watchdog driver
* [`8af96f7af`](https://github.com/siderolabs/talos/commit/8af96f7afdac1c4d5e2697b897b81e2bddd15f66) docs: add ETCD downgrade documentation
* [`44edd205d`](https://github.com/siderolabs/talos/commit/44edd205d5fdffab39b65ee62695a40e22ef188c) docs: add remark about 'exclude-from-external-load-balancers' label
* [`727101926`](https://github.com/siderolabs/talos/commit/7271019263b0dc5b28d2764d19fe531e473222fc) fix(ci): use a random suffix for ami names
* [`d621ce372`](https://github.com/siderolabs/talos/commit/d621ce3726f20ee568ea3b6ac57d9e8dfa0580cc) fix: grype scan
* [`d62e255c2`](https://github.com/siderolabs/talos/commit/d62e255c260810a5f0f2959e32592a3331df28d3) fix: issues with reading GPT
* [`5d0883e14`](https://github.com/siderolabs/talos/commit/5d0883e147163c12a77cd926db799ffed854aedf) feat: update PCI DB module to v0.3.2
* [`3751c8ccf`](https://github.com/siderolabs/talos/commit/3751c8ccfa1bab9fcd435290f36e9012a5626e40) test: wait for service account test job longer
* [`a592eb9f9`](https://github.com/siderolabs/talos/commit/a592eb9f98788883a7ec6d17772e10707230a0d8) feat: update Linux to 6.12.40
* [`4c40e6d3f`](https://github.com/siderolabs/talos/commit/4c40e6d3fb4c2f451a8d7a671df5f6254161bd5d) feat: update etcd to 3.6.4
* [`2bc37bd2c`](https://github.com/siderolabs/talos/commit/2bc37bd2c9679c8055fd7b52eb310f23a329af4e) docs: fix error in kernel module guide
* [`bfc57fb86`](https://github.com/siderolabs/talos/commit/bfc57fb863224f7626f49e5b26be06f77bea2e40) chore: tag aws snapshots created via ci with the image name
* [`06ef7108a`](https://github.com/siderolabs/talos/commit/06ef7108a6050b3a8fd7535f01a469f09042bf56) fix: issue with volume remount on service restart
* [`03efbff18`](https://github.com/siderolabs/talos/commit/03efbff18e420c4fe960f490f91dd9f4751ece04) docs: add SBOM documentation
* [`af8a2869d`](https://github.com/siderolabs/talos/commit/af8a2869dbbec073ffaf72a1378682e109b053ec) fix: do not download artifacts for cron Grype scan
* [`5f442159b`](https://github.com/siderolabs/talos/commit/5f442159b224c96c90badc7176fed17bfb561709) feat: unify disk encryption configuration
* [`38e176e59`](https://github.com/siderolabs/talos/commit/38e176e594edb3d271d98f78417b9fd5ba0c5288) chore(ci): fix datasource versioning
* [`85d6b9198`](https://github.com/siderolabs/talos/commit/85d6b919890a1aa9c4f94d5b18861cc617134ff9) feat: update etcd to v3.5.22
* [`dd7bd2dab`](https://github.com/siderolabs/talos/commit/dd7bd2dab8cf09334e3e353d6a477509bbaa303e) docs: rewrite the getting started and prod docs for v1.10 and v1.11
* [`136a899aa`](https://github.com/siderolabs/talos/commit/136a899aa25b3fdcdd771594668278d563f09192) chore: regenerate release step with signing fixes
* [`450b30d5a`](https://github.com/siderolabs/talos/commit/450b30d5a986563869efdbaa074e82d612f6f2ef) chore(ci): add more nvidia test matrix
* [`451c2c4c3`](https://github.com/siderolabs/talos/commit/451c2c4c39e70c20df58fc31459cd5c789a0e46f) test: add talosctl:latest to the image cache
</p>
</details>

### Dependency Changes

* **github.com/bougou/go-ipmi**                  v0.7.8 -> v0.8.1
* **github.com/cosi-project/runtime**            v1.12.0 -> v1.13.0
* **github.com/klauspost/compress**              v1.18.1 -> v1.18.3
* **github.com/planetscale/vtprotobuf**          79df5c4772f2 -> ba97887b0a25
* **github.com/siderolabs/image-factory**        v0.8.4 -> v0.9.0
* **github.com/siderolabs/omni/client**          v1.3.2 -> v1.4.7
* **github.com/siderolabs/talos**                v1.11.5 -> v1.12.2
* **github.com/siderolabs/talos/pkg/machinery**  v1.12.0-beta.0 -> v1.13.0-alpha.0
* **github.com/spf13/cobra**                     v1.10.1 -> v1.10.2
* **go.uber.org/zap**                            v1.27.0 -> v1.27.1
* **golang.org/x/net**                           v0.47.0 -> v0.49.0
* **golang.org/x/sync**                          v0.18.0 -> v0.19.0
* **google.golang.org/grpc**                     v1.76.0 -> v1.78.0
* **google.golang.org/protobuf**                 v1.36.10 -> v1.36.11

Previous release can be found at [v0.7.1](https://github.com/siderolabs/omni-infra-provider-bare-metal/releases/tag/v0.7.1)

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

