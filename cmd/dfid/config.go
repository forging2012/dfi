// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// For more information, please refer to <http://unlicense.org/>

package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetupConfig() {
	// bind the "bind" flags
	flag.String("bind", "0.0.0.0:5050", "The address and port to listen for dfi protocol connections")
	flag.String("http", "127.0.0.1:8080", "The address and port to listen on for http commands")
	flag.Parse()

	viper.BindPFlag("bind.dfi", flag.Lookup("bind"))
	viper.BindPFlag("bind.http", flag.Lookup("http"))

	viper.SetConfigName("dfid")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dfi")
	viper.AddConfigPath("/etc/dfi")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error loading config file: %s \n", err))
	}

	viper.SetDefault("bind", map[string]string{
		"dfi":  "0.0.0.0:5050",
		"http": "127.0.0.1:8080",
	})

	// someday support postgresql, etc. Hence the map :)
	viper.SetDefault("database", map[string]string{
		"path": "./data/posts.db",
	})

	viper.SetDefault("tor", map[string]interface{}{
		"enabled":    true,
		"control":    10051,
		"socks":      10050,
		"cookiePath": "./tor/",
	})

	viper.SetDefault("socks", map[string]interface{}{"enabled": true, "port": 10050})

	viper.SetDefault("net", map[string]interface{}{
		"maxPeers": 100,
	})

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed, reloading: ", e.Name)
	})
}
