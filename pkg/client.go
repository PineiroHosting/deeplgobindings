package deeplclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	authKeyParamName = "auth_key"
	// unofficial internal HTTP status code for "Quota exceeded"
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

// doApiFunction is an internally used function to execute API functions more easily. The param uri should not begin with
// a slash character.
func (client *Client) doApiFunction(uri, method string, values *url.Values) (resp *http.Response, err error) {
	// add authentication header and encode values
	values.Set(authKeyParamName, string(client.AuthKey))
	// create new http request
	var req *http.Request
	var requestUrl string
	var body io.Reader
	if method == http.MethodPost {
		requestUrl = fmt.Sprintf("%s%s", client.EndpointUrl, uri)
		valuesEncoded := values.Encode()
		if valuesEncoded > maxBodySize {
			return nil, errors.New("body size should not exceed maximum of " + maxBodySize)
		}
		body = strings.NewReader(valuesEncoded)
	} else {
		requestUrl = fmt.Sprintf("%s%s?%s", client.EndpointUrl, uri, values.Encode())
	}
	if req, err = http.NewRequest(method, requestUrl, body); err != nil {
		return
	}
	// add header to allow the server to identify the POST request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	// check status code and wrap response
	switch resp.StatusCode {
	case http.StatusOK:
		return
	case http.StatusBadRequest:
		err = &WrongRequestErr{}
		break
	case http.StatusForbidden:
		err = &AuthFailedErr{}
		resp.Body.Close()
		return
	case http.StatusRequestEntityTooLarge:
		err = &RequestEntityTooLargeErr{}
		break
	case http.StatusTooManyRequests:
		err = &TooManyRequestsErr{}
		break
	case StatusQuotaExceeded:
		err = &QuotaExceededErr{}
		break
	default:
		err = UnwrappedApiResponseCodeErr(resp.StatusCode)
		resp.Body.Close()
		return nil, err
	}
	if jsonErr := json.NewDecoder(resp.Body).Decode(err); jsonErr != nil {
		resp.Body.Close()
		return nil, jsonErr
	}
	resp.Body.Close()
	return
}

// ApiLang is a wrapper type for languages used in requests/responses within the translate function.
type ApiLang string

const (
	// LangEN English (EN)
	LangEN = ApiLang("EN")
	// LangDE German (DE)
	LangDE = ApiLang("DE")
	// LangFR French (FR)
	LangFR = ApiLang("FR")
	// LangES Spanish (ES)
	LangES = ApiLang("ES")
	// LangIT Italian (IT)
	LangIT = ApiLang("IT")
	// LangNL Dutch (NL)
	LangNL = ApiLang("NL")
	// LangPL Polish (PL)
	LangPL = ApiLang("PL")
)

// String returns the very basic string representation of the API language.
func (apiLang ApiLang) String() string {
	return string(apiLang)
}

// LangFromString tries to find and return the matching wrapped API language type.
func LangFromString(apiLangString string) (error, ApiLang) {
	switch apiLangString {
	case "EN":
		return nil, LangEN
	case "DE":
		return nil, LangDE
	case "FR":
		return nil, LangFR
	case "ES":
		return nil, LangES
	case "IT":
		return nil, LangIT
	case "NL":
		return nil, LangNL
	case "PL":
		return nil, LangPL
	default:
		return fmt.Errorf("could not find API language: %s", apiLangString), ""
	}
}

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
}

// TranslationResponse represents the data of the json response of the translate API function.
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
