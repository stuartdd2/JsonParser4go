package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetJsonParsed(getUrl string) (NodeI, error) {
	data, err := GetJson(getUrl)
	if err != nil {
		return nil, err
	}
	n, err := Parse(data)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func GetJson(getUrl string) ([]byte, error) {
	resp, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get data from server. Status is not 200. Code:%d Url:%s", resp.StatusCode, getUrl)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func PostJsonBytes(postUrl string, data []byte) ([]byte, error) {
	return PostData(postUrl, "application/octet-stream", data)
}

func PostJsonText(postUrl string, data []byte) ([]byte, error) {
	return PostData(postUrl, "application/json", data)
}

func PostJsonValue(postUrl string, node NodeI) ([]byte, error) {
	return PostJsonText(postUrl, []byte(node.JsonValue()))
}

func PostJsonValueIndented(tab int, postUrl string, node NodeI) ([]byte, error) {
	return PostJsonText(postUrl, []byte(node.JsonValueIndented(tab)))
}

func PostData(postUrl string, contentType string, data []byte) ([]byte, error) {
	resp, err := http.Post(postUrl, contentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to post data to server. Status is not 201. Code:%d Url:%s", resp.StatusCode, postUrl)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
