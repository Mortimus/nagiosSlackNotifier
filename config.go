package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const configPath = "config.json"
const failPath = "/opt/NagiosBot"

var configuration Configuration

// Configuration stores all our user defined variables
type Configuration struct {
	Debug                  bool   `json:"Debug"`                  // Allow debug out
	LogPath                string `json:"LogPath"`                // Where to write logs to
	SlackHookURL           string `json:"SlackHookURL"`           // URL to notify slack
	TeamsHookURL           string `json:"TeamsHookURL"`           // URL to notify teams
	SlackChannel           string `json:"SlackChannel"`           // Channel in slack to notify
	SlackUsername          string `json:"SlackUsername"`          // Known Media Extensions for action
	SlackIconURL           string `json:"SlackIconURL"`           // Known Media Extensions for action
	SlackNagiosLink        string `json:"SlackNagiosLink"`        // Known Media Extensions for action
	NagiosAckURL           string `json:"NagiosAckURL"`           // Known Media Extensions for action
	AlertSlack             bool   `json:"AlertSlack"`             // Send Notifications to Slack
	AlertTeams             bool   `json:"AlertTeams"`             // Send Notifications to Teams
	ProblemColor           string `json:"ProblemColor"`           // Problem Color for teams
	RecoveryColor          string `json:"RecoveryColor"`          // Recovery Color for teams
	AcknowledgeColor       string `json:"AcknowledgeColor"`       // Acknowledge Color for teams
	FlappingStartColor     string `json:"FlappingStartColor"`     // FlappingStart Color for teams
	FlappingStopColor      string `json:"FlappingStopColor"`      // FlappingStop Color for teams
	FlappingDisabledColor  string `json:"FlappingDisabledColor"`  // FlappingDisabled Color for teams
	DowntimeStartColor     string `json:"DowntimeStartColor"`     // DowntimeStart Color for teams
	DowntimeStopColor      string `json:"DowntimeStopColor"`      // DowntimeStop Color for teams
	DowntimeCancelledColor string `json:"DowntimeCancelledColor"` // DowntimeCancelled Color for teams
	DefaultColor           string `json:"DefaultColor"`           // Default Color for teams
	TeamsSource            string `json:"TeamsSource"`            // Source included in teams
}

func init() {
	readConfig()
	log.Printf("Configuration loaded:\n %+v\n", configuration)
}

func readConfig() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	}
	if _, err := os.Stat(dir + "/" + configPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		dir, _ = os.Getwd()
	}
	if _, err := os.Stat(dir + "/" + configPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		dir = failPath
	}
	if _, err := os.Stat(dir + "/" + configPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Fatal(err)
	}
	file, err := os.OpenFile(dir+"/"+configPath, os.O_RDONLY, 0444)
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return err
	}
	return nil
}

func saveConfig() error {
	marshalledConfig, _ := json.MarshalIndent(configuration, "", "\t")
	err := ioutil.WriteFile(configPath, marshalledConfig, 0644)
	if err != nil {
		return err
	}
	log.Printf("Config Saved to %s\n", configPath)
	return nil
}
