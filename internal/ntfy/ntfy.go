package ntfy

import (
	"bytes"
	"fmt"
	"net/http"

	"rainalert/internal/config"
)

func SendNtfy(cfg config.Config, message string) error {
	url := fmt.Sprintf("https://ntfy.sh/%s", cfg.NtfyTopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return err
	}
	req.Header.Set("Title", "Weather Alert")
	req.Header.Set("Priority", "5")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
