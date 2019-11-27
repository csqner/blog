/*
  @Author : lanyulei
*/

package auth

import (
	"blog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// 加密密码
func EncryptionPassword(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败，错误信息: %s", err)
		return ""
	}
	return string(hashPassword)
}

// 验证密码，正确则返回true，错误则返回false
func DecryptionPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		logger.Errorf("密码解密失败，错误信息: %s", err)
		return false
	}
	return true
}
