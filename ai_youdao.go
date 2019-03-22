package ai_youdao

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

type YDConfig struct {
	ApiUrl string `yaml:"api_url"`
	AppKey string `yaml:"app_key"`
	Secret string `yaml:"app_secret"`
}


var YDApiUrl string
var YDAppKey string
var Secret string
var apiUrl *url.URL
var ug uuid.Generator



func init(){
	var conf YDConfig
	filename, _ := filepath.Abs("./yaml_configs/youdao.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	//fmt.Println(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	YDApiUrl = conf.ApiUrl
	YDAppKey = conf.AppKey
	Secret = conf.Secret
	apiUrl, _ = url.Parse(YDApiUrl)

}
type YdResponse struct {
	Code string `json:"errorCode"`
	Translation []string `json:"translation"`
	Data string
}

func (resp *YdResponse) Translate() string{
	if len(resp.Translation) >=1{
		return resp.Translation[0]
	}else{
		return ""
	}

}

func Query(queryContent string) *YdResponse {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	req := packQueryParams(apiUrl ,queryContent, timeStamp)
	fmt.Println(req.String())
	res, err := http.Get(req.String());
	if err != nil {
		panic(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	var resp YdResponse
	json.Unmarshal(result, &resp)
	resp.Data = string(result)
	return &resp
}

func packQueryParams(req *url.URL, content string, curtime string) *url.URL{
	tmpUuid := uuid.Must(uuid.NewV1())
	salt := fmt.Sprintf("%x",string(tmpUuid[:]))
	q := req.Query()
	q.Set("from", "en")
	q.Set("to", "zh-CHS")
	q.Set("salt", salt )
	q.Set("appKey", YDAppKey)
	q.Set("signType", "v3")
	q.Set("curtime", curtime)
	q.Set("q", content)
	//签名信息，sha256(appKey+input+salt+curtime+密钥)
	signFinger := encrypt(YDAppKey, computeInput(content), salt, curtime, Secret)
	q.Set("sign", signFinger)
	req.RawQuery = q.Encode()
	return req
}

//其中，input的计算方式为：
// input=q前10个字符 + q长度 + q后十个字符（当q长度大于20）
// input=q字符串（当q长度小于等于20）。
func computeInput(content string) string{
	s := []rune(content)
	if len(s) <= 20{
		return content
	}else{
		return string(s[0:10]) + strconv.Itoa(len(s)) + string(s[len(s)-10:])
	}
}
