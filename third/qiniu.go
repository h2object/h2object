package third

import (
	"fmt"
	"time"
	"strings"
	sysio "io"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/io"
	"github.com/h2object/h2object/object"
)

func QiniuKey(uri string) string {
	s := strings.TrimPrefix(uri, "/")
	s = strings.Replace(s, "/", "-", -1)
	return s
}


type QiniuHelper struct{
	AccessKey string
	SecretKey string
	BucketName string
	Domain	string
	cache 	object.Cache
}

func NewQiniuHelper(accessKey, secretKey string, cache object.Cache) *QiniuHelper {
	conf.ACCESS_KEY = accessKey
	conf.SECRET_KEY = secretKey
	return &QiniuHelper{
		AccessKey: accessKey,
		SecretKey: secretKey,
		cache: cache,
	}
}

func (helper *QiniuHelper) Put(bucket string, key string, rd sysio.Reader, size int64) (string, error) {
	token := helper.token(bucket)

	var ret io.PutRet
	var extra = &io.PutExtra{
	// Params:   params,
	// MimeType: contentType,
	// Crc32:    crc32,
	// CheckCrc: CheckCrc,
	}
	if err := io.Put2(nil, &ret, token, key, rd, size, extra); err != nil {
		helper.cache.Delete(fmt.Sprintf("%s:%s", helper.AccessKey, bucket))
		return "", err
	}
	return ret.Key, nil
}

func (helper *QiniuHelper) PutFile(bucket string, key string, fn string) (string, error) {
	token := helper.token(bucket)

	var ret io.PutRet
	var extra = &io.PutExtra{
	// Params:   params,
	// MimeType: contentType,
	// Crc32:    crc32,
	// CheckCrc: CheckCrc,
	}
	if err := io.PutFile(nil, &ret, token, key, fn, extra); err != nil {
		helper.cache.Delete(fmt.Sprintf("%s:%s", helper.AccessKey, bucket))
		return "", err
	}
	return ret.Key, nil
}

func (helper *QiniuHelper) token(bucketName string) string {
	tk, ok := helper.cache.Get(fmt.Sprintf("%s:%s", helper.AccessKey, bucketName))
	if ok {
		return tk.(string)
	} 

	duration := time.Hour
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
		Expires: 3600,
		// CallbackUrl:  callbackUrl,
		// CallbackBody: callbackBody,
		// ReturnUrl:    returnUrl,
		// ReturnBody:   returnBody,
		// AsyncOps:     asyncOps,
		// EndUser:      endUser,
		// Expires:      expires,
	}

	tk2 := putPolicy.Token(nil)
	helper.cache.Set(fmt.Sprintf("%s:%s", helper.AccessKey, bucketName), tk2, duration)
	return tk2
}
