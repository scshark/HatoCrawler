/**
 * @Author: scshark
 * @Description:
 * @File:  Common
 * @Date: 2/6/23 11:07 AM
 */
package lab

import (
	"SecCrawler/utils"
	"bytes"
	"net/http"
)

func GetUrlData(url string, getType string) (string, error) {

	client := utils.CrawlerClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	switch getType {
	case "html":
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("content-type", "text/html")
	case "json":
		req.Header.Set("Cache-Control", "max-age=0")
		req.Header.Set("accept-encoding", "gzip, deflate, br")
		req.Header.Set("content-type", "application/json")
	}

	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Proxy-Connection", "keep-alive")

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(resp.Body)

	return buffer.String(), err
}
