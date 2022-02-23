package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRETKEY = "243223ffslsfsldf" //私钥
)

type MyClaims struct {
	Uid int `json:"uid"`
	Rid int `json:"rid"`
	jwt.StandardClaims
}

func CreateToken(uid int, rid int) (string, error) {
	maxAge := 60 * 100 // 定义失效时间
	var claims = MyClaims{
		uid,
		rid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // 过期时间
			Issuer:    "yangjia",                                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) { // 解析token
		return []byte(SECRETKEY), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ParseToken2(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) { // 解析token
		return []byte(SECRETKEY), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

type tokenAuth interface {
	TokenAuth(string) (*MyClaims, error)
}

type PermissionAuth interface {
	permissionAuth(string) bool
}

type tokenAuthHadle func(string) (*MyClaims, error)
type permissionAuthHadle func(string) bool

func (t tokenAuthHadle) TokenAuth(token string) (*MyClaims, error) {
	return t(token)
}

func (t permissionAuthHadle) permissionAuth(permission string) bool {
	return t(permission)
}

// 身份验证方式1
func tokenAuth1(token string) (*MyClaims, error) {
	return ParseToken(token)
}

// 身份验证方式2
func tokenAuth2(token string) (*MyClaims, error) {
	return ParseToken2(token)
}

func permissionAuth(string) bool {
	return true
}

var AuthMap = make(map[string]tokenAuth)

// 将匹配规则添加到map 中, 可以根据 AuthMap["auth"] 去获取对应的解析方法,具体件认证中间件
// 权限验证的还没有写
func InitAuth() {
	AuthMap["tokenAuth1"] = tokenAuthHadle(tokenAuth1)
	AuthMap["tokenAuth2"] = tokenAuthHadle(tokenAuth2)
}
