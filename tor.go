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
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// For more information, please refer to <http://unlicense.org/>
package dfi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	torc "github.com/postfix/goControlTor"
	log "github.com/sirupsen/logrus"
)

func SetupDFITorService(port, tor int, cookie string) (*torc.TorControl, string, error) {
	control := &torc.TorControl{}

	serviceDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Info(serviceDir)
	servicePort := map[int]string{port: fmt.Sprintf("127.0.0.1:%d", port)}

	err := control.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tor))

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Dialed Tor control port")

	err = control.CookieAuthenticate(cookie)

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Authenticated with Tor, creating service")

	err = control.CreateHiddenService(serviceDir, servicePort)

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Service created")

	onion, err := torc.ReadOnion(serviceDir)
	onion = strings.TrimSpace(onion)

	log.WithField("onion", onion).Info("Connecting to Tor")

	return control, onion, nil
}
