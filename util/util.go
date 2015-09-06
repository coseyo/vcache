package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func MD5(str string) string {
	bytes := []byte(str)
	hasher := md5.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func JsonDecode(data string) (interface{}, error) {
	bytes := []byte(data)
	var m interface{}
	err := json.Unmarshal(bytes, &m)
	return m, err
}

func JsonEncode(m interface{}) (string, error) {
	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}
