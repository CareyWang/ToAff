package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const localConfigPath = "config.json"
var lastUpdateTime int64

func main() {
	configUrl := flag.String("u", "", "Remote config")
	port := flag.Int("p", 8006, "Port")
	flag.Parse()

	// 引入配置
	configMap := loadConfig(*configUrl)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		host := c.Request.Header.Get("X-Forwarded-Host")
		suffix := strings.Split(host, ".")[0]
		redirect := configMap[suffix]
		log.Println(host, suffix, redirect)

		// 未获取到url时，重新加载配置
		if redirect == "" && time.Now().Unix() - lastUpdateTime > 3600 {
			configMap = loadConfig(*configUrl)
		}

		if redirect != "" {
			c.Redirect(http.StatusMovedPermanently, redirect)
		} else {
			c.String(http.StatusNotFound, "Redirect url undefined!")
		}
	})

	router.Run(fmt.Sprintf(":%d", *port))
}

func loadConfig(configUrl string) map[string]string {
	var config string
	if configUrl == "" {
		// TODO: 引入外部配置
		content, err := ioutil.ReadFile(localConfigPath)
		if err != nil {
			panic(err.Error())
		}

		config = string(content)

	} else {
		response, err := http.Get(configUrl)
		defer response.Body.Close()

		if err != nil {
			panic(err.Error())
		}

		body, _ := ioutil.ReadAll(response.Body)
		config = string(body)
	}

	if len(config) == 0 {
		panic("Config not found!")
	}

	var configMap map[string]string
	err := json.Unmarshal([]byte(config), &configMap)
	if err != nil {
		panic("Invalid Config!")
	}

	lastUpdateTime = time.Now().Unix()
	log.Printf("%d: 更新配置!", lastUpdateTime)
	log.Println(config)

	return configMap
}
