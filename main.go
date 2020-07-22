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
)

const localConfigPath = "config.json"

func main() {
	configUrl := flag.String("u", "", "Remote config")
	port := flag.Int("p", 8006, "Port")
	flag.Parse()

	// 引入配置
	var config string
	if *configUrl == "" {
		// TODO: 引入外部配置
		content, err := ioutil.ReadFile(localConfigPath)
		if err != nil {
			panic(err.Error())
		}

		config = string(content)

	} else {
		response, err := http.Get(*configUrl)
		defer response.Body.Close()

		if err != nil {
			panic(err.Error())
		}

		body, _ := ioutil.ReadAll(response.Body)
		config = string(body)
	}

	log.Println(config)
	if len(config) == 0 {
		panic("Config not found!")
	}

	var configMap map[string]string
	err := json.Unmarshal([]byte(config), &configMap)
	if err != nil {
		panic("Invalid Config!")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		host := c.Request.Header.Get("X-Forwarded-Host")
		suffix := strings.Split(host, ".")[0]
		redirect := configMap[suffix]
		log.Println(host, suffix, redirect)
		if redirect != "" {
			c.Redirect(http.StatusMovedPermanently, redirect)
		} else {
			c.JSON(http.StatusNotFound, "Redirect url undefined!")
		}
	})

	router.Run(fmt.Sprintf(":%d", *port))
}
