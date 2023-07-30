package conf

import (
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	_ Conf = (*conf)(nil)

	//Get 读取配置项 fileName 文件名, key 配置项
	Get func(fileName string, key string) any
)

type Conf interface {
	// Load 加载配置文件
	Load()
	// Watch 监听配置文件变化
	Watch()

	//File name 文件名
	File(name string) *viper.Viper

	Get(fileName string, key string) any
}

type conf struct {
	opts *Options

	files sync.Map
}

func New(opts ...Option) Conf {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &conf{
		opts:  options,
		files: sync.Map{},
	}
}

func (c *conf) Load() {
	for _, source := range c.opts.sources {
		fs, err := source.Load()
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			v := viper.New()
			v.SetConfigType(f.Format)
			v.SetConfigFile(f.Path)
			v.AutomaticEnv()

			if err := v.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					log.Printf("Using conf file: %s [%s]\n", viper.ConfigFileUsed(), err)
				}
				panic(err)
			}
			v.AutomaticEnv()

			name := strings.TrimSuffix(path.Base(f.Key), filepath.Ext(f.Key))
			log.Printf("文件名: %s", name)
			c.files.Store(name, v)
		}
	}

	Get = c.Get
}

func (c *conf) Watch() {
	c.files.Range(func(key, value any) bool {
		v := value.(*viper.Viper)
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s", e.Name)
		})
		v.WatchConfig()
		return true
	})
}

func (c *conf) File(name string) *viper.Viper {
	if v, ok := c.files.Load(name); ok {
		return v.(*viper.Viper)
	}
	return nil
}

func (c *conf) Get(fileName string, key string) any {
	return c.File(fileName).Get(key)
}
