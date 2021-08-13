package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

const staticDir = "./static"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 静态目录， 不存在创建
func init() {
	dir, err := os.Stat(staticDir)
	if err != nil {
		createStaticDir(staticDir)
		return
	}
	if !dir.IsDir() {
		createStaticDir(staticDir)
	}
}

// 创建静态文件夹
func createStaticDir(staticDir string) {
	err := os.Mkdir(staticDir, 755)
	if err != nil {
		panic(err)
	}
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 返回图片存储路径
func GetHeadImgUrl(fileName string) string {

	fileSuffix := path.Ext(fileName)

	md5str := strconv.FormatInt(time.Now().Unix(), 10)

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	md5str += string(b)

	return path.Join(staticDir, Md5(md5str)+fileSuffix)
}
