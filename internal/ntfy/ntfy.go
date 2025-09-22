package ntfy

import (
	"bytes"
	"fmt"
	"net/http"

	"rainalert/internal/config"
)

func SendRainAlert(cfg config.Config) error {
	return SendNtfy(cfg, "Rain Alert", "5", "Rain expected tomorrow!")
}

func SendNoRainAlert(cfg config.Config) error {
	return SendNtfy(cfg, "No Rain", "3", "No rain expected tomorrow.")
}

func SendErrorAlert(cfg config.Config, errMsg string) error {
	return SendNtfy(cfg, "Error", "3", errMsg)
}

func SendNtfy(cfg config.Config, title string, priority string, message string) error {
	url := fmt.Sprintf("https://ntfy.sh/%s", cfg.NtfyTopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return err
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
