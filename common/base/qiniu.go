package base

import (
	"github.com/qiniu/api.v7/kodo"
	"qiniupkg.com/api.v7/conf"
)

var (
	Domain    string
	Bucket    string
	AccessKey string
	SecretKey string
)
var QiNiuExpire uint32 = 3600

// QiNiuUploadToken return token for upload pics
func QiNiuUploadToken() string {
	conf.ACCESS_KEY = AccessKey
	conf.SECRET_KEY = SecretKey
	c := kodo.New(0, nil)
	policy := &kodo.PutPolicy{
		Scope:   Bucket,
		Expires: QiNiuExpire,
	}
	return AccessKey + c.MakeUptoken(policy)
}

// QiNiuDownloadUrl return download url for file key
func QiNiuDownloadUrl(key string) string {
	baseUrl := kodo.MakeBaseUrl(Domain, key)
	policy := kodo.GetPolicy{}
	c := kodo.New(0, nil)
	return c.MakePrivateUrl(baseUrl, &policy)
}
