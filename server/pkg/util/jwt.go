package util

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义了JWT的负载结构
type Claims struct {
	UserID               uint   `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	IsRoot               bool   `json:"is_root"`  // 是否是管理员
	jwt.RegisteredClaims        // 嵌入JWT标准声明
}

// GenerateToken 生成JWT Token
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//   - isRoot: 是否是管理员
//
// 返回值:
//   - string: 生成的Token字符串
//   - error: 错误信息
func GenerateToken(userID uint, username string, isRoot bool) (string, error) {
	// 从环境变量获取JWT密钥
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	// 解析Token过期时间配置
	expirationStr := os.Getenv("JWT_EXPIRATION")
	if expirationStr == "" {
		expirationStr = "24h" // 默认24小时过期
	}
	expiration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return "", err
	}

	// 创建JWT声明(Claims)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		IsRoot:   isRoot,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                 // 生效时间
		},
	}

	// 使用HS256算法创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名Token并返回
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析并验证JWT Token
// 参数:
//   - tokenString: 待解析的Token字符串
//
// 返回值:
//   - *Claims: 解析出的声明信息
//   - error: 错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 从环境变量获取JWT密钥
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token并提取声明信息
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
