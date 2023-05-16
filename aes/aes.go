package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

// 用于加密 解密操作

// AES,即高级加密标准（Advanced Encryption Standard）,
// 是一个对称分组密码算法,旨在取代DES成为广泛使用的标准.AES中常见的有三种解决方案,分别为AES-128,AES-192和AES-256.
// AES加密过程涉及到4种操作:字节替代（SubBytes）,行移位（ShiftRows）,列混淆（MixColumns）和轮密钥加（AddRoundKey）.
// 解密过程分别为对应的逆操作.由于每一步操作都是可逆的,按照相反的顺序进行解密即可恢复明文.加解密中每轮的密钥分别由初始密钥扩展得到.
// 算法中16字节的明文,密文和轮密钥都以一个4x4的矩阵表示. AES 有五种加密模式:
// 电码本模式（Electronic Codebook Book (ECB)）,密码分组链接模式（Cipher Block Chaining (CBC)）,
// 计算器模式（Counter (CTR)）,密码反馈模式（Cipher FeedBack (CFB)）和输出反馈模式（Output FeedBack (OFB)）

func check(CKey string) cipher.Block {
	if CKey == `` {
		log.Println("密钥不能为空 ")
		panic(`密钥不能为空`)
	}
	Key := []byte(CKey)
	// 分组秘钥
	block, err := aes.NewCipher(Key)
	if err != nil {
		log.Println("key 长度必须 16/24/32长度: ", err)
		panic(err)
	}
	return block
}

// 加密
func Encrypt(origData []byte, CKey string) string {
	block := check(CKey)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	padding := blockSize - len(origData)%blockSize
	repeat := bytes.Repeat([]byte{byte(padding)}, padding)
	origData = append(origData, repeat...)

	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(CKey)[:blockSize])
	// 创建数组
	c := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(c, origData)
	// 使用RawURLEncoding 不要使用StdEncoding(放在url参数中会导致错误)
	return base64.RawURLEncoding.EncodeToString(c)
}

// 解密
func Decrypt(ciphertext, CKey string) []byte {
	// 分组秘钥
	block := check(CKey)
	// 使用RawURLEncoding 不要使用StdEncoding(放在url参数中会导致错误)
	decodeString, _ := base64.RawURLEncoding.DecodeString(ciphertext)

	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, []byte(CKey)[:blockSize])
	// 创建数组
	orig := make([]byte, len(decodeString))
	// 解密
	blockMode.CryptBlocks(orig, decodeString)

	// 去补全码
	length := len(orig)
	u := int(orig[length-1])

	return orig[:(length - u)]
}
