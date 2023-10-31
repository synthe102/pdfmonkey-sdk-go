package pdfmonkeysdkgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/caarlos0/env/v10"
)

const (
	PdfMonkeyAPIEndpoint = "https://api.pdfmonkey.io/api/v1"
)

type Client struct {
	APIKey   string `env:"PDFMONKEY_API_KEY,notEmpty"`
	Endpoint string `env:"PDFMONKEY_API_ENDPOINT"`
}

type CurrentUser struct {
	ID                  string    `json:"id"`
	AuthToken           string    `json:"auth_token"`
	AvailableDocuments  int       `json:"available_documents"`
	CreatedAt           time.Time `json:"created_at"`
	CurrentPlan         string    `json:"current_plan"`
	CurrentPlanInterval string    `json:"current_plan_interval"`
	DesiredName         string    `json:"desired_name"`
	Email               string    `json:"email"`
	Lang                string    `json:"lang"`
	PayingCustomer      bool      `json:"paying_customer"`
	TrialEndsOn         string    `json:"trial_ends_on"`
	UpdatedAt           time.Time `json:"updated_at"`
	BlockResources      bool      `json:"block_resources"`
	ShareLinks          bool      `json:"share_links"`
}

type Document struct {
	ID                 string           `json:"id,omitempty"`
	AppID              string           `json:"app_id,omitempty"`
	Checksum           string           `json:"checksum,omitempty"`
	CreatedAt          time.Time        `json:"created_at,omitempty"`
	DocumentTemplateID string           `json:"document_template_id,omitempty"`
	DownloadURL        string           `json:"download_url,omitempty"`
	FailureCause       any              `json:"failure_cause,omitempty"`
	Filename           string           `json:"filename,omitempty"`
	GenerationLogs     []GenerationLogs `json:"generation_logs,omitempty"`
	Meta               string           `json:"meta,omitempty"`
	Payload            string           `json:"payload,omitempty"`
	PreviewURL         string           `json:"preview_url,omitempty"`
	PublicShareLink    string           `json:"public_share_link,omitempty"`
	Status             string           `json:"status,omitempty"`
	UpdatedAt          time.Time        `json:"updated_at,omitempty"`
}
type GenerationLogs struct {
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type DocumentCard struct {
	ID                         string    `json:"id,omitempty"`
	AppID                      string    `json:"app_id,omitempty"`
	CreatedAt                  time.Time `json:"created_at,omitempty"`
	DocumentTemplateID         string    `json:"document_template_id,omitempty"`
	DocumentTemplateIdentifier string    `json:"document_template_identifier,omitempty"`
	DownloadURL                string    `json:"download_url,omitempty"`
	FailureCause               any       `json:"failure_cause,omitempty"`
	Filename                   string    `json:"filename,omitempty"`
	Meta                       string    `json:"meta,omitempty"`
	PublicShareLink            string    `json:"public_share_link,omitempty"`
	Status                     string    `json:"status,omitempty"`
	UpdatedAt                  time.Time `json:"updated_at,omitempty"`
}

type DocumentCardList []DocumentCard

type ClientOption func(c *Client)

func NewClient(opts ...ClientOption) (*Client, error) {
	client := &Client{
		Endpoint: PdfMonkeyAPIEndpoint,
	}
	for _, opt := range opts {
		opt(client)
	}
	if err := env.Parse(client); err != nil {
		return nil, err
	}
	return client, nil
}

func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.APIKey = key
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.Endpoint = endpoint
	}
}

type GetCurrentUserResponse struct {
	CurrentUser CurrentUser `json:"current_user"`
}

func (c *Client) GetCurrentUser() (*GetCurrentUserResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.Endpoint+"/current_user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out *GetCurrentUserResponse

	err = json.Unmarshal(b, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type ListDocumentsInput struct {
	Page struct {
		Number int32 `json:"number,omitempty"`
	} `json:"page,omitempty"`
	Q struct {
		DocumentTemplateID []string  `json:"document_template_id,omitempty"`
		WorkspaceID        string    `json:"workspace_id,omitempty"`
		Status             string    `json:"status,omitempty"`
		UpdatedSince       time.Time `json:"update_since,omitempty"`
	} `json:"q,omitempty"`
}

type ListDocumentsOutput struct {
	DocumentCardList DocumentCardList `json:"documents"`
}

func (c *Client) ListDocuments(input *ListDocumentsInput) (*ListDocumentsOutput, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, c.Endpoint+"/documents", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to list documents, error code: %v", res.StatusCode)
	}

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out *ListDocumentsOutput
	err = json.Unmarshal(resb, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type GetDocumentCardOutput struct {
	DocumentCard DocumentCard `json:"document_card"`
}

func (c *Client) GetDocumentCard(documentId *string) (*GetDocumentCardOutput, error) {
	req, err := http.NewRequest(http.MethodGet, c.Endpoint+fmt.Sprintf("/document_cards/%s", *documentId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch document card, error code: %v", res.StatusCode)
	}

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out *GetDocumentCardOutput
	err = json.Unmarshal(resb, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type GetDocumentOutput struct {
	Document Document `json:"document"`
}

func (c *Client) GetDocument(documentId *string) (*GetDocumentOutput, error) {
	req, err := http.NewRequest(http.MethodGet, c.Endpoint+fmt.Sprintf("/documents/%s", *documentId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch documents, error code: %v", res.StatusCode)
	}

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out *GetDocumentOutput
	err = json.Unmarshal(resb, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type CreateDocumentInput struct {
	Document Document `json:"document"`
}

type CreateDocumentOutput struct {
	Document Document `json:"document"`
}

func (c *Client) CreateDocument(input *CreateDocumentInput) (*CreateDocumentOutput, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.Endpoint+"/documents", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create document, error code: %v", res.StatusCode)
	}

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out *CreateDocumentOutput
	err = json.Unmarshal(resb, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type UpdateDocumentInput struct {
	Document Document `json:"document"`
}

func (c *Client) UpdateDocument(input *UpdateDocumentInput) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.Endpoint+fmt.Sprintf("/documents/%s", input.Document.ID), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to edit document with ID: %s", input.Document.ID)
	}

	return nil
}

func (c *Client) DeleteDocument(documentId *string) error {
	req, err := http.NewRequest(http.MethodDelete, c.Endpoint+fmt.Sprintf("/documents/%s", *documentId), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to delete document with ID: %s", *documentId)
	}

	return nil
}
