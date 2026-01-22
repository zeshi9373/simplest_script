package tool

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// 密码加密
func Password(plainpwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainpwd), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash)
	return encodePWD
}

// 密码校验
func CheckPassword(plainpwd, cryptedpwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cryptedpwd), []byte(plainpwd)) //验证（对比）
	return err == nil
}
