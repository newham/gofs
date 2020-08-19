package api

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type CSV struct {
	lock  sync.Mutex
	path  string
	lines []Line
}

type Line []string

const tagDeleted = "deleted"

func (l *Line) Id() string {
	if len(*l) > 0 {
		return (*l)[0]
	}
	return ""
}

func (l *Line) ToString() string {
	return strings.Join(*l, ",")
}

func (l *Line) IsDel() bool {
	return l.Id() == tagDeleted
}

func (l *Line) Text() string {
	return strings.Join(*l, ",")
}

func (l *Line) Print() {
	println(l.Text())
}

func NewLine(id string, data ...string) Line {
	l := []string{id}
	return append(l, data...)
}

func NewCSV(path string) (*CSV, error) {
	Mkfile(path) //不存在就创建一份
	c := &CSV{path: path, lines: []Line{}}
	if err := c.Load(path); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CSV) Get(id string) Line {
	var line Line
	c.Map(func(i int, l Line) bool {
		if l.Id() == id {
			line = l
			return true
		}
		return false
	})
	return line
}

func (c *CSV) Map(f func(int, Line) bool) {
	for i, l := range c.lines {
		if l.IsDel() {
			continue
		}
		if f(i, l) {
			break
		}
	}
}

func (c *CSV) Del(id string) {
	c.Map(func(i int, l Line) bool {
		if l.Id() == id {
			c.lines[i] = []string{tagDeleted}
		}
		return false
	})
}

func (c *CSV) Put(line Line) {
	if line == nil {
		return
	}
	put := true
	c.Map(func(i int, l Line) bool {
		if l.Id() == line.Id() {
			c.lines[i] = line
			put = false
		}
		return false
	})
	if put {
		c.lines = append(c.lines, line)
	}

}

func (c *CSV) ToBytes() []byte {
	str := ""
	for _, l := range c.lines {
		str += l.ToString() + "\n"
	}
	return []byte(str)
}

func (c *CSV) Print() {
	c.Map(func(i int, l Line) bool {
		fmt.Printf("NO.%-3d", i)
		l.Print()
		return false
	})
}

func (c *CSV) SaveTo(path string) error {
	if path == "" {
		return errors.New("no path")
	}
	//eb := Base64Encode(c.ToBytes()) //encrypt
	eb := c.ToBytes()
	c.lock.Lock()
	err := OverwriteBytes(path, eb)
	c.lock.Unlock()
	return err
}

func (c *CSV) Save() error {
	return c.SaveTo(c.path)
}

func (c *CSV) Load(path string) error {
	if path == "" {
		return errors.New("no path")
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		println(err.Error())
		return err
	}
	//db := Base64Decode(b) //encrypt
	db := b
	reader := bufio.NewReader(bytes.NewReader(db))
	if reader == nil {
		return errors.New("empty file")
	}
	for {
		lineStr, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line := strings.Split(strings.Trim(lineStr, " \n"), ",")
		c.lines = append(c.lines, line)
	}
	return nil
}
