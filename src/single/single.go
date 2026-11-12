package single

import (
	"os"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

type Script struct {
	pid string
	fp  *os.File
}

func Single() *Script {
	return &Script{
		pid: utils.PIDPath,
		fp:  nil,
	}
}

func (s *Script) CreatePidFile() {
	f, err := os.OpenFile(s.pid, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	s.fp = f
}

func (s *Script) RemovePidFile() {
	err := s.fp.Close()
	if err != nil {
		panic(err)
	}

	err = os.Remove(s.pid)
	if err != nil {
		panic(err)
	}
	s.fp = nil
}

func (s *Script) IsRunning() bool {
	_, err := os.Stat(s.pid)
	if err != nil {
		return false
	}
	return true
}
