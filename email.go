package sendlayer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

type EmailsService struct {
	client *Client
}

func NewEmailsService(client *Client) *EmailsService {
	return &EmailsService{client: client}
}

func (e *EmailsService) validateEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func (e *EmailsService) readAttachment(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		resp, err := e.client.HTTP.Get(filePath)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(content), nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func (e *EmailsService) normalizeRecipient(val interface{}) ([]EmailAddress, error) {
	var out []EmailAddress
	switch v := val.(type) {
	case string:
		if !e.validateEmail(v) {
			return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid email: %s", v)}}
		}
		out = append(out, EmailAddress{Email: v})
	case EmailAddress:
		if !e.validateEmail(v.Email) {
			return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid email: %s", v.Email)}}
		}
		out = append(out, v)
	case []string:
		for _, s := range v {
			if !e.validateEmail(s) {
				return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid email: %s", s)}}
			}
			out = append(out, EmailAddress{Email: s})
		}
	case []EmailAddress:
		for _, addr := range v {
			if !e.validateEmail(addr.Email) {
				return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid email: %s", addr.Email)}}
			}
			out = append(out, addr)
		}
	case map[string]string:
		email, ok := v["email"]
		if !ok || !e.validateEmail(email) {
			return nil, &SendLayerValidationError{SendLayerError{"Invalid email in map"}}
		}
		out = append(out, EmailAddress{Email: email, Name: v["name"]})
	case []map[string]string:
		for _, m := range v {
			email, ok := m["email"]
			if !ok || !e.validateEmail(email) {
				return nil, &SendLayerValidationError{SendLayerError{"Invalid email in map slice"}}
			}
			out = append(out, EmailAddress{Email: email, Name: m["name"]})
		}
	default:
		return nil, &SendLayerValidationError{SendLayerError{"Unsupported recipient type"}}
	}
	return out, nil
}

func (e *EmailsService) Send(
	from interface{},
	to interface{},
	subject string,
	text string,
	html string,
	cc interface{},
	bcc interface{},
	replyTo interface{},
	attachments []Attachment,
	headers map[string]string,
	tags []string,
) (*EmailResponse, error) {
	if text == "" && html == "" {
		return nil, &SendLayerValidationError{SendLayerError{"Either 'text' or 'html' content must be provided."}}
	}
	fromDetails, err := e.normalizeRecipient(from)
	if err != nil || len(fromDetails) == 0 {
		return nil, err
	}
	toList, err := e.normalizeRecipient(to)
	if err != nil || len(toList) == 0 {
		return nil, err
	}
	payload := EmailRequest{
		From:        fromDetails[0],
		To:          toList,
		Subject:     subject,
		ContentType: "Text",
	}
	if html != "" {
		payload.ContentType = "HTML"
		payload.HTMLContent = html
	} else {
		payload.PlainContent = text
	}
	if cc != nil {
		ccList, err := e.normalizeRecipient(cc)
		if err != nil {
			return nil, err
		}
		payload.CC = ccList
	}
	if bcc != nil {
		bccList, err := e.normalizeRecipient(bcc)
		if err != nil {
			return nil, err
		}
		payload.BCC = bccList
	}
	if replyTo != nil {
		replyToList, err := e.normalizeRecipient(replyTo)
		if err != nil {
			return nil, err
		}
		payload.ReplyTo = replyToList
	}
	if len(attachments) > 0 {
		for i, att := range attachments {
			if att.Path == "" || att.Type == "" {
				return nil, &SendLayerValidationError{SendLayerError{"Attachment path and type are required"}}
			}
			content, err := e.readAttachment(att.Path)
			if err != nil {
				return nil, err
			}
			filename := filepath.Base(att.Path)
			payload.Attachments = append(payload.Attachments, struct {
				Content     string `json:"Content"`
				Type        string `json:"Type"`
				Filename    string `json:"Filename"`
				Disposition string `json:"Disposition"`
				ContentId   int    `json:"ContentId"`
			}{
				Content:     content,
				Type:        att.Type,
				Filename:    filename,
				Disposition: "attachment",
				ContentId:   i + 1,
			})
		}
	}
	if headers != nil {
		payload.Headers = headers
	}
	if tags != nil {
		payload.Tags = tags
	}
	respBody, _, err := e.client.doRequest("POST", "email", payload, nil)
	if err != nil {
		return nil, err
	}
	var resp EmailResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
