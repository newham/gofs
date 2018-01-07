/******************************************
*FileName: conf.go
*Author: Liu han
*Date: 2016-12-9
*Description: read conf file
*******************************************/

package api

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var ConfigFile = "./conf/app.conf"

// var ConfigFile = "./test.conf" //for test

type ConfigInterface interface {
	//Set(key, val string) error   // support section::key type in given key when using ini type.
	String(key string) string    // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string //get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultval string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultval []string) []string //get string slice
	DefaultInt(key string, defaultval int) int
	DefaultInt64(key string, defaultval int64) int64
	DefaultBool(key string, defaultval bool) bool
	DefaultFloat(key string, defaultval float64) float64
	//DIY(key string) (interface{}, error)
	//GetSection(section string) (map[string]string, error)
	//SaveConfigFile(filename string) error
}

type Key struct {
	Name  string
	Value string
}
type Config struct {
	File string
	Keys map[string]string
}

var AppConfig ConfigInterface

func init() {
	config := &Config{ConfigFile, make(map[string]string)}
	if err := config.Prase(); err != nil {
		panic(err)
	}
	AppConfig = config
}

func (c *Config) Prase() error {
	isEnd := false

	f, err := os.Open(c.File)
	defer f.Close()
	if err != nil {
		return errors.New("Open file" + ConfigFile + " failed")
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if line != "" {
				isEnd = true
			} else {
				break
			}
		}
		line = strings.TrimSpace(line)
		if isCommentOut(line) {
			continue
		}

		firstIndex := strings.Index(line, "=")
		if firstIndex < 1 {
			continue
		} else {
			c.Keys[strings.Trim(line[:firstIndex], "\" ")] = strings.Trim(line[firstIndex+1:], "\" ")
		}

		if isEnd {
			break
		}
	}
	return nil
}

func isCommentOut(line string) bool {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "*") {
		return true
	} else {
		return false
	}
}

func (c *Config) String(key string) string {
	return c.Keys[key]
}
func (c *Config) Strings(key string) []string {
	if c.Keys[key] == "" {
		return make([]string, 0)
	} else {
		return strings.Split(c.Keys[key], " ")
	}
}
func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.Keys[key])
}
func (c *Config) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.Keys[key], 10, 64)
}
func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.Keys[key])
}
func (c *Config) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.Keys[key], 64)
}

func (c *Config) DefaultString(key string, defaultval string) string {
	if c.String(key) == "" {
		return defaultval
	} else {
		return c.String(key)
	}
}
func (c *Config) DefaultStrings(key string, defaultval []string) []string {
	if len(c.Strings(key)) < 1 {
		return defaultval
	} else {
		return c.Strings(key)
	}
}
func (c *Config) DefaultInt(key string, defaultval int) int {
	if value, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return value
	}
}
func (c *Config) DefaultInt64(key string, defaultval int64) int64 {
	if value, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return value
	}
}
func (c *Config) DefaultBool(key string, defaultval bool) bool {
	if value, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return value
	}
}
func (c *Config) DefaultFloat(key string, defaultval float64) float64 {
	if value, err := c.Float(key); err != nil {
		return defaultval
	} else {
		return value
	}
}
