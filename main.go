package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func main() {

	headers := make(http.Header)
	headers.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36`)

	uri := &url.URL{Scheme: "https", Host: "www17.muenchen.de", Path: "/Passverfolgung/"}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Do(&http.Request{Method: "GET", URL: uri, Header: headers})
	if err != nil {
		log.Fatalf("Request1 failed %s", err)
	}

	/* Search CSRF id */
	var matchNc = regexp.MustCompile(`__ncforminfo" value="(.*)"\/>`)
	var ncFormInfo string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "__ncforminfo") {
			se := matchNc.FindStringSubmatch(scanner.Text())
			if len(se) > 1 {
				ncFormInfo = se[1]
			}
		}
	}
	resp.Body.Close()

	if ncFormInfo == "" {
		//log.Fatal("Cant CSRF id")
	} else {
		fmt.Println("CSRF id", ncFormInfo)
	}

	data := url.Values{}
	data.Add("__ncforminfo", ncFormInfo)
	data.Add("ifNummer", "39393")
	data.Add("pbAbfragen", "Abfragen")

	fmt.Println(jar.Cookies(uri))
	headers.Add("Content-Type", "application/x-www-form-urlencoded")
	headers.Add("Referer", "https://www17.muenchen.de/Passverfolgung/")
	headers.Add("Origin", "https://www17.muenchen.de")
	headers.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err = client.Do(&http.Request{URL: uri, Method: "POST", Header: headers, Body: ioutil.NopCloser(bytes.NewReader([]byte(data.Encode())))})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTTP Status Code ist " , resp.Status)
	fmt.Println("HTTP Respnonse Header sind ", resp.Header)

	//io.Copy(os.Stdout, resp.Body)

  resp.Body.Close()

}
