package storage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
)

// PutPolicy 表示文件上传的上传策略，参考 https://developer.qiniu.com/kodo/manual/1206/put-policy
type PutPolicy struct {
	Scope               string `json:"scope"`
	Expires             uint64 `json:"deadline"` // 截止时间（以秒为单位）
	IsPrefixalScope     int    `json:"isPrefixalScope,omitempty"`
	InsertOnly          uint16 `json:"insertOnly,omitempty"` // 若非0, 即使Scope为 Bucket:Key 的形式也是insert only
	DetectMime          uint8  `json:"detectMime,omitempty"` // 若非0, 则服务端根据内容自动确定 MimeType
	FsizeMin            int64  `json:"fsizeMin,omitempty"`
	FsizeLimit          int64  `json:"fsizeLimit,omitempty"`
	MimeLimit           string `json:"mimeLimit,omitempty"`
	ForceSaveKey        bool   `json:"forceSaveKey,omitempty"`
	SaveKey             string `json:"saveKey,omitempty"`
	CallbackFetchKey    uint8  `json:"callbackFetchKey,omitempty"`
	CallbackURL         string `json:"callbackUrl,omitempty"`
	CallbackHost        string `json:"callbackHost,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	CallbackBodyType    string `json:"callbackBodyType,omitempty"`
	ReturnURL           string `json:"returnUrl,omitempty"`
	ReturnBody          string `json:"returnBody,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyURL string `json:"persistentNotifyUrl,omitempty"`
	PersistentPipeline  string `json:"persistentPipeline,omitempty"`
	EndUser             string `json:"endUser,omitempty"`
	DeleteAfterDays     int    `json:"deleteAfterDays,omitempty"`
	FileType            int    `json:"fileType,omitempty"`
}

// UploadToken 方法用来进行上传凭证的生成
// 该方法生成的过期时间是现对于现在的时间
func (p *PutPolicy) UploadToken(cred *auth.Credentials) string {
	return p.uploadToken(cred)
}

func (p PutPolicy) uploadToken(cred *auth.Credentials) (token string) {
	if p.Expires == 0 {
		p.Expires = 3600 // 默认一小时过期
	}
	p.Expires += uint64(time.Now().Unix())
	putPolicyJSON, _ := json.Marshal(p)
	token = cred.SignWithData(putPolicyJSON)
	return
}

func getAkBucketFromUploadToken(token string) (ak, bucket string, err error) {
	items := strings.Split(token, ":")
	// KODO-11919
	if len(items) == 5 && items[0] == "" {
		items = items[2:]
	} else if len(items) != 3 {
		err = errors.New("invalid upload token, format error")
		return
	}

	ak = items[0]
	policyBytes, dErr := base64.URLEncoding.DecodeString(items[2])
	if dErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	putPolicy := PutPolicy{}
	uErr := json.Unmarshal(policyBytes, &putPolicy)
	if uErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	bucket = strings.Split(putPolicy.Scope, ":")[0]
	return
}
