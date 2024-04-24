package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

// md5加密密码
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 定义需求（claims），也就是需要通过 jwt 传输的数据
type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// 定义密钥
var myKey = []byte("gin-gorm-oj-key")

// 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	userClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 验证
	if !claims.Valid {
		return nil, fmt.Errorf("Analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// 发送验证码
// 邮箱授权码xxxxx   base64:xxxxxxx
// xxxxxxxx
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "xxxxx@qq.com"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码:<b>" + code + "</b>")
	return e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "xxxxxx@qq.com", "邮箱授权码xxxx", "smtp.qq.com:465"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})

}

// 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()

}

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// 代码保存
func CodeSave(code []byte) (string, error) {
	// 在code目录下创建目录
	// 代码提交给一个唯一标识作为目录名
	dirName := "code/" + GetUUID()
	// 该目录下有个main.go文件
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	// 创建main.go文件
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	// 写入代码
	f.Write(code)
	defer f.Close()
	return path, nil
}
