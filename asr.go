package qianfan

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ASRRequest 模型请求的参数结构体，但并非每个模型都完整支持如下参数，具体是否支持以 API 文档为准
type AsrRequest struct {
	DevPid  int `json:"dev_pid"`
	Rate    int `json:"rate"`
	Len     int `json:"len"`
	Channel int `json:"channel"`

	Format string `json:"format"`
	Token  string `json:"token"`
	Cuid   string `json:"cuid"`
	Speech string `json:"speech"`
}

// {"corpus_no":"7345860864562507542","err_msg":"success.","err_no":0,"result":["北京科技馆。"],"sn":"18036583701710341513"}
type AsrResponse struct {
	CorpusNo string   `json:"corpus_no"`
	ErrMsg   string   `json:"err_msg"`
	Sn       string   `json:"sn"`
	ErrNo    int      `json:"err_no"`
	Result   []string `json:"result"`
}

// 用于 ASR 模型请求的结构体
type ASR struct{}

func (s *ASR) NewASRRequest(audioFile string) (*AsrRequest, error) {
	asrReq := &AsrRequest{
		DevPid:  1537,
		Rate:    16000,
		Channel: 1,
		Cuid:    "123456PHP",
	}
	if audioFile == "" {
		return asrReq, errors.New("no audio file")
	}

	//设置文件格式
	fileExt := filepath.Ext(audioFile)
	if fileExt != "" {
		asrReq.Format = fileExt[1:]
	}

	fp, err := os.Open(audioFile)
	if err != nil {
		return asrReq, err
	}

	audioByte, err := io.ReadAll(fp)
	if err != nil {
		return asrReq, err
	}

	asrReq.Len = len(audioByte)
	asrReq.Speech = base64.StdEncoding.EncodeToString(audioByte)

	auth := GetAuthManager()
	token, err := auth.GetAccessToken(GetConfig().SpeechAK, GetConfig().SpeechSK)
	if err != nil {
		return asrReq, err
	}
	asrReq.Token = token

	return asrReq, nil
}

// 发送请求
func (s *ASR) Do(ctx context.Context, request *AsrRequest) (string, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	asrURL := GetConfig().ASRBaseURL
	fmt.Println("asrURL:", asrURL)
	client := &http.Client{}
	client.Timeout = 10 * time.Second
	httpReq, err := http.NewRequest("POST", asrURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	bytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	asrRes := new(AsrResponse)
	err = json.Unmarshal(bytes, asrRes)
	if err != nil {
		return "", err
	}

	if asrRes.ErrNo != 0 {
		return "", errors.New(asrRes.ErrMsg)
	}

	// {"corpus_no":"7345860864562507542","err_msg":"success.","err_no":0,"result":["北京科技馆。"],"sn":"18036583701710341513"}
	resp := strings.Join(asrRes.Result, "\r\n")
	return resp, nil
}

func NewASR() *ASR {
	asr := ASR{}

	return &asr
}
