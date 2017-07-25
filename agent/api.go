package agent

import (
	"encoding/json"
	"net/http"
)

type ApiAgent struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Active  bool   `json:"active"`
	Enabled bool   `json:"enabled"`
	Busy    bool   `json:"busy"`
}

func GetAgentStatus(ag *Agent, username, password string) (*ApiAgent, error) {
	url := ag.FormatBambooURL("rest/api/1.0/agent?online=true&os_auth=basic")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	var aar []ApiAgent
	err = json.NewDecoder(resp.Body).Decode(&aar)
	if err != nil {
		return nil, err
	}

	for _, aa := range aar {
		if aa.ID == ag.ID {
			return &aa, nil
		}
	}

	return nil, nil
}
