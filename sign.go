package ai_youdao

import (
	"crypto/sha256"
	"fmt"
	"strings"
)


func encrypt(params ...interface{}) string{
	var data []string
	for _, item := range params{
		data = append(data, item.(string))
	}
	encString := strings.Join(data, "")
	fmt.Println(encString)
	encData := []byte(encString)
	encryptData := sha256.Sum256(encData)
	return fmt.Sprintf("%x", string(encryptData[:]))
}

