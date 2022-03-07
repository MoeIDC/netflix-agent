package utils

import (
	"crypto/tls"
	"github.com/evsio0n/log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func TestUnblock() bool {
	var netflixV6 string
	ns, err := net.LookupHost("www.netflix.com")
	if err != nil {
		log.Error("Err: " + err.Error())
		return false
	}

	for _, n := range ns {
		if net.ParseIP(n).To4() == nil {
			netflixV6 = "[" + n + "]"
		}
		break
	}

	res := RequestIP("https://www.netflix.com/title/70143836", netflixV6)
	log.Info(res)
	return !strings.Contains(res, "Error") && !strings.Contains(res, "Ban")
}

func RequestIP(requrl string, ip string) string {
	urlValue, err := url.Parse(requrl)
	if err != nil {
		log.Error("Url Parse error: " + err.Error())
		return "Error"
	}
	host := urlValue.Host
	if ip == "" {
		ip = host
	}
	newrequrl := strings.Replace(requrl, host, ip, 1)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: host},
			// goodryb pull
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				LocalAddr: &net.TCPAddr{
					IP: net.ParseIP(""),
				},
			}).DialContext,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       5 * time.Second,
	}
	req, err := http.NewRequest("GET", newrequrl, nil)
	if err != nil {
		//return errors.New(strings.ReplaceAll(err.Error(), newrequrl, requrl))
		log.Error("http request err: " + err.Error())
		return "Error"
	}
	req.Host = host
	req.Header.Set("USER-AGENT", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		//return errors.New(strings.ReplaceAll(err.Error(), newrequrl, requrl))
		log.Error("http request err: " + err.Error())
		return "Error"
	}
	defer resp.Body.Close()

	Header := resp.Header

	if Header["X-Robots-Tag"] != nil {
		if Header["X-Robots-Tag"][0] == "index" {
			return "us"
		}
	}

	if Header["Location"] == nil {
		return "Ban"
	} else {
		return strings.Split(Header["Location"][0], "/")[3]
	}
}
