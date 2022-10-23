package checkRsync

import (
	"bytes"
	"fmt"
	"mirror-api/model/kvDB/mirror"
	"mirror-api/util"
	"mirror-api/util/logger"
	"os/exec"
	"strings"
)

func Run() {
	if err := ProcessCheck(); err != nil {
		util.SendFailed("checkRsync/ps.go - ProcessCheck()", err)
	}

}

func ProcessCheck() error {
	cMirrors := mirror.GetAll()
	for _, v := range cMirrors {
		status, err := ProcessCheckByMirror(v)
		if err != nil {
			logger.L.Errorw(err.Error(), "func", "ProcessCheck()", "extra", "ProcessCheckByMirror(v)")
		}
		v.Status = status
		mirror.SetOneByID(v.Id, v)
	}

	return nil
}

func ProcessCheckByMirror(m mirror.Mirror) (string, error) {
	psCmd := exec.Command("/usr/bin/ps", "-ef")
	var outB, errB bytes.Buffer
	psCmd.Stdout = &outB
	psCmd.Stderr = &errB
	err := psCmd.Run()
	if err != nil {
		logger.L.Error(errB.String())
		return mirror.StatusCheckFailed, fmt.Errorf("psCmd.Run(): %s", err)
	}
	outStr := outB.String()
	split := strings.Split(outStr, "\n")

	for _, currPs := range split {
		if strings.Contains(currPs, "grep --color=auto") || !strings.Contains(currPs, "rsync") || strings.Contains(currPs, "--daemon") {
			continue
		}

		if strings.Contains(currPs, m.Id) {
			return mirror.StatusRunning, nil
		}
	}

	return mirror.StatusNotRunning, nil
}
