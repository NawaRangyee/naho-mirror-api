package mirror

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
	"mirror-api/config"
	"mirror-api/util"
	"os"
	"time"
)

var hostname = ""

var mirrorCache *cache.Cache

func init() {
	mirrorCache = cache.New(12*time.Hour, 1*time.Hour)
}

var (
	ErrNotFound = fmt.Errorf("not found")
)

const (
	StatusCheckFailed = "failed_to_check"
	StatusRunning     = "running"
	StatusNotRunning  = "not_running"
)

func init() {
	InitFromConfig()
}

func InitFromConfig() {
	f, err := os.Open(*config.ConfigFilePath)
	if err != nil {
		util.SendFailed("model/kvDB/mirror/store.go - InitFromConfig() - os.Open(*config.ConfigFilePath)", err)
		return
	}
	defer f.Close()

	d := &Data{}
	err = json.NewDecoder(f).Decode(d)
	if err != nil {
		util.SendFailed("model/kvDB/mirror/store.go - InitFromConfig() - json.NewDecoder(f).Decode(d)", err)
		return
	}

	for _, m := range d.Mirrors {
		mirrorCache.Set(m.Id, m, cache.DefaultExpiration)
	}

}

func GetAll() map[string]Mirror {
	result := make(map[string]Mirror)

	for k, v := range mirrorCache.Items() {
		result[k] = v.Object.(Mirror)
	}

	return result
}

func GetOneByID(id string) (Mirror, error) {
	m, e := mirrorCache.Get(id)
	if !e {
		return Mirror{}, ErrNotFound
	}

	return m.(Mirror), nil
}

func GetHostname() string {
	return hostname
}

func SetOneByID(id string, m Mirror) {
	mirrorCache.Set(id, m, cache.DefaultExpiration)
}
