package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const splash string = `

■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■
■                                           ■
■   █▀▄▀█ ██    ▄  █ █▄▄▄▄ ▄███▄   █▀▄▀█    ■
■   █ █ █ █ █  █   █ █  ▄▀ █▀   ▀  █ █ █    ■
■   █ ▄ █ █▄▄█ ██▀▀█ █▀▀▌  ██▄▄    █ ▄ █    ■
■   █   █ █  █ █   █ █  █  █▄   ▄▀ █   █    ■
■      █     █    █    █   ▀███▀      █     ■
■     ▀     █    ▀    ▀              ▀      ■
■          ▀                                ■
■      ▄▄▄▄▄   ▄█▄    ██      ▄             ■
■     █     ▀▄ █▀ ▀▄  █ █      █            ■
■   ▄  ▀▀▀▀▄   █   ▀  █▄▄█ ██   █           ■
■    ▀▄▄▄▄▀    █▄  ▄▀ █  █ █ █  █           ■
■              ▀███▀     █ █  █ █           ■
■                       █  █   ██           ■
■                      ▀                    ■
■          Mahrem Scan Go Port v0.9         ■
■              by Anzirra team              ■
■                                           ■
■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■

`

var limit int
var name string
var delay float32

func buildRandomPath(length int) string {
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"

	var output strings.Builder

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}

func buildUri(mahremUrl url.URL, path string) string {
	mahremUrl.Path = path
	return mahremUrl.String()
}

func mahremRequest(mahremUri string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", mahremUri, nil)

	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ""
	}

	return string(body)
}

func filterResult(body string) string {
	r, _ := regexp.Compile("\"image-container.*<img class=\".*\" src=\"(.+)\" crossorigin=\"anonymous\"")

	if r.MatchString(body) {
		if len(r.FindStringSubmatch(body)[1]) > 1 {
			return r.FindStringSubmatch(body)[1]
		}
	}

	return ""
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println(splash)

	name = fmt.Sprintf("mahrem%d.html", time.Now().Unix())

	fmt.Println("File name:", name)

	if _, err := os.Stat(name); os.IsNotExist(err) {
		f, _ := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(fmt.Sprintf(`
		<div style="
			text-shadow: 5px 5px 3px rgba(0, 0, 0, .3); 
			color: #222222; 
			display: block;
			text-align: center;
		">
			<pre>
			%s
			</pre>
		</div>
		<div></div>
		`, splash))
	}

	mahremUrl := url.URL{
		Scheme: "https",
		Host:   "prnt.sc",
	}

	for i := 0; i < 500; i++ {
		imgarea := ""
		itemPath := ""
		itemUri := ""
		itemContent := ""

		for true {
			itemPath = buildRandomPath(rand.Intn(7-5) + 5)
			itemUri = buildUri(mahremUrl, itemPath)

			itemContent = filterResult(mahremRequest(itemUri))

			if strings.HasPrefix(itemContent, "//") {
				continue
			}

			if len(itemContent) < 1 {
				continue
			}

			imgarea = fmt.Sprintf(`
			
			<div style="
				width:250px;
				margin: 10px;
				text-align: center;
			">
				<a target="_blank" href="%s">
					<img src="%s" width="250" align="middle" />
				</a>
			</div>
			`, itemUri, itemContent)

			fmt.Println(itemContent)

			time.Sleep(1 * time.Second)

			break
		}

		f, _ := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)

		f.WriteString(fmt.Sprintf(`
		<div style="
			width: 280px;
			margin: 10px auto;
			background-color: crimson;
			border-radius: 10px;
			padding: 10px;
			display: inline-block;
		">
			<div style="width:200px;">
				<a target="_blank" href="%s" style="
					text-decoration: none;
					color: #ffffff;
					font-weight: bold;
					display: block;
					margin-bottom: 10px;
					padding: 10px 20px;
				">
					%s
				</a>
			</div>
			%s
		</div>
		
		`, itemUri, itemPath, imgarea))

		defer f.Close()
	}
}
