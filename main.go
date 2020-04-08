package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

const appType string = "application/json"

const nagiosACKservice string = "34"
const nagiosACKhost string = "33"

const (
	_           = iota // App name+path - IGNORED
	SERVICE int = iota // Service = Host or Service
	NAGIOSHOST
	SLACKCHANNEL
	NOTIFICATIONTYPE     // Notification Type
	SERVICEDESC          // Service Description
	HOSTALIAS            // Host Alias
	HOSTADDRESS          // Host Address
	SERVICESTATE         // Service State
	LONGDATETIME         // Timestamp of issue
	SERVICEOUTPUT        // Unknown
	NOTIFICATIONCOMMENTS // unknown
	MAXINPUT             // For error checking
)

const (
	SERVICEMODE string = "SERVICE" // This is a service not a host
	HOSTMODE           = "HOST"    // This is a host not a service
)

// NagiosAlert contains raw info from nagios
type NagiosAlert struct { // TODO: make these their real types for more integration
	NotificationType     string
	ServiceDesc          string
	HostAlias            string
	HostAddress          string
	ServiceState         string
	LongDateTime         string
	ServiceOutput        string
	NotificationComments string
}

func (n *NagiosAlert) fromArgs(args []string) {
	n.NotificationType = args[NOTIFICATIONTYPE]
	n.ServiceDesc = args[SERVICEDESC]
	n.HostAlias = args[HOSTALIAS]
	n.HostAddress = args[HOSTADDRESS]
	n.ServiceState = args[SERVICESTATE]
	n.LongDateTime = args[LONGDATETIME]
	n.ServiceOutput = args[SERVICEOUTPUT]
	n.NotificationComments = args[NOTIFICATIONCOMMENTS]
}

func (n NagiosAlert) getUrgencyColor() string {
	// good, warning, danger, or hex code
	/*
		PROBLEM	A service or host has just entered (or is still in) a problem state. If this is a service notification, it means the service is either in a WARNING, UNKNOWN or CRITICAL state. If this is a host notification, it means the host is in a DOWN or UNREACHABLE state.
		RECOVERY	A service or host recovery has occurred. If this is a service notification, it means the service has just returned to an OK state. If it is a host notification, it means the host has just returned to an UP state.
		ACKNOWLEDGEMENT	This notification is an acknowledgement notification for a host or service problem. Acknowledgement notifications are initiated via the web interface by contacts for the particular host or service.
		FLAPPINGSTART	The host or service has just started flapping.
		FLAPPINGSTOP	The host or service has just stopped flapping.
		FLAPPINGDISABLED	The host or service has just stopped flapping because flap detection was disabled..
		DOWNTIMESTART	The host or service has just entered a period of scheduled downtime. Future notifications will be supressed.
		DOWNTIMESTOP	The host or service has just exited from a period of scheduled downtime. Notifications about problems can now resume.
		DOWNTIMECANCELLED	The period of scheduled downtime for the host or service was just cancelled. Notifications about problems can now resume.
	*/
	switch n.NotificationType {
	case "PROBLEM":
		return configuration.ProblemColor
	case "RECOVERY":
		return configuration.RecoveryColor
	case "ACKNOWLEDGEMENT":
		return configuration.AcknowledgeColor
	case "FLAPPINGSTART":
		return configuration.FlappingStartColor
	case "FLAPPINGSTOP":
		return configuration.FlappingStopColor
	case "FLAPPINGDISABLED":
		return configuration.FlappingDisabledColor
	case "DOWNTIMESTART":
		return configuration.DowntimeStartColor
	case "DOWNTIMESTOP":
		return configuration.DowntimeStopColor
	case "DOWNTIMECANCELLED":
		return configuration.DowntimeCancelledColor
	default:
		return configuration.DefaultColor
	}
}

func isServiceMode() bool {
	if os.Args[SERVICE] == SERVICEMODE {
		return true
	}
	return false
}

// JSONMarshal see https://stackoverflow.com/questions/24656624/golang-display-character-not-ascii-like-not-0026
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func main() {
	f, err := os.OpenFile(configuration.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	if len(os.Args) < MAXINPUT {
		missing := MAXINPUT - len(os.Args)
		if configuration.Debug {
			log.Println("Not enough arguments, missing " + strconv.Itoa(missing))
		}
		// append missing strings
		for i := 0; i < missing; i++ {
			os.Args = append(os.Args, "")
		}
	}
	if configuration.Debug {
		if os.Args[SERVICE] == "0" {
			log.Println("Host mode")
		} else {
			log.Println("Service mode")
		}
	}
	var alert NagiosAlert
	alert.fromArgs(os.Args)

	if configuration.Debug {
		log.Println(alert)
	}
	if configuration.AlertSlack {
		var slackMessage Message
		slackMessage.alert(alert)
		slackMessage.send()
	}
	if configuration.AlertTeams {
		var teamsMessage TeamsMessage
		teamsMessage.alert(alert)
		teamsMessage.send()
	}

}
