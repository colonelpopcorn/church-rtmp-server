package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xgfone/ngconf"
)

type ConfigController struct{}

const NGINX_PATH = "/usr/local/nginx/conf/nginx.conf"

var UNAUTHORIZED_RESPONSE = gin.H{
	SUCCESS_KEY:          false,
	RESPONSE_MESSAGE_KEY: "Unauthorized",
}

// Stream obj

type NginxConf struct {
	Content string `json:"content"`
}

func (cc *ConfigController) GetConfiguration(c *gin.Context) {
	claims := c.GetStringMap(userInfoKey)
	if !(claims["isAdmin"].(bool)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Unauthorized",
		})
		return
	}
	content, err := ioutil.ReadFile(NGINX_PATH)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Failed to read nginx conf for editing.",
			"ioUtilError":        err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		SUCCESS_KEY:          true,
		RESPONSE_MESSAGE_KEY: "Succesfully fetched nginx conf for editing.",
		"content":            string(content),
	})
}

func (cc *ConfigController) UpdateConfiguration(c *gin.Context) {
	claims := c.GetStringMap(userInfoKey)
	if !(claims["isAdmin"].(bool)) {
		c.JSON(http.StatusUnauthorized, UNAUTHORIZED_RESPONSE)
		return
	}
	var content NginxConf
	c.BindJSON(&content)
	log.Println(content.Content)
	if content.Content == "" {
		log.Fatal("I don't know jimbo...")
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Content is empty, not saving file",
		})
		return
	}
	// TODO: verify valid nginx conf
	confIsValid := testNginxConf(content.Content)
	if !confIsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Invalid config",
		})
		return
	}
	err := ioutil.WriteFile(NGINX_PATH, []byte(content.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Cannot save content",
			"ioUtilError":        err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		SUCCESS_KEY:          true,
		RESPONSE_MESSAGE_KEY: "Successfully saved modified nginx conf.",
	})

}

func testNginxConf(content string) (b bool) {
	_, err := ngconf.Decode(content)
	return err == nil
}
