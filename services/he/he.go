package he

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nam2184/mymy/models/db"
)

type HEService struct {
	BaseEndpoint string
}

func NewHEService() *HEService {
	return &HEService{
		BaseEndpoint: os.Getenv("HE_API_URL"),
	}
}

func (h *HEService) ClassifyNSFWContent(message db.EncryptedMessage) (bool, error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return false, fmt.Errorf("Error marshaling message to JSON %v:", err)
	}

	resp, err := http.Post(h.BaseEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("Error making request to JSON %v:", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusForbidden:
		return false, nil
	case http.StatusUnprocessableEntity:
		return false, fmt.Errorf("Not processable for some reason")
	default:
		return false, fmt.Errorf("Unexpected status code %d from %s\n", resp.StatusCode, message.Content)
	}
}
