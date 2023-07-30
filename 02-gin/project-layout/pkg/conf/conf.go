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

var _ Conf = (*conf)(nil)

type Conf interface {
	Load()
	Watch()

	File(name string) *viper.Viper
}

type conf struct {
	opts *Options
	// files   map[string]*viper.Viper
	files sync.Map
}

func New(opts ...Option) Conf {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &conf{
		opts: options,
		// files: make(map[string]*viper.Viper),
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
			// c.files[name] = v
			c.files.Store(name, v)
		}
	}

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
	//for _, v := range c.files {
	//	v.OnConfigChange(func(e fsnotify.Event) {
	//		log.Printf("Config file changed: %s", e.Name)
	//	})
	//	v.WatchConfig()
	//}
}

// File 根据文件名获取对应配置对象
// 如果要读取 foo.yaml 配置，可以 File("foo").Get("bar")
func (c *conf) File(name string) *viper.Viper {
	if v, ok := c.files.Load(name); ok {
		return v.(*viper.Viper)
	}
	// return c.files[name]
	return nil
}
