# OpenWrtをELECOM WAB-I1750-PSで動かす
~~https://github.com/openwrt/openwrt/pull/3661 にてzpcさんにより、OpenWrt公式へのELECOM WAB-I1750-PS向けのFWがビルドできるPRが出ていますが、上流へのPRは中止したそうです( https://zpc.st/hardware/elecom/wab-i1750-ps/ にzpcさんが解析した内容が載っています )

どうやらこれで動くようなので、このコードを最新版にパッチしてビルドします~~


[zpcさん・大破さんによって作成されたWAB-I1750-PS上でOpenWrtが動くようになるパッチ](https://github.com/openwrt/openwrt/pull/14706)がOpenWrt公式へマージされました。

まだ正式版としてはリリースされていないため、openWrt公式からはLuCIが入っていないsnapshotビルドの公開しか行われておりません。LuCI入りのビルド済みバイナリをReleasesに置いておくので**下のお約束を読んだ上で**ダウンロードしてご自由にご利用ください


## お約束

- **日本国内の場合、法律の関係上OpenWrtで無線LANを飛ばす行為は非推奨です。もし吹く場合でも、調べてから自己責任で飛ばしてください**
- 以下の方法を参考にビルドしたものを用いて / 配布するビルド済みのバイナリを用いて OpenWrt化して何らかの損害が出た場合でも、責任は一切取れません。**自己責任でお願いします**

## 参考情報

### 元のFWに戻す方法

大破さんにより、OpenWrtを入れたWAB-I1750-PSを戻す方法が[コミットメッセージ](https://github.com/openwrt/openwrt/commit/b18edb1bfa34420fde1404d9d1e619c889557154)に記されています

> Revert to OEM firmware:
> 
> 1. Download the latest OEM firmware
> 2. Remove 128 bytes(0x80) header from firmware image
> 3. Decode by xor with a pattern "8844a2d168b45a2d" (hex val)
> 4. Upload the decoded firmware to the device
> 5. Flash to "firmware" partition by mtd command
> 6. Reboot

ELECOMから提供されたFWはそのままは適用できず、細工をする必要があります。2.と3.を勝手にやるプログラムを作成したので、それをReleasesから拾ってきてやります

TODO

### ビルド方法
公式に入ったので、下記の内容はほとんど[[OpenWrt Wiki] Quick image building guide](https://openwrt.org/docs/guide-developer/toolchain/beginners-build-guide) と同じものです

ここではUbuntu22.04を前提としています(WSLは公式にサポートされていません : [詳細](https://openwrt.org/docs/guide-developer/toolchain/wsl))

#### GitHubから最新Stableを拾ってくる
```
git clone -b v23.05.2 https://github.com/openwrt/openwrt.git
cd openwrt/
```

#### feedsをインストール
```
./scripts/feeds update -a
./scripts/feeds install -a
```

#### makeの前準備
```
make manuconfig
```

各種チェックが終わった後TUIが開きます
そこで `Target System` を `Atheros ATH79` に `Target Profile` を `ELECOM WAB-I1750-PS` にして、WAB-I1750-PS用のFWが出てくるようにします
OpenWrtを管理するWeb-GUIアプリケーションのLuCIを入れるので、 `LuCI` → `1.Collections` → `luci` を `*` にします(`M`だと、インストール時にインターネットから拾う設定になってしまいます)
後は日本語化をするために `2.Modules` → `Translations` → `Japanese (ja)`を入れたり入れなかったり、他のプロトコルをサポートするようにしたり、カスタマイズしましょう(パッケージのインストールはインストール後でもできるけどね, **ROMが8MBしかないので注意!**)

#### make!!!!!!!!!!!!!!
**make -j{並列処理させたいコア数} で並列処理が可能です**
結構時間かかるので、気長に待ちます
```
make
```

makeし終わったら、 `bin/targets/ath79/generic/` 配下に色々バイナリが入っています
そのうち末尾が `factory.bin` となっているもの(私の環境では `openwrt-ath79-generic-elecom_wab-i1750-ps-squashfs-factory.bin` )をファームウェアアップデート画面より投入すれば、OpenWrt化ができます!

