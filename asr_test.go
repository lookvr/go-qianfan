package qianfan

import (
	"context"
	"fmt"
	"testing"
)

func TestASR(t *testing.T) {
	GetConfig().SpeechAK = "ectldv5dZHEg913UmF4G0zyE"
	GetConfig().SpeechSK = "IXoyvjii6Ym88yiYC4fKOWL0fza87E7L"

	audioFile := "./16k.wav"

	asr := NewASR()
	req, err := asr.NewASRRequest(audioFile)
	if err != nil {
		t.Fatalf("NewASRRequest:%v \n", err)
		return
	}

	res, err := asr.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("http dorequest:%v \n", err)
		return
	}

	fmt.Println("res:", res)
}
