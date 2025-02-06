package service

import (
    "fmt"
	"bytes"
	"strings"
	"os/exec"
	"gopkg.in/ini.v1"
)

type ServiceState string

const (
	Fail2banRunning ServiceState = "running"
	Fail2banStop    ServiceState = "stop"
	Fail2banError   ServiceState = "error"
)

type Fail2banStatus struct {
    Installed   bool            `json:"installed"`
    State       ServiceState    `json:"state"`
    Error       string          `json:"error"`
}

type Fail2banService struct {}

func (s *Fail2banService) GetStatus() (*Fail2banStatus, error) {
    status := &Fail2banStatus{}
    cmd := exec.Command("systemctl", "status", "fail2ban")
    var outb bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb
    err := cmd.Run()

    if err != nil {
        if strings.Contains(errb.String(), "service could not be found") {
            return status, nil
        }
        status.Error = errb.String()
        status.State = Fail2banError
        return status, nil
    }
    // TODO: Check if the service is running
// 	var lines []string
//     lines = strings.Split(outb.String(), "\n")
//     fmt.Println(outb.String())
//     fmt.Println(lines)

	return status, nil
}

func (s *Fail2banService) InstallService() error {
    cmd := exec.Command("/bin/bash", "./x-ui.sh", "dev_install_ip_limit")
    var outb bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb

    err := cmd.Run()

    if err != nil {
        fmt.Println(errb.String())
        fmt.Println(err)
        return err
    }
    fmt.Println(outb.String())
    return nil
}


type Fail2banConfig struct {
    BanTime     string     `json:"banTime"`
    LogPath     string     `json:"logPath"`
}

func (s *Fail2banService) GetConfig() (*Fail2banConfig, error) {
    config := &Fail2banConfig{}
    configData, err := ini.Load("./config/fail2ban/jail.d/3x-ipl.conf")

    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        return nil, err
    }

    section := configData.Section("3x-ipl")
    config.BanTime = section.Key("bantime").String()
    config.LogPath = section.Key("logpath").String()

    return config, nil
}