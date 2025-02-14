package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

/*
安全问题id
0 = 没有安全提问
1 = 母亲的名字
2 = 爷爷的名字
3 = 父亲出生的城市
4 = 您其中一位老师的名字
5 = 您个人计算机的型号
6 = 您最喜欢的餐馆名称
7 = 驾驶执照的最后四位数字
*/

const (
	action     = "login"
	username   = ""                           //用户名
	password   = ""  //密码md5
	questionid = ""                                 //安全问题ID，默认0为未设置
	answer     = ""                         //安全问题答案
	sendkey    = "" //Server酱sendkey
)

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Formhash   string `json:"formhash"`
	Mark       string `json:"mark"`
	Cookie     string
	Signsubmit string
}

var r Response

var wg sync.WaitGroup

var getCookie, _ = cookiejar.New(nil)
var client = &http.Client{Jar: getCookie}

func main() {
	pattern := flag.String("p", "s", "pattern: s / i / a")
	flag.Parse()

	if *pattern == "s" {
		login()
		r.Signsubmit = "true"
		ajaxsign(r, client)
	} else if *pattern == "i" {
		login()
		gethomepage(client)
	} else if *pattern == "a" {
		login()
		for true {
			wg.Add(2)
			r.Signsubmit = "true"
			go ajaxsign(r, client)
			go gethomepage(client)
			wg.Wait()
		}

	}

}

func login() {
	LoginData := url.Values{"action": {action}, "username": {username}, "password": {password}, "questionid": {questionid}, "answer": {answer}}
	req, err := http.NewRequest("POST", "https://www.t00ls.com/login.json", strings.NewReader(LoginData.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)

	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&r)
	if err != nil {
		return
	}
	if r.Status != "success" {
		fmt.Println("登陆失败，一小时后重试 ")
		time.Sleep(time.Hour)
		login()
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

}

// t00ls签到
func ajaxsign(r Response, client *http.Client) {
	for true {
		singData := url.Values{"signsubmit": {r.Signsubmit}, "formhash": {r.Formhash}}
		req, err := http.NewRequest("POST", "https://www.t00ls.com/ajax-sign.json", strings.NewReader(singData.Encode()))
		if err != nil {
			wg.Done()
		}
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			wg.Done()
		}
		var sign Response
		json.NewDecoder(resp.Body).Decode(&sign)
		if sign.Status == "success" {
			fmt.Println(username + "签到成功")
			push(time.Now().Format("2006/01/02 15:04") + username + " 签到成功")
			time.Sleep(time.Hour * 24)
		} else if sign.Message == "alreadysign" {
			fmt.Println(username + " 今日已完成签到。")
			push(time.Now().Format("2006/01/02 15:04") + username + " 今日已完成签到。")
			time.Sleep(time.Hour * 24)
		} else {
			fmt.Println("签到失败，1小时后重试。")
			time.Sleep(time.Hour)
			ajaxsign(r, client)
		}
	}

}

// 定时访问
func gethomepage(client *http.Client) {
	for true {
		//当前小时 每天 2-6点不访问
		dayhour := time.Now().Hour()
		fmt.Println("当前时间:", time.Now(), "; 当前小时:", dayhour)

		if dayhour < 2 || dayhour > 6 {

			req, err := http.NewRequest("GET", "https://www.t00ls.com/space-uid-15021.html", nil)
			if err != nil {
				wg.Done()
			}
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
			req.Header.Add("Content-Type", "text/html")
			resp, err := client.Do(req)
			if err != nil || resp == nil {
				wg.Done()
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp != nil && resp.StatusCode == 200 {
				//随机因子
				rand.Seed(time.Now().UnixNano())
				//获得随机分钟访问一下 10-20分钟;
				randTime := rand.Intn(10) + 10
				fmt.Println("下一次间隔分钟:", randTime)
				time.Sleep(time.Minute * time.Duration(randTime))
			}
			continue
		}
		fmt.Println("凌晨2-6点不访问,等待30分后再试。")
		time.Sleep(time.Minute * 30)

	}

}

// 方糖推送
func push(msg string) {
	urls := "https://sctapi.ftqq.com/" + sendkey + ".send?title=t00ls签到&desp=" + url.QueryEscape(msg)
	_, err := http.Get(urls)
	if err != nil {
		return
	}
	//url := "https://api.bot.wgpsec.org/push/" + sendkey + "txt=" + url.QueryEscape(msg)
	//http.Get(url)
}
