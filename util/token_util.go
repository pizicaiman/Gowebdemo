package util

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
	"log"
	"time"
)

// 秘钥，如果要对接，问问上游生成的时候，有没有 秘钥
var mySigningKey = []byte("asfasfdafasdfdasfa.")

//TokenHandle 解析token

func TokenHandle(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	data := token.Claims.(jwt.MapClaims)["data"]
	var tokenData = new(TokenData)
	if data != nil {
		var info = data.(string)
		err = json.Unmarshal([]byte(info), &tokenData)
	}

	return tokenData.UserId, err
}

// CreateToken is created token
func CreateToken(data TokenData) (string, error) {
	dataByte, err := json.Marshal(data)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(dataByte),
		"exp":  time.Now().Unix() + 1000*5,
		"iss":  "ibc_business",
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("加密后的token字符串", tokenString)
	return tokenString, nil
}

// TokenData is 用户对象
type TokenData struct {
	UserId   string
	Age      int32
	NickName string
	Name     string
	Phone    string
}
