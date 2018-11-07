package util

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/json-iterator/go"
)

func MD5(str string) string {
	bytes := []byte(str)
	hasher := md5.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func JsonDecode(data string) (interface{}, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes := []byte(data)
	var m interface{}
	err := json.Unmarshal(bytes, &m)
	return m, err
}

func JsonEncode(m interface{}) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}
