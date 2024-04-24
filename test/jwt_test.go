package test

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// 定义需求（claims），也就是需要通过 jwt 传输的数据
type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

// 定义密钥
var myKey = []byte("gin-gorm-oj-key")

// 生成token
func TestGenerateToken(t *testing.T) {
	userClaim := &UserClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}
	// jwt.NewWithClaims 方法根据 Claims 结构体创建 Token 示例。
	// 此时token 为*jwt.Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	// 通过密钥生成token
	// 通过SignedString方法转化为string类型
	SignedString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(SignedString)
}

// 解析token
func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE"
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}
