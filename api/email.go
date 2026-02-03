package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"unicode/utf8"
)

type EmailRequest struct {
	Message string `json:"message"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
}

type EmailConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
	ToEmail     string
}

const (
	MaxMessageLength = 10000 // 10KB max for message
	MaxSubjectLength = 200   // Max subject line length
	MaxNameLength    = 100   // Max name length
	MaxPhoneLength   = 20    // Max phone length
	MaxEmailLength   = 254   // RFC 5321 max email length
)

var (
	headerInjectionPattern = regexp.MustCompile(`(?i)[\r\n]+(to|cc|bcc|from|subject):\s*`)
	htmlTagPattern         = regexp.MustCompile(`<[^>]*>`)
	mimeBoundaryPattern    = regexp.MustCompile(`(?i)content-type|content-transfer-encoding|multipart|boundary`)
	emailPattern           = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phonePattern           = regexp.MustCompile(`^[\d\s\(\)\-\+\.]+$`)
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

var emailConfig EmailConfig

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailReq EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := validateEmailRequest(&emailReq); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	if err := sendEmail(emailConfig, emailReq); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
		"status":  "success",
	})
}

func validateEmailRequest(req *EmailRequest) error {
	if req.Email == "" {
		return ValidationError{Field: "email", Message: "email is required"}
	}
	if req.Subject == "" {
		return ValidationError{Field: "subject", Message: "subject is required"}
	}
	if req.Message == "" {
		return ValidationError{Field: "message", Message: "message is required"}
	}

	if len(req.Email) > MaxEmailLength {
		return ValidationError{Field: "email", Message: "email is too long"}
	}
	if len(req.Subject) > MaxSubjectLength {
		return ValidationError{Field: "subject", Message: fmt.Sprintf("subject exceeds %d characters", MaxSubjectLength)}
	}
	if len(req.Message) > MaxMessageLength {
		return ValidationError{Field: "message", Message: fmt.Sprintf("message exceeds %d characters", MaxMessageLength)}
	}
	if len(req.Name) > MaxNameLength {
		return ValidationError{Field: "name", Message: fmt.Sprintf("name exceeds %d characters", MaxNameLength)}
	}
	if len(req.Phone) > MaxPhoneLength {
		return ValidationError{Field: "phone", Message: fmt.Sprintf("phone exceeds %d characters", MaxPhoneLength)}
	}

	if !emailPattern.MatchString(req.Email) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}

	if req.Phone != "" && !phonePattern.MatchString(req.Phone) {
		return ValidationError{Field: "phone", Message: "invalid phone format"}
	}

	if headerInjectionPattern.MatchString(req.Subject) {
		return ValidationError{Field: "subject", Message: "invalid characters detected"}
	}
	if headerInjectionPattern.MatchString(req.Message) {
		return ValidationError{Field: "message", Message: "invalid characters detected"}
	}
	if headerInjectionPattern.MatchString(req.Name) {
		return ValidationError{Field: "name", Message: "invalid characters detected"}
	}

	if mimeBoundaryPattern.MatchString(req.Message) {
		return ValidationError{Field: "message", Message: "suspicious content detected"}
	}
	if mimeBoundaryPattern.MatchString(req.Subject) {
		return ValidationError{Field: "subject", Message: "suspicious content detected"}
	}

	if htmlTagPattern.MatchString(req.Message) {
		return ValidationError{Field: "message", Message: "HTML tags are not allowed"}
	}
	if htmlTagPattern.MatchString(req.Subject) {
		return ValidationError{Field: "subject", Message: "HTML tags are not allowed"}
	}
	if htmlTagPattern.MatchString(req.Name) {
		return ValidationError{Field: "name", Message: "HTML tags are not allowed"}
	}

	if !utf8.ValidString(req.Email) || !utf8.ValidString(req.Subject) ||
		!utf8.ValidString(req.Message) || !utf8.ValidString(req.Name) ||
		!utf8.ValidString(req.Phone) {
		return ValidationError{Field: "general", Message: "invalid character encoding"}
	}

	req.Email = sanitizeString(req.Email)
	req.Subject = sanitizeString(req.Subject)
	req.Message = sanitizeString(req.Message)
	req.Name = sanitizeString(req.Name)
	req.Phone = sanitizeString(req.Phone)

	return nil
}

func sanitizeString(s string) string {
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)

	var result strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' || (r >= 32 && r < 127) || r >= 160 {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func sendEmail(config EmailConfig, req EmailRequest) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	message := composeMessage(config, req)

	addr := config.Host + ":" + config.Port
	err := smtp.SendMail(
		addr,
		auth,
		config.SenderEmail,
		[]string{config.ToEmail},
		[]byte(message),
	)

	return err
}

func composeMessage(config EmailConfig, req EmailRequest) string {
	fromHeader := fmt.Sprintf("\"%s\" <%s>", config.SenderName, config.SenderEmail)
	subjectLine := fmt.Sprintf("AlanGeorge.dev: %s", req.Subject)
	emailBody := fmt.Sprintf(`AlanGeorge.dev Contact Form Submission
=======================

Name: %s
Email: %s
Phone: %s

Message:
--------
%s

=======================
This is a plain text email with no attachments.
Sent via contact form API.
`, req.Name, req.Email, req.Phone, req.Message)

	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		fromHeader,
		config.ToEmail,
		subjectLine,
		emailBody,
	)

	return message
}
