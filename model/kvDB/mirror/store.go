package mirror

import (
	"fmt"
	json "github.com/json-iterator/go"
	"mirror-api/config"
	"mirror-api/util"
	"os"
)

var hostname = ""
var mirrorMap map[string]Mirror

var (
	ErrNotFound = fmt.Errorf("not found")
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

	newMirrorMap := make(map[string]Mirror)
	for _, m := range d.Mirrors {
		newMirrorMap[m.Id] = m
	}

	mirrorMap = newMirrorMap
}

func GetAll() map[string]Mirror {
	return mirrorMap
}

func GetOneByID(id string) (*Mirror, error) {
	m, e := mirrorMap[id]
	if !e {
		return nil, ErrNotFound
	}

	return &m, nil
}

func GetHostname() string {
	return hostname
}
