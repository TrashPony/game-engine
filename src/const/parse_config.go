package _const

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var Config appConfig

type appConfig struct {
	configFields map[string]string
}

func (c *appConfig) initConfig() {

	c.configFields = make(map[string]string)

	configFile, err := ioutil.ReadFile("./main.ini")
	if err != nil {
		log.Fatal("фаил с конфигурацией ./main.ini - не найден")
	}

	configLines := strings.Split(string(configFile), "\n")
	for i := 0; i < len(configLines); i++ {

		if strings.Contains(configLines[i], "[") && strings.Contains(configLines[i], "]") {
			continue
		}

		if configLines[i] != "" {

			configLine := strings.Split(configLines[i], "=")

			if len(configLine) > 1 {
				c.configFields[strings.TrimSpace(configLine[0])] = strings.TrimSpace(configLine[1])
			}
		}
	}
}

func (c *appConfig) GetParams(paramsKey string) string {

	if c.configFields == nil {
		c.initConfig()
	}

	return c.configFields[paramsKey]
}

func (c *appConfig) GetIntParams(paramsKey string) int {

	if c.configFields == nil {
		c.initConfig()
	}

	result, err := strconv.Atoi(c.configFields[paramsKey])
	if err != nil {
		panic("parse int config: " + err.Error())
	}

	return result
}
