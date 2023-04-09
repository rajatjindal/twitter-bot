package twitter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
)

// HandleCRCResponse handles the crc response
func (b *Bot) HandleCRCResponse(w http.ResponseWriter, r *http.Request) {
	crcToken := r.FormValue("crc_token")
	v := computeHmac256(crcToken, b.config.Tokens.ConsumerToken)
	send := fmt.Sprintf(`{"response_token": "sha256=%s"}`, v)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(send))
}

func getConfirmCRCEndpoint(webhookID, environment string) string {
	crcConfirm := fmt.Sprintf("account_activity/all/%s/webhooks/%s.json", environment, webhookID)
	return fmt.Sprintf("%s/%s", twitterAPIBase, crcConfirm)
}

// triggerCRC sends put request to twitter for manually triggering CRC
func (b *Bot) TriggerCRC(webhookID string) error {
	req, err := http.NewRequest("PUT", getConfirmCRCEndpoint(webhookID, b.config.WebhookConfig.Environment), nil)
	if err != nil {
		return err
	}

	resp, err := b.asOwnerOfApp.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
