package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha1"
	s "strings"
	"sort"
  "time"
	"strconv"
	"bytes"

)

const twitterApiBase = "https://api.twitter.com";

func PercentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {

		if shouldEscape(b) {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {

			buf.WriteByte(b)
		}
	}
	return buf.String()
}


func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '.', '_', '~':
		return false
	}
	return true
}

func getTimeAsString() string {
	timeN := time.Now()
	now := strconv.FormatInt(timeN.Unix(), 10)
	return now
}
func getNonce() string {
	base := RandStringBytesMaskImprSrcSB(32)
	return base64.StdEncoding.EncodeToString([]byte(base))
}
func GetSignature(input, key string) string {
    key_for_sign := []byte(key)
    h := hmac.New(sha1.New, key_for_sign)
    h.Write([]byte(input))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getOAuthString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, PercentEncode(k) + "=\"" + PercentEncode(params[k]) + "\"")
	}
	sort.Strings(keys)
	return "OAuth " + s.Join(keys, ", ")
}

func signParams(method string, reqUrl string, params map[string]string, key string, token string) string {
	var signingKey = ""
	if len(token) > 0 {
		signingKey = PercentEncode(key) + "&" + PercentEncode(token)
	} else {
		signingKey = PercentEncode(key) + "&"
	}
	var str = s.ToUpper(method) + "&" + PercentEncode(reqUrl)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, PercentEncode(k) + "=" + PercentEncode(params[k]))
	}
	sort.Strings(keys)
	str += "&" + PercentEncode(s.Join(keys, "&"))
	return GetSignature(str, signingKey)
}


func PostTwitterUpdate(token, secret, status string) (bool, string) {
	var reqUrl = twitterApiBase + "/1.1/statuses/update.json"
	nonce := getNonce()
	time := getTimeAsString()
 	m := map[string]string{
		"oauth_consumer_key": os.Getenv("TWITTER_CONSUMER_KEY"),
		"oauth_nonce": nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp": time,
		"oauth_version": "1.0",
		"oauth_token": token,
		"status": status,
	}
	signature := signParams("POST", reqUrl, m, os.Getenv("TWITTER_CONSUMER_SECRET"), secret)
	m["oauth_signature"] = signature
	delete(m, "status")
	authHeader := getOAuthString(m)
	client := &http.Client{}
	request, _ := http.NewRequest("POST", reqUrl + "?status=" + PercentEncode(status),nil)
	request.Header.Add("Authorization", authHeader)

	resp, err := client.Do(request)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		return false, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return true, string(body)
}
