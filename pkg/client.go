package deeplclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// StatusQuotaExceeded is the unofficial internal HTTP status code for "Quota exceeded"
	StatusQuotaExceeded  = 456
	translateFunctionUri = "translate"
	usageFunctionUri     = "usage"
	// maximum body size - see https://www.deepl.com/docs-api/accessing-the-api/limits/
	maxBodySize = 128 * 1024
)

// Client allows easy access to the DeepL API by providing methods for each API function.
type Client struct {
	*http.Client
	// AuthKey stores the authentication key required to get access DeepL's API.
	AuthKey     []byte
	EndpointUrl string
}

// handleApiError is an internally used function to parse the status of a finished HTTP request. If any error occurred
// (detected by status code or non-valid JSON response), it will be parsed into a known client API error
func handleApiError(resp *http.Response) (returnResponse bool, err error) {
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusBadRequest:
		err = &WrongRequestErr{}
		break
	case http.StatusForbidden:
		err = &AuthFailedErr{}
		return true, nil
	case http.StatusRequestEntityTooLarge:
		err = &RequestEntityTooLargeErr{}
		break
	case http.StatusTooManyRequests:
		err = &TooManyRequestsErr{}
		break
	case http.StatusNotFound:
		err = &NotFoundErr{}
		break
	case StatusQuotaExceeded:
		err = &QuotaExceededErr{}
		break
	default:
		err = UnwrappedApiResponseCodeErr(resp.StatusCode)
		return false, err
	}
	if jsonErr := json.NewDecoder(resp.Body).Decode(err); jsonErr != nil {
		return false, err
	}

	return true, nil
}

// doApiFunctionWithMultipartForm is an internally used function to execute API functions which require upload
// of complex structures like files etc. The param uri should not begin with a slash character
func (client *Client) doApiFunctionWithMultipartForm(uri, method string, boundary string, body *bytes.Buffer) (
	resp *http.Response, err error) {
	// create new http request
	var req *http.Request
	requestUrl := fmt.Sprintf("%s%s", client.EndpointUrl, uri)
	if req, err = http.NewRequest(method, requestUrl, body); err != nil {
		return
	}
	// add header to allow the server to identify the POST request and auth key
	req.Header.Set("Authorization", "DeepL-Auth-Key "+string(client.AuthKey))
	req.Header.Set("Content-Type", `multipart/form-data; boundary="`+boundary+`"`)

	if resp, err = client.Do(req); err != nil {
		return nil, err
	}

	// check status code and wrap response
	returnResponse, err := handleApiError(resp)
	// in case response is not valid/confusing, omit it
	if !returnResponse {
		resp = nil
	}
	return
}

// doApiFunction is an internally used function to execute API functions more easily. The param uri should not begin with
// a slash character.
func (client *Client) doApiFunction(uri, method string, values *url.Values) (resp *http.Response, err error) {
	// create new http request
	var req *http.Request
	var requestUrl string
	var body io.Reader
	if method == http.MethodPost {
		requestUrl = fmt.Sprintf("%s%s", client.EndpointUrl, uri)
		valuesEncoded := values.Encode()
		if len(valuesEncoded) > maxBodySize {
			return nil, errors.New("body size should not exceed maximum of " + strconv.Itoa(maxBodySize))
		}
		body = strings.NewReader(valuesEncoded)
	} else {
		requestUrl = fmt.Sprintf("%s%s?%s", client.EndpointUrl, uri, values.Encode())
	}
	if req, err = http.NewRequest(method, requestUrl, body); err != nil {
		return
	}
	// add header to allow the server to identify the POST request and auth key
	req.Header.Set("Authorization", "DeepL-Auth-Key "+string(client.AuthKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if resp, err = client.Do(req); err != nil {
		return nil, err
	}

	// check status code and wrap response
	returnResponse, err := handleApiError(resp)
	// in case response is not valid/confusing, omit it
	if !returnResponse {
		resp = nil
	}
	return
}

// ApiLang is a wrapper type for languages used in requests/responses within the translation function.
type ApiLang string

const (
	// 	LangBG Bulgarian
	LangBG = ApiLang("BG")
	// LangCS Czech
	LangCS = ApiLang("CS")
	// LangDA Danish
	LangDA = ApiLang("DA")
	// LangDE German
	LangDE = ApiLang("DE")
	// LangEL Greek
	LangEL = ApiLang("EL")
	// LangEN English
	LangEN = ApiLang("EN")
	// LangES Spanish
	LangES = ApiLang("ES")
	// LangET Estonian
	LangET = ApiLang("ET")
	// LangFI Finnish
	LangFI = ApiLang("FI")
	// LangFR French
	LangFR = ApiLang("FR")
	// LangHU Hungarian
	LangHU = ApiLang("HU")
	// LangID Indonesian
	LangID = ApiLang("ID")
	// LangIT Italian
	LangIT = ApiLang("IT")
	// LangJA Japanese
	LangJA = ApiLang("JA")
	// LangLT Lithuanian
	LangLT = ApiLang("LT")
	// LangLV Latvian
	LangLV = ApiLang("LV")
	// LangNL Dutch
	LangNL = ApiLang("NL")
	// LangPL Polish
	LangPL = ApiLang("PL")
	// LangPT Portuguese (all Portuguese varieties mixed)
	LangPT = ApiLang("PT")
	// LangRO Romanian
	LangRO = ApiLang("RO")
	// LangRU Russian
	LangRU = ApiLang("RU")
	// LangSK Slovak
	LangSK = ApiLang("SK")
	// LangSL Slovenian
	LangSL = ApiLang("SL")
	// LangSV Swedish
	LangSV = ApiLang("SV")
	// LangTR Turkish
	LangTR = ApiLang("TR")
	// LangUK Ukrainian
	LangUK = ApiLang("UK")
	// LangZH Chinese
	LangZH = ApiLang("ZH")
)

// String returns the very basic string representation of the API language.
func (apiLang ApiLang) String() string {
	return string(apiLang)
}

// LangFromString tries to find and return the matching wrapped API language type.
func LangFromString(apiLangString string) (error, ApiLang) {
	switch apiLangString {
	case "BG":
		return nil, LangBG
	case "CS":
		return nil, LangCS
	case "DA":
		return nil, LangDA
	case "DE":
		return nil, LangDE
	case "EL":
		return nil, LangEL
	case "EN":
		return nil, LangEN
	case "ES":
		return nil, LangES
	case "ET":
		return nil, LangET
	case "FI":
		return nil, LangFI
	case "FR":
		return nil, LangFR
	case "HU":
		return nil, LangHU
	case "ID":
		return nil, LangID
	case "IT":
		return nil, LangIT
	case "JA":
		return nil, LangJA
	case "LT":
		return nil, LangLT
	case "LV":
		return nil, LangLV
	case "NL":
		return nil, LangNL
	case "PL":
		return nil, LangPL
	case "PT":
		return nil, LangPT
	case "RO":
		return nil, LangRO
	case "RU":
		return nil, LangRU
	case "SK":
		return nil, LangSK
	case "SL":
		return nil, LangSL
	case "SV":
		return nil, LangSV
	case "TR":
		return nil, LangTR
	case "UK":
		return nil, LangUK
	case "ZH":
		return nil, LangZH
	default:
		return fmt.Errorf("could not find API language: %s", apiLangString), ""
	}
}

// ApiFormality is used to set whether the translation should lean towards formal or informal language.
type ApiFormality string

const (
	// FormalityDefault is the default used.
	FormalityDefault = ApiFormality("default")
	// FormalityMore sets a more formal language.
	FormalityMore = ApiFormality("more")
	// FormalityLess sets a more informal language.
	FormalityLess = ApiFormality("more")
)

// TranslationRequest contains the payload data for each translation request.
type TranslationRequest struct {
	// field comments partially taken from official DeepL API documentation

	// Text to be translated. Only UTF8-encoded plain text is supported. The parameter may be specified multiple times
	// and translations are returned in the same order as in the input. Each of the parameter values may contain
	// multiple sentences.
	Text string
	// SourceLang is the language of the text to be translated. If parameter is omitted, the API will detect the
	// language of the text and translate it.
	SourceLang ApiLang
	// TargetLang determines the language into which you want to translate.
	TargetLang ApiLang
	// TagHandling sets which kind of tags should be handled. Comma-separated list of one or more values. See official
	// DeepL API documentation for more details about tag handling.
	TagHandling []string
	// NonSplittingTags contains a list of XML tags which never split sentences. See official DeepL API documentation
	// for more details about tag handling.
	NonSplittingTags []string
	// IgnoreTags contains a list of XML tags whose content is never translated. See official DeepL API documentation
	// for more details about tag handling.
	IgnoreTags []string
	// DoNotSplitSentences sets whether the translation engine should first split the input into sentences. False by
	// default.
	//
	// For applications which are sending one sentence per text parameter, it is advisable to set this to false, in
	// order to prevent the engine from splitting the sentence unintentionally.
	DoNotSplitSentences bool
	// Sets whether the translation engine should preserve some aspects of the formatting, even if it would usually
	// correct some aspects. False by default.
	//
	// The formatting aspects controlled by the setting include:
	// 	punctuation at the beginning and end of the sentence,
	// 	upper/lower case at the beginning of the sentence.
	PreserveFormatting bool
	// Formality Sets whether the translated text should lean towards formal or informal language. This feature
	// currently only works for target languages DE (German), FR (French), IT (Italian), ES (Spanish), NL (Dutch),
	// PL (Polish), PT-PT, PT-BR (Portuguese) and RU (Russian).
	Formality ApiFormality
}

// TranslationResponse represents the data of the json response of the translation API function.
type TranslationResponse struct {
	// Translations contains all requested translations and their results.
	Translations []struct {
		// DetectedSourceLanguage contains the ApiLang detected by the DeepL API.
		DetectedSourceLanguage ApiLang `json:"detected_source_language"`
		// Text contains the translated text.
		Text string `json:"text"`
	} `json:"translations"`
}

// Translate translates the requested text and returns the translated text or an error if something went wrong.
func (client *Client) Translate(req *TranslationRequest) (resp *TranslationResponse, err error) {
	// parse url values for HTTP request
	values := &url.Values{}
	if len(req.Text) == 0 {
		return resp, errors.New("'Text' field of translation request cannot be empty")
	}
	values.Add("text", req.Text)
	if req.SourceLang != "" {
		values.Add("source_lang", req.SourceLang.String())
	}
	if len(req.TargetLang) == 0 {
		return resp, errors.New("'TargetLang' field of translation request cannot be omitted")
	}
	values.Add("target_lang", req.TargetLang.String())
	if len(req.TagHandling) > 0 {
		values.Add("tag_handling", strings.Join(req.TagHandling, ","))
	}
	if len(req.NonSplittingTags) > 0 {
		values.Add("non_splitting_tags", strings.Join(req.NonSplittingTags, ","))
	}
	if len(req.IgnoreTags) > 0 {
		values.Add("ignore_tags", strings.Join(req.IgnoreTags, ","))
	}
	// inverted functionality from DeepL in order to be able to use the default values of Go
	if req.DoNotSplitSentences {
		values.Add("split_sentences", "0")
	}
	// do not get confused with the different handling for both booleans
	if req.PreserveFormatting {
		values.Add("preserve_formatting", "1")
	}
	if len(req.Formality) > 0 {
		values.Add("formality", string(req.Formality))
	}
	var httpResp *http.Response
	httpResp, err = client.doApiFunction(translateFunctionUri, http.MethodPost, values)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	resp = &TranslationResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return
}

// UsageResponse represents the data of the json response of the usage API function.
type UsageResponse struct {
	// CharacterCount contains the amount of characters translated so far in the current billing period.
	CharacterCount int64 `json:"character_count"`
	// CharacterLimit contains the maximum volume of characters that can be translated in the current billing period.
	CharacterLimit int64 `json:"character_limit"`
}

// GetUsage returns the usage information for the current billing period.
func (client *Client) GetUsage() (resp *UsageResponse, err error) {
	// execute api function
	var httpResp *http.Response
	httpResp, err = client.doApiFunction(usageFunctionUri, http.MethodGet, &url.Values{})
	// check for error
	if err != nil {
		return
	}
	// parse answer into UsageResponse struct
	defer httpResp.Body.Close()
	resp = &UsageResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return
}
