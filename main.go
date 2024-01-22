package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "time"
)

// Config represents the overall configuration including servers and settings.
type Config struct {
    PingInterval    int      `json:"pingInterval"`
    TeamsWebhookURL string   `json:"teamsWebhookURL"`
    Servers         []Server `json:"servers"`
}

// Server represents a server with its address.
type Server struct {
    Address string `json:"address"`
}

// TeamsMessage represents the JSON structure for Teams message.
type TeamsMessage struct {
    Text string `json:"text"`
}

func main() {
    // Load the configuration from servers.json
    config, err := loadConfig("servers.json")
    if err != nil {
        fmt.Printf("Error loading config: %s\n", err)
        return
    }

    for {
        for _, server := range config.Servers {
            // Ping the server
            err := pingServer(server.Address)
            if err != nil {
                fmt.Printf("Server %s is down: %s\n", server.Address, err)
                // Send a message to Teams
                sendMessageToTeams(config.TeamsWebhookURL, fmt.Sprintf("Server %s is down", server.Address))
            } else {
                fmt.Printf("Server %s is up\n", server.Address)
            }
        }
        time.Sleep(time.Duration(config.PingInterval) * time.Second)
    }
}

// loadConfig loads the configuration from a JSON file.
func loadConfig(filename string) (Config, error) {
    var config Config
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return Config{}, err
    }
    err = json.Unmarshal(file, &config)
    if err != nil {
        return Config{}, err
    }
    return config, nil
}

// pingServer pings the given server address.
func pingServer(address string) error {
    cmd := exec.Command("ping", address, "-c 4")
    if err := cmd.Run(); err != nil {
        return err
    }
    return nil
}

// sendMessageToTeams sends a message to a Microsoft Teams channel.
func sendMessageToTeams(webhookURL, message string) {
    teamsMessage := TeamsMessage{Text: message}

    bytesRepresentation, err := json.Marshal(teamsMessage)
    if err != nil {
        fmt.Printf("Could not marshal Teams message: %s\n", err)
        return
    }

    _, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(bytesRepresentation))
    if err != nil {
        fmt.Printf("Could not send message to Teams: %s\n", err)
        return
    }

    fmt.Println("Message sent to Teams")
}
