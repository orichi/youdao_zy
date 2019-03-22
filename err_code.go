package ai_youdao

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type ErrorCode struct {
	Description string
	Codes map[string]string
}

var ErrCode ErrorCode

func init() {
	filename, _ := filepath.Abs("./yaml_configs/yd_err_code.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &ErrCode)
	if err != nil {
		panic(err)
	}
}

func (errCode ErrorCode) Code(code string) string{
	return errCode.Codes[code]
}
