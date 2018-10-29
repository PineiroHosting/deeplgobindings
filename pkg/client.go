package deeplgobindings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	deeplBaseApiUrl  = "https://api.deepl.com/v2/"
	authKeyparamName = "auth_key"
)

// DeeplClient allows easy access to the DeepL API by providing methods for each API function.
type DeeplClient struct {
	*http.Client
	// AuthKey stores the authentication key required to get access DeepL's API.
	AuthKey []byte
}

// doApiFunction is an internally used function to execute API functions more easily. The param uri should not begin with
// a slash character.
func (client *DeeplClient) doApiFunction(uri string, values *url.Values) (resp *http.Response, err error) {
	// add authentication header and encode values
	values.Set(authKeyparamName, string(client.AuthKey))
	body := strings.NewReader(values.Encode())
	// create new http request
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", deeplBaseApiUrl, uri), body); err != nil {
		return
	}
	req.ContentLength = body.Size()
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
		break
	default:
		err = UnwrappedApiResponseCodeErr(resp.StatusCode)
		resp.Body.Close()
		return nil, err
	}
	if jsonErr := json.NewDecoder(resp.Body).Decode(err); err != nil {
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
	return fmt.Sprintf("API-lang:%s", string(apiLang))
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
	//
	// Text to be translated. Only UTF8-encoded plain text is supported. The parameter may be specified multiple times
	// and translations are returned in the same order as in the input. Each of the parameter values may contain
	// multiple sentences.
	Text string
	// Language of the text to be translated. If parameter is omitted, the API will detect the language of the text and
	// translate it.
	SourceLang ApiLang
	// The language into which you want to translate.
	TargetLang         ApiLang
	// Sets which kind of tags should be handled. Comma-separated list of one or more values. See official DeepL API
	// documentation for more details about tag handling.
	TagHandling        []string
	// Comma-separated list of XML tags which never split sentences. See official DeepL API documentation for more
	// details about tag handling.
	NonSplittingTags   []string
	// Comma-separated list of XML tags whose content is never translated. See official DeepL API documentation for more
	//	// details about tag handling.
	IgnoreTags         []string
	// Sets whether the translation engine should first split the input into sentences. True by default.
	//
	// For applications which are sending one sentence per text parameter, it is advisable to set this to false, in
	// order to prevent the engine from splitting the sentence unintentionally.
	SplitSentences     bool
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
		Text                   string `json:"text"`
	} `json:"translations"`
}
