package alexa

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func ValidateAlexaRequest(r *http.Request, rawBody string, body AlexaRequest) error {
	signatureCertChainURL := r.Header.Get("SignatureCertChainUrl")
	signature := r.Header.Get("Signature")
	timestamp := body.Request.Timestamp

	if err := validateURL(signatureCertChainURL); err != nil {
		return err
	}
	if err := validateTimestamp(timestamp); err != nil {
		return err
	}
	certificate, err := downloadCertificate(signatureCertChainURL)
	if err != nil {
		return err
	}
	if err := validateCertificate(certificate); err != nil {
		return err
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	return validateSignature(certificate, signatureBytes, []byte(rawBody))
}

func validateURL(signatureCertChainURL string) error {
	if !strings.HasPrefix(signatureCertChainURL, "https://s3.amazonaws.com/echo.api/") {
		return errors.New("invalid certificate chain URL")
	}
	return nil
}

func validateTimestamp(timestamp string) error {
	requestTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}
	if time.Since(requestTime).Seconds() > 150 {
		return errors.New("request timestamp is too old")
	}
	return nil
}

func downloadCertificate(url string) (*x509.Certificate, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}
	return x509.ParseCertificate(block.Bytes)
}

func validateCertificate(cert *x509.Certificate) error {
	if !strings.Contains(cert.Subject.CommonName, "echo-api.amazon.com") {
		return errors.New("certificate is not issued by Amazon")
	}
	return nil
}

func validateSignature(cert *x509.Certificate, signature, data []byte) error {
	return cert.CheckSignature(x509.SHA1WithRSA, data, signature)
}
