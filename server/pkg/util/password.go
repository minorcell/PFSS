package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// defaultBcryptCost 定义密码哈希的计算成本
	defaultBcryptCost = 12
	// saltLength 定义盐值的字节长度
	saltLength = 16
)

// GenerateSalt 生成随机盐值
// 返回值:
//   - string: 生成的盐值(Base64编码)
//   - error: 错误信息
func GenerateSalt() (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword 对密码进行哈希处理
// 参数:
//   - password: 需要哈希的原始密码
//
// 返回值:
//   - string: 哈希后的密码
//   - error: 错误信息
func HashPassword(password string) (string, error) {
	// 使用bcrypt生成密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultBcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword 验证密码是否匹配哈希值
// 参数:
//   - password: 待验证的原始密码
//   - hash: 存储的密码哈希
//
// 返回值:
//   - bool: 验证结果(true表示匹配)
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
