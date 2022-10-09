package deeplclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	documentTranslateFunctionUri          = "document"
	documentTranslateResultFunctionSubUri = "result"
)

// DocumentTranslationStartRequest contains the payload data for each document translation request
type DocumentTranslationStartRequest struct {
	// SourceLang is the language of the text to be translated. If parameter is omitted, the API will detect the
	// language of the text and translate it.
	SourceLang ApiLang

	// TargetLang determines the language into which you want to translate.
	TargetLang ApiLang

	// File is the file to be translated (only .docx, .pptx, .pdf, .htm(l) and .txt are allowed)
	File []byte

	// Filename describes the name of the file to be translated
	Filename string

	// Formality Sets whether the translated text should lean towards formal or informal language. This feature
	// currently only works for target languages DE (German), FR (French), IT (Italian), ES (Spanish), NL (Dutch),
	// PL (Polish), PT-PT, PT-BR (Portuguese) and RU (Russian).
	Formality ApiFormality

	// GlossaryId A unique ID assigned to a glossary.
	GlossaryId ApiLang
}

// DocumentTranslationStartResponse represents the data of the json response of the document translation API function
type DocumentTranslationStartResponse struct {
	// DocumentId is a unique ID which is being used in subsequent API requests regarding the uploaded document
	DocumentId string `json:"document_id"`

	// DocumentKey is an encryption key which is necessary to decrypt this document on download
	DocumentKey string `json:"document_key"`
}

// DocumentTranslationStatusRequest and DocumentTranslationStartResponse share the same fields
type DocumentTranslationStatusRequest DocumentTranslationStartResponse

// DocumentTranslationStatus is a wrapper type for various states that can occur during translation
type DocumentTranslationStatus string

const (
	StatusQueued      = DocumentTranslationStatus("queued")
	StatusTranslating = DocumentTranslationStatus("translating")
	StatusDone        = DocumentTranslationStatus("done")
	StatusError       = DocumentTranslationStatus("error")
)

// DocumentTranslationStatusResponse represents the data of the json response of the document translation status API
// function
type DocumentTranslationStatusResponse struct {
	// DocumentID is the unique ID of the document
	DocumentId string

	// Status defines the current state of the translation process
	Status DocumentTranslationStatus

	// SecondsRemaining describes the estimated time until the translation is done
	SecondsRemaining uint

	// BilledCharacters is the amount of characters billed
	BilledCharacters uint

	// ErrorMessage describes an error during translation, if one occurred (if not the value is nil)
	ErrorMessage *string
}

// DocumentTranslationDownloadRequest and DocumentTranslationStartResponse share the same fields
type DocumentTranslationDownloadRequest DocumentTranslationStartResponse

// StartDocumentTranslate starts the translation process of a given document and returns document information or an
// error if something went wrong
func (client *Client) StartDocumentTranslate(req *DocumentTranslationStartRequest) (
	resp *DocumentTranslationStartResponse, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if len(req.File) == 0 {
		return resp, errors.New("'File' field most not be empty")
	}
	if len(req.Filename) == 0 {
		return resp, errors.New("'Filename' field must not be empty")
	}
	filePart, err := writer.CreateFormFile("file", req.Filename)
	if err != nil {
		return resp, err
	}

	if _, err = filePart.Write(req.File); err != nil {
		return resp, err
	}

	if req.SourceLang != "" {
		if err = writer.WriteField("source_lang", req.SourceLang.String()); err != nil {
			return resp, err
		}
	}
	if len(req.TargetLang) == 0 {
		return resp, errors.New("'TargetLang' field of translation request cannot be omitted")
	}
	if req.SourceLang != "" {
		if err = writer.WriteField("target_lang", req.TargetLang.String()); err != nil {
			return resp, err
		}
	}
	if len(req.Formality) != 0 {
		if err = writer.WriteField("formality", string(req.Formality)); err != nil {
			return resp, err
		}
	}
	if len(req.GlossaryId) > 0 {
		if err = writer.WriteField("glossary_id", string(req.GlossaryId)); err != nil {
			return resp, err
		}
	}
	// without closing the writer, the body cannot be read
	writer.Close()

	var httpResp *http.Response
	httpResp, err = client.doApiFunctionWithMultipartForm(documentTranslateFunctionUri, http.MethodPost,
		writer.Boundary(), body)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	resp = &DocumentTranslationStartResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return
}

// CheckDocumentTranslationStatus returns the current status of the document translation or an error, if something
// went wrong
func (client *Client) CheckDocumentTranslationStatus(req *DocumentTranslationStatusRequest) (
	resp *DocumentTranslationStatusResponse, err error) {
	values := &url.Values{}

	if len(strings.TrimSpace(req.DocumentId)) == 0 {
		return resp, errors.New("'DocumentId' must not be empty")
	}

	if len(strings.TrimSpace(req.DocumentKey)) == 0 {
		return resp, errors.New("'DocumentKey' must not be empty")
	}
	values.Add("document_key", req.DocumentKey)

	var httpResp *http.Response
	httpResp, err = client.doApiFunction(documentTranslateFunctionUri+"/"+req.DocumentId, http.MethodPost, values)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	resp = &DocumentTranslationStatusResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		if resp.ErrorMessage != nil {
			err = errors.New(*resp.ErrorMessage)
		} else {
			// this may happen if the filename does not match the file content (e.g. txt-filename when file is a
			// Microsoft Word file)
			if resp.Status == StatusError {
				err = errors.New("an unspecified error occurred during translation")
			}
		}
	}
	return
}

// DownloadTranslatedDocument returns the translated document or an error, if something went wrong
func (client *Client) DownloadTranslatedDocument(req *DocumentTranslationDownloadRequest) (result []byte, err error) {
	values := &url.Values{}

	if len(strings.TrimSpace(req.DocumentId)) == 0 {
		return nil, errors.New("'DocumentId' must not be empty")
	}

	if len(strings.TrimSpace(req.DocumentKey)) == 0 {
		return nil, errors.New("'DocumentKey' must not be empty")
	}
	values.Add("document_key", req.DocumentKey)

	var httpResp *http.Response
	httpResp, err = client.doApiFunction(
		documentTranslateFunctionUri+"/"+req.DocumentId+"/"+documentTranslateResultFunctionSubUri,
		http.MethodPost, values)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	result, err = io.ReadAll(httpResp.Body)
	return
}
