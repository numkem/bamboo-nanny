package agent

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"

	"path/filepath"

	"github.com/magiconair/properties"
)

type Agent struct {
	ID               int    `xml:"agentDefinition>id"`
	Name             string `xml:"agentDefinition>name"`
	BambooURL        string
	WorkingDirectory string `xml:"buildWorkingDirectory"`
}

func (a *Agent) FormatBambooURL(url string) string {
	ss := strings.Split(a.BambooURL, "/")
	return fmt.Sprintf("%s//%s/%s", ss[0], ss[2], url)
}

func NewAgentFromDir(dir string) (*Agent, error) {
	// Load configuration
	wrapperFile := filepath.Join(dir, "conf", "wrapper.conf")
	prop, err := properties.LoadFile(wrapperFile, properties.UTF8)
	if err != nil {
		return nil, err
	}
	props := prop.Map()

	var agent Agent
	agent.BambooURL = props["wrapper.app.parameter.2"]

	// Decode the xml configuration
	f, err := os.Open(path.Join(dir, "bamboo-agent.cfg.xml"))
	if err != nil {
		return nil, err
	}

	err = xml.NewDecoder(f).Decode(&agent)
	return &agent, err
}
