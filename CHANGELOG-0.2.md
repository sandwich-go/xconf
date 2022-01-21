### Unreleased (2022-01-22 00:55:41)

#### 🤖  Tools
  * **sem**: make changelog ([47f9729](https://github.com/sandwich-go/xconf/commit/47f9729cc635287cca5d7085f8626801699064a7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-22 00:55:30 &#43;0800 &#43;0800</small>)

### v0.2.6 🌈 (2022-01-22 00:29:23)

#### 🐛  Bug Fixed
  * rm debug info ([6a1bab0](https://github.com/sandwich-go/xconf/commit/6a1bab0332e2febcf740c198db9007ccdb987785)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 20:06:57 &#43;0800 &#43;0800</small>)

#### 🚀  New Feature
  * xcmd support ([accc582](https://github.com/sandwich-go/xconf/commit/accc582fe6aecee2e5f8997df5dd7526b8653706)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-18 18:17:11 &#43;0800 &#43;0800</small>)

#### 🛠  Refactor
  * xcmd Usage ([ae45d4f](https://github.com/sandwich-go/xconf/commit/ae45d4f359b27a6f871879b73276bc2bf95f1ea5)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-22 00:29:23 &#43;0800 &#43;0800</small>)
  * xcmd ([6cafeea](https://github.com/sandwich-go/xconf/commit/6cafeeacd14d71b57b6661755b18f6f2bade8dc4) , [2c3879d](https://github.com/sandwich-go/xconf/commit/2c3879d33dacc211f6482099118a7f9cb0907a9b)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-22 00:03:03 &#43;0800 &#43;0800</small>)
  * add config interface hooker ([f00e718](https://github.com/sandwich-go/xconf/commit/f00e7187c7ee66ecd0f8f93d0fbd2f61e4d2c470)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 20:36:12 &#43;0800 &#43;0800</small>)
  * add ParseMetaKeyFlagFiles ([f653648](https://github.com/sandwich-go/xconf/commit/f653648e02cffe0d5f652c3e2bc546f8bcf0cee2)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 20:14:14 &#43;0800 &#43;0800</small>)
  * xflag ignore check ([4010f05](https://github.com/sandwich-go/xconf/commit/4010f050bd2e6321a70494052cde306b407ad588)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 19:12:44 &#43;0800 &#43;0800</small>)
  * add readme ([1558802](https://github.com/sandwich-go/xconf/commit/155880240863a860aa72dd1fe2022a6f3686ac66)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 18:09:52 &#43;0800 &#43;0800</small>)
  * rm depend math package ([d4a0f9c](https://github.com/sandwich-go/xconf/commit/d4a0f9cf6e7da1a8ddaab41ec086c98e753acc1c)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 14:17:49 &#43;0800 &#43;0800</small>)
  * print as - if flag not define by xonf ([0619ac1](https://github.com/sandwich-go/xconf/commit/0619ac19f551cd0f57144a124efb7fedae6c9541)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 14:10:44 &#43;0800 &#43;0800</small>)
  * add comment ([5d0a559](https://github.com/sandwich-go/xconf/commit/5d0a559cf60e650e9910990e0c7f99b853632049)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 14:04:22 &#43;0800 &#43;0800</small>)
  * support string value alias for flag and env ([63b0e10](https://github.com/sandwich-go/xconf/commit/63b0e10b3bb403c5cfafa40443404915022697b0)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 13:22:30 &#43;0800 &#43;0800</small>)
  * support string value alias,but flag env not support now ([65bd505](https://github.com/sandwich-go/xconf/commit/65bd505d34301d443c0650a5dba19f2df1d9072b)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 12:32:21 &#43;0800 &#43;0800</small>)
  * support string value alias ([d4b089d](https://github.com/sandwich-go/xconf/commit/d4b089de4df4bebee799142c5a169ff048e39597)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-20 12:27:33 &#43;0800 &#43;0800</small>)
  * do not return error help ([6362353](https://github.com/sandwich-go/xconf/commit/6362353ea054de95ce298adde7f903c9f98e0b48)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-19 18:49:59 &#43;0800 &#43;0800</small>)
  * xcmd comments ([39d1523](https://github.com/sandwich-go/xconf/commit/39d152372a9ab6478abafabd2183eaefbc3a958a)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-19 16:43:13 &#43;0800 &#43;0800</small>)
  * xcmd support pre middleware ([2835e56](https://github.com/sandwich-go/xconf/commit/2835e56edb55fc533a46ec2d9727fb57c87a86ff)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-19 10:21:59 &#43;0800 &#43;0800</small>)
  * set root command ([7aac31b](https://github.com/sandwich-go/xconf/commit/7aac31bdb615e2decd698d49a9f1b2440b832aec)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-18 18:25:35 &#43;0800 &#43;0800</small>)
  * using xconf.Usage as default error output ([e4e5aa4](https://github.com/sandwich-go/xconf/commit/e4e5aa45e8122cfe15676b84ce38347738f25c97)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 17:54:37 &#43;0800 &#43;0800</small>)
  * support print xconf usage ([03eb8de](https://github.com/sandwich-go/xconf/commit/03eb8def087378f3d68cc029a01e38f8edf8381c)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 17:20:38 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([5f659a7](https://github.com/sandwich-go/xconf/commit/5f659a72d95fa61a6fedc66ebeed695bba86175f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 17:17:53 &#43;0800 &#43;0800</small>)

### v0.2.5 (2022-01-13 17:16:47)

#### 🛠  Refactor
  * support print xconf usage ([4b12d7f](https://github.com/sandwich-go/xconf/commit/4b12d7f53ea6c9e8aeae3c3ad3ad08ddb40bb59d)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 17:16:47 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([526256d](https://github.com/sandwich-go/xconf/commit/526256d1f6760df647d786ce383eea509f241663)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 13:59:27 &#43;0800 &#43;0800</small>)

#### 💪  Commit
  * Merge branch 'version/0.2' ([9026d0e](https://github.com/sandwich-go/xconf/commit/9026d0e9ed7f77044066f8cd2d1bd516e69a154b)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 13:59:47 &#43;0800 &#43;0800</small>)
  * go lint ([6ee6cb9](https://github.com/sandwich-go/xconf/commit/6ee6cb9e8c1421d42919bb0bb627c854472dc5b1)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 13:18:00 &#43;0800 &#43;0800</small>)

### v0.2.4 (2022-01-13 13:06:47)

#### 🛠  Refactor
  * usage using self UnquoteUsage ([4fdcbeb](https://github.com/sandwich-go/xconf/commit/4fdcbebc79638e1869e4e33bced175fee8e60d41)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 12:53:25 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([81b55a1](https://github.com/sandwich-go/xconf/commit/81b55a1de88c08e5b39db5269e7ffed00fc097ab) , [68bbaf1](https://github.com/sandwich-go/xconf/commit/68bbaf1351918206d4e10ba02c7eca3849cd721a)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 12:38:56 &#43;0800 &#43;0800</small>)

#### 💪  Commit
  * update readme ([1615b2a](https://github.com/sandwich-go/xconf/commit/1615b2abac650b057ffc0db66204ae3f0818b3bf)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 13:06:47 &#43;0800 &#43;0800</small>)
  * Merge branch 'version/0.3' ([398fea3](https://github.com/sandwich-go/xconf/commit/398fea366c4edc3431af00eaf5850903ab658060)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 13:03:32 &#43;0800 &#43;0800</small>)

### v0.2.3 (2022-01-13 12:37:48)

#### 🐛  Bug Fixed
  * parse arg even if args=0, we need create flag info into the flagset ([359ec07](https://github.com/sandwich-go/xconf/commit/359ec073aab43d4f3be5c71e66c91ec7fcf569ff) , [e98ea2c](https://github.com/sandwich-go/xconf/commit/e98ea2cefb8186888e87bf501cc853f50c99ae30)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 12:37:48 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([1abf9af](https://github.com/sandwich-go/xconf/commit/1abf9afa77b99ed6e215a494d8dc0ef9259dc7b4)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 12:27:58 &#43;0800 &#43;0800</small>)

### v0.2.2 (2022-01-13 12:25:27)

#### 🐛  Bug Fixed
  * global Usage support ([e6c8e8b](https://github.com/sandwich-go/xconf/commit/e6c8e8b65026dd9d6a55fdc8fcf3456ed9a8b53f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 11:03:12 &#43;0800 &#43;0800</small>)

#### 🛠  Refactor
  * usage support Deprecated info ([b227595](https://github.com/sandwich-go/xconf/commit/b2275951558d92126351e5346d05fad88b6ea373)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-13 12:25:27 &#43;0800 &#43;0800</small>)

### v0.2.1 (2022-01-12 20:59:15)

#### 🐛  Bug Fixed
  * Usage should  not support passed in,using new FlagSet ([b8af8e8](https://github.com/sandwich-go/xconf/commit/b8af8e835ef7bbc5433c1d76fa95ea12d64ed2c4)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 20:59:15 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([d2b5228](https://github.com/sandwich-go/xconf/commit/d2b522888077fe896eeace1966a895696530262e)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 20:19:17 &#43;0800 &#43;0800</small>)

### v0.2.0 (2022-01-12 20:12:16)

#### 🐛  Bug Fixed
  * usage print option global usage info ([a84e1b2](https://github.com/sandwich-go/xconf/commit/a84e1b293145ae9530bf49a6ec4ba0ab31291139)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 19:11:42 &#43;0800 &#43;0800</small>)
  * get char width using runewidth package ([8d754d0](https://github.com/sandwich-go/xconf/commit/8d754d056376f1e089a66900e392f7a19ad464c7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 18:42:41 &#43;0800 &#43;0800</small>)
  * usage print ([5465e6a](https://github.com/sandwich-go/xconf/commit/5465e6aa1c12b37c42029bf34a0b96f5fe4d9b37)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 17:49:09 &#43;0800 &#43;0800</small>)
  * check flagset nil ([94a3d77](https://github.com/sandwich-go/xconf/commit/94a3d770cf5f45fa287eb6521aeba2d13e68e7f5)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 13:40:33 &#43;0800 &#43;0800</small>)
  * go lint ([ef1c0ad](https://github.com/sandwich-go/xconf/commit/ef1c0ad0e54b0cb171a9458b4161ae7f35ca1dc7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 14:47:12 &#43;0800 &#43;0800</small>)
  * samole for replit ([cdbc628](https://github.com/sandwich-go/xconf/commit/cdbc628060fcb7a40962e3f6ab38492bd2e1e4b6)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 14:42:20 &#43;0800 &#43;0800</small>)
  * type alias, do not get Underlying type name for Anonymous field ([3a1cc4e](https://github.com/sandwich-go/xconf/commit/3a1cc4e7e9686e1df5727f8e30c25b0196fb191d)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 13:57:47 &#43;0800 &#43;0800</small>)
  * check err ([f21d663](https://github.com/sandwich-go/xconf/commit/f21d6639ea39a9addf5eea81cb9624c4fc3d09c7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 12:21:18 &#43;0800 &#43;0800</small>)
  * flag set parse ([b423db9](https://github.com/sandwich-go/xconf/commit/b423db99075309b64d4a4d223937ab6fea9e1511)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 12:16:14 &#43;0800 &#43;0800</small>)
  * ci support go 1.14  using ioutil.WriteFile ([03ea69b](https://github.com/sandwich-go/xconf/commit/03ea69beb913a9ffb4722f2b28daaaeb54e0f9ab)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 16:14:53 &#43;0800 &#43;0800</small>)
  * ci support go 1.15 1.16 ([cb12294](https://github.com/sandwich-go/xconf/commit/cb122943bfe26a777cc85fa1cde6e879b731191c) , [e1b25e2](https://github.com/sandwich-go/xconf/commit/e1b25e2fcb653118622a7414d8fe08fff50eea01)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 16:12:54 &#43;0800 &#43;0800</small>)
  * lint ([7ee0ab4](https://github.com/sandwich-go/xconf/commit/7ee0ab4affff26d0376a5162a58ae13ac65e6466)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 16:02:14 &#43;0800 &#43;0800</small>)
  * rm io.ReadAll ([48e19bd](https://github.com/sandwich-go/xconf/commit/48e19bd76f918f05155a3dc43cc438c5e85613dd)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 15:50:25 &#43;0800 &#43;0800</small>)
  * watcher logic ([ae1416c](https://github.com/sandwich-go/xconf/commit/ae1416c3dbe3d44effb5b1cc798f346a718e2b0d)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 15:45:11 &#43;0800 &#43;0800</small>)
  * rm unused code ([987bbf2](https://github.com/sandwich-go/xconf/commit/987bbf2d7f2b064578534f2663f369b52f252b9f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 14:20:40 &#43;0800 &#43;0800</small>)
  * race when test and rm ioutil ([6aa98e3](https://github.com/sandwich-go/xconf/commit/6aa98e342475e88c5cc8786a7d6095ce5ba2edab)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 14:14:17 &#43;0800 &#43;0800</small>)
  * race when test ([91775a0](https://github.com/sandwich-go/xconf/commit/91775a08e43a7eb0c8ffbf7ee2f5b898f4e71f18)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 14:11:37 &#43;0800 &#43;0800</small>)
  * panic when no tag ([fb843af](https://github.com/sandwich-go/xconf/commit/fb843afc924ec2ab2c4695582fe7552e06acee63)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 13:05:00 &#43;0800 &#43;0800</small>)

#### 🚀  New Feature
  * --help --h support print xconf Usage ([130a387](https://github.com/sandwich-go/xconf/commit/130a387fc1c3fa10fb2ded24c8881db20b0ae23f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 20:11:19 &#43;0800 &#43;0800</small>)

#### 🛠  Refactor
  * replit sample ([bcbbd14](https://github.com/sandwich-go/xconf/commit/bcbbd14e7fe134eea19348e3414f6a35ac1116d2)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 12:51:03 &#43;0800 &#43;0800</small>)
  * check if leaf node ([8610eb4](https://github.com/sandwich-go/xconf/commit/8610eb454009338f49eb6b80eee78f0921b71c3a) , [724f7de](https://github.com/sandwich-go/xconf/commit/724f7def8643e0f7974cdbd2e26dfd0416312512)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 12:45:51 &#43;0800 &#43;0800</small>)
  * using latest optiongen ([66cbd4f](https://github.com/sandwich-go/xconf/commit/66cbd4f45b83e2f271185b2065a40e52f1ed44fd) , [b12dd42](https://github.com/sandwich-go/xconf/commit/b12dd42ed7010189704d1279edff47943fa24fae) , [aaf3347](https://github.com/sandwich-go/xconf/commit/aaf33473d89db59fa68628667f04a95f2fb3ebc8)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 18:20:29 &#43;0800 &#43;0800</small>)
  * simple annotation support ([9b88991](https://github.com/sandwich-go/xconf/commit/9b88991b7462de6597fcf264892b0cd77ec92221)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 18:15:14 &#43;0800 &#43;0800</small>)
  * filter empty files ([c47c8f5](https://github.com/sandwich-go/xconf/commit/c47c8f55cf40a47982559c801d647e63b9f42f3f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 12:34:38 &#43;0800 &#43;0800</small>)
  * gen comments ([f7e1d33](https://github.com/sandwich-go/xconf/commit/f7e1d334457425cd12909c7f60d5effde55872d9) , [22c8755](https://github.com/sandwich-go/xconf/commit/22c87558fad95257189cd765b8082438708273a6)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 12:28:01 &#43;0800 &#43;0800</small>)
  * keep meta data into datMeta ([ce0a74d](https://github.com/sandwich-go/xconf/commit/ce0a74d7f89523f1b20e254660e6d732d82a2771) , [fba98a0](https://github.com/sandwich-go/xconf/commit/fba98a0bf2e3af86c1594f2458249bd675a67b83)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 17:49:14 &#43;0800 &#43;0800</small>)
  * install usage for default FlagSet ([3438a21](https://github.com/sandwich-go/xconf/commit/3438a21355d79e673c7fb29455362bae74c6bba8) , [beaf5ff](https://github.com/sandwich-go/xconf/commit/beaf5ff06a8c268a4af9bebf7c9b61b636cb1fea)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 17:32:48 &#43;0800 &#43;0800</small>)
  * usage change ([92191e2](https://github.com/sandwich-go/xconf/commit/92191e2b59d4459bdcf10ed358fcd97f695aaba2)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 17:06:41 &#43;0800 &#43;0800</small>)
  * update doc ([2d5761e](https://github.com/sandwich-go/xconf/commit/2d5761ec57371cf0395516e33f0f7b41934215c2)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 15:14:51 &#43;0800 &#43;0800</small>)
  * rm print log ([f29fa34](https://github.com/sandwich-go/xconf/commit/f29fa34552c8c65f5d43d4506c3e298995c5fe20)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 15:13:43 &#43;0800 &#43;0800</small>)
  * add simple gray rule support ([b507e59](https://github.com/sandwich-go/xconf/commit/b507e59e6640f32262dc3b9a9412ab6946c851c1)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 15:13:19 &#43;0800 &#43;0800</small>)
  * rm auto options ([a691e72](https://github.com/sandwich-go/xconf/commit/a691e721dda5af7ff3316aa3a56bebab05281fb6)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 13:34:18 &#43;0800 &#43;0800</small>)
  * add FieldFlagSetCreateIgnore ([7eb9a8a](https://github.com/sandwich-go/xconf/commit/7eb9a8af435d31dd616ce298bfb662239b9110ab)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 13:15:43 &#43;0800 &#43;0800</small>)
  * add magic codec ([3574199](https://github.com/sandwich-go/xconf/commit/3574199c5aa9b94a6dc1f004683c7d788ebeacaf)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 13:08:42 &#43;0800 &#43;0800</small>)
  * comment ([94ee300](https://github.com/sandwich-go/xconf/commit/94ee3002c818a11fdb5ce6c4cb35288d4853f8f9)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 12:54:07 &#43;0800 &#43;0800</small>)
  * chain move into codec ([d24cbda](https://github.com/sandwich-go/xconf/commit/d24cbda2ae2cce3984ac05bf82243ddd61b6cd29)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 12:53:26 &#43;0800 &#43;0800</small>)
  * secconf reader ([92399be](https://github.com/sandwich-go/xconf/commit/92399be3e1c0d09b1d304dcce17392b9938833fd)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 12:52:43 &#43;0800 &#43;0800</small>)
  * secconf chain support ([fde9667](https://github.com/sandwich-go/xconf/commit/fde96673d0ebb5318a1e1efdd2110241b35192b8) , [3255bf8](https://github.com/sandwich-go/xconf/commit/3255bf8a2464d967a5139472e36b42dc6c3d988e)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 12:50:27 &#43;0800 &#43;0800</small>)
  * support sec conf ([7dc32e6](https://github.com/sandwich-go/xconf/commit/7dc32e66662eb6c97299c4a453960f9557d82f2c)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-04 20:47:22 &#43;0800 &#43;0800</small>)
  * auto recover watch ([1079f51](https://github.com/sandwich-go/xconf/commit/1079f51021eddefd3e8f74e749e3208c6cf220a7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-04 15:55:08 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([680bb0e](https://github.com/sandwich-go/xconf/commit/680bb0ea11a0384d8fe3af587473f7c534cd7ed9) , [5d41300](https://github.com/sandwich-go/xconf/commit/5d413002cb192d25a73a22bca78fd8b9c5da71c7) , [dd61867](https://github.com/sandwich-go/xconf/commit/dd61867b7aeeaa42418af5e120dadc1c8dfa944d) , [17c08d3](https://github.com/sandwich-go/xconf/commit/17c08d3ebaaa373f27eb6831635de1f644c5bc80) , [9778e46](https://github.com/sandwich-go/xconf/commit/9778e4691d7ed694e6921aca4c9d43bede12989b) , [d340965](https://github.com/sandwich-go/xconf/commit/d34096507b2f609e0c70fdd98819423204ace5fc) , [b5d03eb](https://github.com/sandwich-go/xconf/commit/b5d03ebb560524f7061bf2021d3b0f56eafaa0d7) , [b408a91](https://github.com/sandwich-go/xconf/commit/b408a918db57df8def0c8ffd663a85532c8e4e92) , [ac8ecc8](https://github.com/sandwich-go/xconf/commit/ac8ecc805f8015ffbbf54f23e1695ed6eb0998f2) , [eea6063](https://github.com/sandwich-go/xconf/commit/eea6063ac5cd785415b6de00d0af45a7e17622c5) , [6ac07f4](https://github.com/sandwich-go/xconf/commit/6ac07f4aa6c20f1fe0f5278ab3b14d3f480e8f32)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 20:12:16 &#43;0800 &#43;0800</small>)

#### 📝  Add Docs
  * add sample ([247d398](https://github.com/sandwich-go/xconf/commit/247d398b59c84287c5c357ca9d40182751de2afd)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 18:09:27 &#43;0800 &#43;0800</small>)
  * fix doc bug ([4d4c9d2](https://github.com/sandwich-go/xconf/commit/4d4c9d2bd9c3fad23e889ffed29203d8f81c493c)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 17:46:42 &#43;0800 &#43;0800</small>)
  * add more comments ([77745a7](https://github.com/sandwich-go/xconf/commit/77745a7515129cca9401da58f5c2812c84850704) , [0e95b10](https://github.com/sandwich-go/xconf/commit/0e95b1039c7b2781be455582237321dd91c017c2)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 13:45:26 &#43;0800 &#43;0800</small>)
  * add doc ([ab12db3](https://github.com/sandwich-go/xconf/commit/ab12db318148d472c949e03b88c385a6912ed8a4)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-05 13:32:30 &#43;0800 &#43;0800</small>)

#### 💪  Commit
  * support usage ([a59c663](https://github.com/sandwich-go/xconf/commit/a59c66373f67b43e7049aafe5d29540011002be1)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 13:06:10 &#43;0800 &#43;0800</small>)
  * using latest option gen ([1a2383a](https://github.com/sandwich-go/xconf/commit/1a2383ae443e33476c7161d86b013080ae87f2af) , [0c3bc65](https://github.com/sandwich-go/xconf/commit/0c3bc65f2d3c336f2c1445989baa916f553f5c6a)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-12 11:55:17 &#43;0800 &#43;0800</small>)
  * use latest optiongen ([91e757e](https://github.com/sandwich-go/xconf/commit/91e757e9cda17b4581204f7ec7e80801efbc6f01)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-11 17:10:36 &#43;0800 &#43;0800</small>)
  * use warning to print parse from default ([cc22cdb](https://github.com/sandwich-go/xconf/commit/cc22cdb2ec5f0856d1048c5178fb3537a58a0999)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-11 17:09:25 &#43;0800 &#43;0800</small>)
  * optiongen ([f0a5500](https://github.com/sandwich-go/xconf/commit/f0a5500e87b233a99d611bb8ce2d06c3d51dcf1f)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-11 12:35:40 &#43;0800 &#43;0800</small>)
  * go-lint ([5f3867a](https://github.com/sandwich-go/xconf/commit/5f3867a32f1255ad10f59f3a81838652a57fc4a0) , [45dd21d](https://github.com/sandwich-go/xconf/commit/45dd21d199b244099785881d1fafadf93f556e02)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-11 12:28:15 &#43;0800 &#43;0800</small>)
  * Merge branch 'main' of github.com:sandwich-go/xconf ([29cf89c](https://github.com/sandwich-go/xconf/commit/29cf89c693d7295d1b1de4eaadb4a3521fb5bfdf) , [6f54c82](https://github.com/sandwich-go/xconf/commit/6f54c82b3b8779319a1be3f1326e7974618e5e12) , [82bd895](https://github.com/sandwich-go/xconf/commit/82bd895080990e32477ce950532a9ef743f72856) , [2ce163a](https://github.com/sandwich-go/xconf/commit/2ce163a5709124413621dfed7c43d4e4017eb90c) , [63b8b17](https://github.com/sandwich-go/xconf/commit/63b8b17545ba3957fea7cd13b2a0e9716bab8b1b)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-11 12:22:07 &#43;0800 &#43;0800</small>)
  * Update README.md ([f962cd5](https://github.com/sandwich-go/xconf/commit/f962cd5c57f255db8ebb1b0826036331c24d7a38) , [a043ee7](https://github.com/sandwich-go/xconf/commit/a043ee72dd9833c0ae766d46d0e5d6ad238124a7) , [c0b647e](https://github.com/sandwich-go/xconf/commit/c0b647e1b43ccc5771e0a32d1181cee8d877d7e5) , [ad6e92f](https://github.com/sandwich-go/xconf/commit/ad6e92f268cc7e968bd3642139c30c9eefa98caa) , [356bcdb](https://github.com/sandwich-go/xconf/commit/356bcdb87b38a8a343a760feeeb80417d20fd357)) (<small>[timestee](19310233&#43;timestee@users.noreply.github.com)@2022-01-11 11:23:44 &#43;0800 &#43;0800</small>)
  * using latest optiongen ([11db871](https://github.com/sandwich-go/xconf/commit/11db8714f4cabdea14b6b3d7852fe1e7a02a3367)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 21:02:49 &#43;0800 &#43;0800</small>)
  * add git ignore ([14d380e](https://github.com/sandwich-go/xconf/commit/14d380ebda76b704a770f47f5ba78f2af2aa0802)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 20:39:16 &#43;0800 &#43;0800</small>)
  * add etcd filesystem support ([fb79caf](https://github.com/sandwich-go/xconf/commit/fb79caf006300ddd9f5af32748ece2a4ad2a5ee7)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 19:07:23 &#43;0800 &#43;0800</small>)
  * move providers to xconf-providers ([ca0027f](https://github.com/sandwich-go/xconf/commit/ca0027f736f52ba2252b488bdba78b6d36dc0dff)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 19:04:40 &#43;0800 &#43;0800</small>)
  * rm travissupport ([b92b505](https://github.com/sandwich-go/xconf/commit/b92b505284cf4b2a1517b778c8aaeabb9ceb51aa)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 13:48:02 &#43;0800 &#43;0800</small>)
  * replit sample ([02ad89f](https://github.com/sandwich-go/xconf/commit/02ad89fa10e0c1d7e29755dbb0c2f41719fc6bcb)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-10 11:13:29 &#43;0800 &#43;0800</small>)
  * fix : xutil package ([42583c3](https://github.com/sandwich-go/xconf/commit/42583c3f6b90fb19930feee15035c5d44d789a86) , [d0a93f2](https://github.com/sandwich-go/xconf/commit/d0a93f2fee9248ddf410b2723805396885c4386a)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-08 01:10:17 &#43;0800 &#43;0800</small>)
  * fix : just print warning when flagset is ComandLine ([5dc30b9](https://github.com/sandwich-go/xconf/commit/5dc30b913ae30937303a60537104eb8e0a8c7134)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-07 19:55:48 &#43;0800 &#43;0800</small>)
  * not support 15 16 tmp ([4fa5001](https://github.com/sandwich-go/xconf/commit/4fa5001b4cedba656de64ab0847487070f5341a1)) (<small>[timestee](19310233&#43;timestee@users.noreply.github.com)@2022-01-06 14:25:11 &#43;0800 &#43;0800</small>)
  * on branch change to main and verison/* ([18597f4](https://github.com/sandwich-go/xconf/commit/18597f47e94fc72398833578500824b5b7598a47)) (<small>[timestee](19310233&#43;timestee@users.noreply.github.com)@2022-01-06 14:04:44 &#43;0800 &#43;0800</small>)
  * add ci ([901009c](https://github.com/sandwich-go/xconf/commit/901009c3dfab9541b748590da72eacc42a80e631)) (<small>[timestee](19310233&#43;timestee@users.noreply.github.com)@2022-01-06 13:58:11 &#43;0800 &#43;0800</small>)
  * init ci.yml ([d7e312d](https://github.com/sandwich-go/xconf/commit/d7e312d0f3e5d1d6c46e6bc56ea13d07472add04)) (<small>[timestee](19310233&#43;timestee@users.noreply.github.com)@2022-01-06 13:52:53 &#43;0800 &#43;0800</small>)
  * add comments ([b9c18bf](https://github.com/sandwich-go/xconf/commit/b9c18bf6bf4cdc507d71bc155a06a032dc2f7720) , [62c943e](https://github.com/sandwich-go/xconf/commit/62c943e47d79f88d29b73669633a5852a00f79ca)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 11:35:03 &#43;0800 &#43;0800</small>)
  * add .travis and .sembumprc ([4baee75](https://github.com/sandwich-go/xconf/commit/4baee751939030f8be9445937cdca73ae67edbd9) , [ba90289](https://github.com/sandwich-go/xconf/commit/ba90289fbfb4a7651c76b4a38d69e61b92a8836b) , [0d495db](https://github.com/sandwich-go/xconf/commit/0d495db3c4adf74fdee5481cd1654e747c14a571)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-06 11:01:40 &#43;0800 &#43;0800</small>)
  * FieldFlagSetCreateIgnore ([a55c1df](https://github.com/sandwich-go/xconf/commit/a55c1dfafd936060ab302f04f20a275ff9e69999)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-04 14:59:10 &#43;0800 &#43;0800</small>)
  * deprecated log ([034d740](https://github.com/sandwich-go/xconf/commit/034d7404d9c51bdab6580bf82fd67e7ce4e903ab)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-04 14:51:13 &#43;0800 &#43;0800</small>)
  * add xconf_inherit_files xconf_files doc ([c05b240](https://github.com/sandwich-go/xconf/commit/c05b24052a966905f16c39186a533f5e0e2b55ab)) (<small>[hui.wang](hui.wang@funplus.com)@2022-01-04 13:53:24 &#43;0800 &#43;0800</small>)



