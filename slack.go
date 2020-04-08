package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Message to be sent to slack as json
type Message struct {
	Text        string        `json:"text"`
	Channel     string        `json:"channel,omitempty"`
	UserName    string        `json:"username,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

// Field for slack messages
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Action are buttons in slack
type Action struct {
	Name    string        `json:"name,omitempty"`
	Type    string        `json:"type"`
	Text    string        `json:"text"`
	URL     string        `json:"url"`
	Value   string        `json:"value,omitempty"`
	Style   string        `json:"style,omitempty"`
	Confirm ConfirmFields `json:"confirm,omitempty"`
}

// ConfirmFields are confirmation boxes in slack
type ConfirmFields struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	OkText      string `json:"ok_text"`
	DismissText string `json:"dismiss_text"`
}

// Attachment allows for fancy Slack messages
type Attachment struct {
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	PreText    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	ImageURL   string   `json:"image_url"`
	Fields     []Field  `json:"fields"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
	Timestamp  int64    `json:"ts"`
	MarkdownIn []string `json:"mrkdwn_in"`
	Actions    []Action `json:"actions"`
}

func (m *Message) init() {
	m.Channel = configuration.SlackChannel
	m.UserName = configuration.SlackUsername
	m.IconURL = configuration.SlackIconURL
}

// returns the fallback, title, and value
func (m *Message) serviceAlert(nagios NagiosAlert) (string, string, string) {
	var fallback string
	var title string
	var value string
	if isServiceMode() {
		fallback = nagios.HostAlias + " (" + nagios.HostAddress + ")" + " - " + nagios.ServiceDesc + " - " + nagios.NotificationType + " - " + nagios.LongDateTime + "\n" +
			nagios.NotificationComments + " " + nagios.ServiceOutput
		title = nagios.HostAlias + " (" + nagios.HostAddress + ")" + " - " + nagios.ServiceDesc + " - " + nagios.NotificationType
		value = nagios.NotificationComments + " " + nagios.ServiceOutput
	} else {
		fallback = nagios.HostAlias + " (" + nagios.HostAddress + ")" + " - " + nagios.NotificationType + " - " + nagios.LongDateTime + "\n" +
			nagios.ServiceOutput
		title = nagios.HostAlias + " (" + nagios.HostAddress + ") - " + nagios.NotificationType
		value = nagios.ServiceOutput
	}
	if configuration.Debug {
		log.Println("Fallback:", fallback)
		log.Println("Title:", title)
		log.Println("Value:", value)
	}
	return fallback, title, value
}

func (m *Message) genACKurl(nagios NagiosAlert) string {
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

func (m *Message) alert(nagios NagiosAlert) {
	m.init()
	var attach Attachment
	var fld Field
	attach.Color = nagios.getUrgencyColor()
	fld.Short = false
	attach.Fallback, fld.Title, fld.Value = m.serviceAlert(nagios)
	attach.Fields = append(attach.Fields, fld)
	var action Action
	action.Text = "ACKNOWLEDGE"
	action.Type = "button"
	action.URL = m.genACKurl(nagios)
	action.Style = "primary"
	if nagios.getUrgencyColor() == configuration.ProblemColor { // we can only ack bad things
		var conf ConfirmFields
		conf.Title = "Are you sure?"
		conf.Text = "This acknowledges the alarm on Nagios, thus muting notifications until a change occurs."
		conf.OkText = "Yes"
		conf.DismissText = "No"
		action.Confirm = conf
		// attach.Actions = append(attach.Actions, action)
		attach.AddAction(action)
	}

	attach.Footer = configuration.SlackNagiosLink
	m.AddAttachment(&attach)
}

// AddAttachment to a Slack Message
func (m *Message) AddAttachment(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

// AddAction to a Slack Attachment
func (a *Attachment) AddAction(action Action) {
	a.Actions = append(a.Actions, action)
}

func (m Message) send() {
	// out, err := json.Marshal(m)
	out, err := JSONMarshal(m, true)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post(configuration.SlackHookURL, appType, bytes.NewReader(out))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if configuration.Debug {
		fmt.Println(resp.Status)
		fmt.Println(string(out))
	}
}
