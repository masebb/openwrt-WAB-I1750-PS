package main

import (
	"encoding/hex"
	"fmt"
	"os"
)
func main() {
  if len(os.Args) != 2 {
    fmt.Println("第二引数にELECOM WAB-I1750-PSのファームウェアファイルのパスを指定してください")
    return
  }
	data, err := os.ReadFile(os.Args[1])
  if err != nil {
    fmt.Println("ファイル読み込み時にエラーが発生しました : ",err)
    return
  }
  // 128Byte読み飛ばし
  data = data[128:]
	keyHex := "8844a2d168b45a2d"
	key, err := hex.DecodeString(keyHex)
  if err != nil {
    fmt.Println(err)
    return
  }
	keyLen := len(key)
	decryptedData := make([]byte, len(data))
	for i, b := range data {
		decryptedData[i] = b ^ key[i%keyLen]
	}
	file, err := os.Create("result.bin")
  if err != nil {
    fmt.Println("ファイル作成時にエラーが発生しました : ",err)
    return
  }
	_, err = file.Write(decryptedData)
  if err != nil {
    fmt.Println("ファイル書き込み時にエラーが発生しました : ",err)
    return
  }
  fmt.Println("ファイルの復号化が完了しました : result.bin")
}
