# OpenWrtをELECOM WAB-I1750-PSで動かす

https://github.com/openwrt/openwrt/pull/3661 にてzpcさんにより、OpenWrt公式へのELECOM WAB-I1750-PS向けのFWがビルドできるPRが出ていますが、上流へのPRは中止したそうです( https://zpc.st/hardware/elecom/wab-i1750-ps/ にzpcさんが解析した内容が載っています )

どうやらこれで動くようなので、このコードを最新版にパッチしてビルドします

また、ビルド済みの物をReleasesにおいて置くので、**下のお約束を読んだ上で**ダウンロードしてください

**OpenWrtをビルドしたことがない人はこの一生に二度とないOpenWrtのビルドを体験するといいかもしれません**


## お約束

- 以下の方法を参考にビルドしたものを用いて / 配布するビルド済みのバイナリを用いて OpenWrt化して何らかの損害が出た場合でも、責任は一切取れません。**自己責任でお願いします**
- **日本国内の場合、法律の関係上OpenWrtで無線LANを飛ばす行為は非推奨です。もし吹く場合でも、調べてから自己責任で飛ばしてください**
- **私はOpenWrtについて全然詳しくありません。** もし、誤っている部分等ありましたらPRくれるとうれしいです

## 前提
Linux環境が必要です(WSLは公式にサポートされていません : [詳細](https://openwrt.org/docs/guide-developer/toolchain/wsl))

[[OpenWrt Wiki] Build system setup](https://openwrt.org/docs/guide-developer/toolchain/install-buildsystem) に従って各ディストリビューションでビルド用パッケージのインストールをします

Ubuntu 22.04だと以下のコマンドを実行します

```bash
sudo apt update
sudo apt install build-essential clang flex bison g++ gawk \
gcc-multilib g++-multilib gettext git libncurses-dev libssl-dev \
python3-distutils rsync unzip zlib1g-dev file wget
```

## 作業

[[OpenWrt Wiki] Quick image building guide](https://openwrt.org/docs/guide-developer/toolchain/beginners-build-guide) を参考にしつつビルドします

```bash
# GitHubから最新Stableを拾ってくる
git clone -b v23.05.2 https://github.com/openwrt/openwrt.git
cd openwrt/
# PRのpatchを拾ってきます
wget https://patch-diff.githubusercontent.com/raw/openwrt/openwrt/pull/3661.patch
# 適応します(-p0だと動かないので-p1、正直なんでこうなのかわかっていない)
patch -p1 < 3661.patch

# **4年前のPRなので、一部Diffが解決できません**

# 私の場合、`target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/10-ath9k-eeprom`と`target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/11-ath10k-caldata`が失敗しました
# patching file target/linux/ath79/dts/qca9558_devolo_dvl1xxx.dtsi
# patching file target/linux/ath79/dts/qca9558_edimax_wap1750.dtsi
# patching file target/linux/ath79/dts/qca955x.dtsi
# patching file target/linux/ath79/patches-5.10/300-MIPS-ath79-add-missing-QCA955x-UART1-registers.patch
# patching file target/linux/ath79/patches-5.10/301-MIPS-ath79-export-QCA955X-UART1-reference-clock.patch
# patching file target/linux/ath79/patches-5.4/300-MIPS-ath79-add-missing-QCA955x-UART1-registers.patch
# patching file target/linux/ath79/patches-5.4/301-MIPS-ath79-export-QCA955X-UART1-reference-clock.patch
# patching file target/linux/ath79/dts/qca9558_elecom_wab-i1750-ps.dts
# patching file target/linux/ath79/generic/base-files/etc/board.d/02_network
# Hunk #1 succeeded at 139 with fuzz 1 (offset 30 lines).
# Hunk #2 succeeded at 675 with fuzz 1 (offset 114 lines).
# patching file target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/10-ath9k-eeprom
# Hunk #1 FAILED at 38.
# Hunk #2 succeeded at 45 with fuzz 2 (offset -8 lines).
# 1 out of 2 hunks FAILED -- saving rejects to file target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/10-ath9k-eeprom.rej
# patching file target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/11-ath10k-caldata
# Hunk #1 FAILED at 67.
# Hunk #2 succeeded at 73 with fuzz 2 (offset -10 lines).
# 1 out of 2 hunks FAILED -- saving rejects to file target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/11-ath10k-caldata.rej
# patching file target/linux/ath79/image/generic.mk
# Hunk #1 succeeded at 1193 with fuzz 2 (offset 206 lines).

# なので、.rejのファイル/GitHubのPR等々を見ながら人間Patchします

vim target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/10-ath9k-eeprom
# 私の場合、48行目に`elecom,wab-i1750-ps|\`を加えました 
#         nec,wf1200cr|\
#         nec,wg1200cr|\
#         wd,mynet-n600|\
#         wd,mynet-n750)
#                 caldata_extract "art" 0x1000 0x440
#                 ath9k_patch_mac $(mtd_get_mac_ascii devdata "wlan24mac")
#                 ;;  
#         engenius,ecb1200|\
#         engenius,ecb1750)
#                 caldata_extract "art" 0x1000 0x440
#                 ath9k_patch_mac $(macaddr_add $(mtd_get_mac_ascii u-boot-env athaddr) 1)
#                 ;;  
# +++     elecom,wab-i1750-ps|\
#         engenius,ecb1200|\
#         engenius,ecb1750)
#                 caldata_extract "art" 0x1000 0x440
#                 ath9k_patch_mac $(macaddr_add $(mtd_get_mac_ascii u-boot-env athaddr) 1)
#                 ;;  
#         enterasys,ws-ap3705i)
#                 caldata_extract "calibrate" 0x1000 0x440
#                 ath9k_patch_mac $(mtd_get_mac_ascii u-boot-env0 RADIOADDR1)
#                 ;;  

vim target/linux/ath79/generic/base-files/etc/hotplug.d/firmware/11-ath10k-caldata
# 私の場合、76行目に`elecom,wab-i1750-ps|\`を加えました
#         devolo,dvl1750i|\
#         devolo,dvl1750x)
#                 caldata_extract "art" 0x5000 0x844
#                 ath10k_patch_mac $(macaddr_add $(mtd_get_mac_binary art 0x0) -1)
#                 ;;
#         engenius,ecb1200|\
#         engenius,ecb1750)
#                 caldata_extract "art" 0x5000 0x844
#                 ath10k_patch_mac $(mtd_get_mac_ascii u-boot-env athaddr)
#                 ;;
# +++     elecom,wab-i1750-ps|\
#         elecom,wrc-1750ghbk2-i)
#                 caldata_extract "art" 0x5000 0x844
#                 ;;
#         engenius,ecb1200|\
#         engenius,ecb1750)
#                 caldata_extract "art" 0x5000 0x844
#                 ath10k_patch_mac $(mtd_get_mac_ascii u-boot-env athaddr)
#                 ;;
#         engenius,ews511ap)
#                 caldata_extract "art" 0x5000 0x844
#                 ath10k_patch_mac $(macaddr_add $(cat /sys/class/net/eth0/address) 1)
#                 ;;
#         extreme-networks,ws-ap3805i)
#                 caldata_extract "art" 0x5000 0x844
#                 ath10k_patch_mac $(mtd_get_mac_ascii cfg1 RADIOADDR0)
#                 ;;

# 超適当にやっているので、ELECOM WAB-I1750-PS のFWはビルドできますが、ほかのハードウェアのFWをビルドするのにはこのコードを再利用しないほうがいいでしょう

# feedsをインストール
./scripts/feeds update -a
./scripts/feeds install -a

# makeの前準備
make manuconfig
# 各種チェックが終わった後TUIが開きます
# そこで `Target System` を `Atheros ATH79` に `Target Profile` を `ELECOM WAB-I1750-PS` にして、WAB-I1750-PS用のFWが出てくるようにします
# OpenWrtを管理するWeb-GUIアプリケーションのLuCIを入れるので、 `LuCI` → `1.Collections` → `luci` を `*` にします(`M`だと、インストール時にインターネットから拾う設定になってしまいます)
# 後は日本語化をするために `2.Modules` → `Translations` → `Japanese (ja)`を入れたり入れなかったり、他のプロトコルをサポートするようにしたり、カスタマイズしましょう(パッケージのインストールはインストール後でもできるけどね, **ROMが8MBしかないので注意!**)

# make!!!!!!!!!!!!!!
# **make -j{並列処理させたいコア数} で並列処理が可能です**
# 結構時間かかるので、気長に待ちます
make

# makeし終わったら、 `bin/targets/ath79/generic/` 配下に色々バイナリが入っています
# そのうち末尾が `factory.bin` となっているもの(私の環境では `openwrt-ath79-generic-elecom_wab-i1750-ps-squashfs-factory.bin` )をファームウェアアップデート画面より投入すれば、OpenWrt化ができます!

```

## 戻す方法
TODO

Resetボタン押してみたら文鎮化しました
