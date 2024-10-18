package logger

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cutter
// Go日志文件切割
// 实现了 io.Writer 接口
type Cutter struct {
	level     string        // 日志级别（debug,info,warn,error,dpanic,panic,fatal）
	layout    string        // 时间格式 2006-01-02 15:04:05
	formats   []string      // 自定义参数([]string{Director,"2006-01-02","business"(此参数可不写),level+".log"})
	director  string        // 日志存放文件夹
	expireDay int           // 日志过期时间
	file      *os.File      // 文件句柄
	mutex     *sync.RWMutex // 读写锁
}

func (c *Cutter) Sync() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.file != nil {
		return c.file.Sync()
	}
	return nil
}

func (c *Cutter) Write(bytes []byte) (n int, err error) {
	c.mutex.Lock()
	defer func() {
		if c.file != nil {
			_ = c.file.Close()
			c.file = nil
		}
		c.mutex.Unlock()
	}()
	length := len(c.formats)
	values := make([]string, 0, 3+length)
	values = append(values, c.director)
	if c.layout != "" {
		values = append(values, time.Now().Format(c.layout))
	}
	for i := 0; i < length; i++ {
		values = append(values, c.formats[i])
	}
	values = append(values, c.level+".log")
	filename := filepath.Join(values...)
	director := filepath.Dir(filename)
	if err := os.MkdirAll(director, os.ModePerm); err != nil {
		return 0, err
	}
	if err := c.cleanUp(); err != nil {
		return 0, err
	}

	c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	return c.file.Write(bytes)
}

// CleanUp
// 检测文件夹文件创建时间，超过预定时间的文件删除
func (c *Cutter) cleanUp() error {
	if c.expireDay <= 0 {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -c.expireDay)
	return filepath.Walk(c.director, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.ModTime().Before(cutoff) && path != c.director {
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// 选项模式

type CutterOption func(*Cutter)

func NewCutter(dir string, level string, opts ...CutterOption) *Cutter {
	cutter := &Cutter{director: dir, level: level, expireDay: 1, mutex: new(sync.RWMutex)}
	for _, opt := range opts {
		opt(cutter)
	}
	return cutter
}

func WithExpireDay(days int) CutterOption {
	return func(writer *Cutter) {
		writer.expireDay = days
	}
}

func WithLayout(layout string) CutterOption {
	return func(writer *Cutter) {
		writer.layout = layout
	}
}

func WithFormats(formats ...string) CutterOption {
	return func(writer *Cutter) {
		writer.formats = formats
	}
}
