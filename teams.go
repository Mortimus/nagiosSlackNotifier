package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// TeamsMessage are buttons in slack
type TeamsMessage struct {
	Type            string                 `json:"@type,omitempty"`
	Context         string                 `json:"@contect"`
	ThemeColor      string                 `json:"themeColor"`
	Summary         string                 `json:"summary"`
	Sections        []TeamsSection         `json:"sections,omitempty"`
	PotentialAction []TeamsPotentialAction `json:"potentialAction,omitempty"`
}

// TeamsSection defines a section in a teams message
type TeamsSection struct {
	ActivityTitle    string       `json:"activityTitle"`
	ActivitySubtitle string       `json:"activitySubtitle"`
	Facts            []TeamsFacts `json:"facts,omitempty"`
	Markdown         bool         `json:"markdown"`
}

// TeamsFacts defines a value-key pair in a section
type TeamsFacts struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TeamsPotentialAction defines a potential action in a teams message
type TeamsPotentialAction struct {
	Type    string         `json:"@type"`
	Name    string         `json:"name"`
	Targets []TeamsTargets `json:"targets"`
}

// TeamsTargets defines links based on target OS Supported operating system values are default, windows, iOS and android
type TeamsTargets struct {
	OS  string `json:"os"`
	URI string `json:"uri"`
}

func (m TeamsMessage) send() {
	// out, err := json.Marshal(m)
	out, err := JSONMarshal(m, true)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post(configuration.TeamsHookURL, appType, bytes.NewReader(out))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if configuration.Debug {
		fmt.Println(resp.Status)
		fmt.Println(string(out))
	}
}

func (m *TeamsMessage) alert(nagios NagiosAlert) {
	m.Type = "MessageCard"
	m.Context = "http://schema.org/extensions"
	m.ThemeColor = nagios.getUrgencyColor()
	m.Summary = fmt.Sprintf("%s (%s) - %s - %s - %s \n %s", nagios.HostAlias, nagios.HostAddress, nagios.ServiceDesc, nagios.NotificationType, nagios.LongDateTime, nagios.ServiceOutput)
	m.Sections = []TeamsSection{
		TeamsSection{
			ActivityTitle:    nagios.HostAlias,
			ActivitySubtitle: nagios.HostAddress,
			Facts: []TeamsFacts{
				TeamsFacts{
					Name:  "TYPE",
					Value: nagios.NotificationType,
				},
				TeamsFacts{
					Name:  "OUTPUT",
					Value: nagios.ServiceOutput,
				},
				TeamsFacts{
					Name:  "SOURCE",
					Value: configuration.TeamsSource,
				},
			},
			Markdown: true,
		},
	}
	if nagios.getUrgencyColor() == configuration.ProblemColor {
		m.PotentialAction = []TeamsPotentialAction{
			TeamsPotentialAction{
				Type: "OpenUri",
				Name: "Acknowledge",
				Targets: []TeamsTargets{
					TeamsTargets{
						OS:  "default",
						URI: m.genACKurl(nagios),
					},
				},
			},
		}
	}
}

func (m *TeamsMessage) genACKurl(nagios NagiosAlert) string {
	var ackURL string
	if isServiceMode() {
		ackURL = configuration.NagiosAckURL + nagiosACKservice + "&host=" + url.QueryEscape(nagios.HostAlias) + "&service=" + url.QueryEscape(nagios.ServiceDesc)
	} else {
		ackURL = configuration.NagiosAckURL + nagiosACKhost + "&host=" + url.QueryEscape(nagios.HostAlias)
	}
	if configuration.Debug {
		fmt.Println("ackURL:", ackURL)
	}
	return ackURL
}
