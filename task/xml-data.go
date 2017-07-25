package task

import (
	"io/ioutil"
	"os"
	"regexp"

	"path/filepath"

	"github.com/numkem/bamboo-nanny/agent"
)

const folderRegex = "[A-Z]+-[A-Z,0-9]+-[A-Z,0-9]+"

type XMLDataCleanup struct{}

func (t XMLDataCleanup) Run(a *agent.Agent) error {
	files, err := ioutil.ReadDir(a.WorkingDirectory)
	if err != nil {
		return err
	}

	for _, f := range files {
		match, err := regexp.MatchString(folderRegex, f.Name())
		if err != nil {
			return err
		}

		if match {
			p := filepath.Join(a.WorkingDirectory, f.Name())
			err := os.RemoveAll(p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
