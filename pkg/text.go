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
	translateFunctionUri = "translate"
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
	// GlossaryId Specify the glossary to use for the translation. Important: This requires the source_lang parameter to
	// be set and the language pair of the glossary has to match the language pair of the request.
	GlossaryId ApiLang
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
	if len(req.GlossaryId) > 0 {
		values.Add("glossary_id", string(req.GlossaryId))
	}
	var httpResp *http.Response
	httpResp, err = client.doApiFunction(translateFunctionUri, http.MethodPost, values)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("could not close response of translation request%e\n", err)
		}
	}(httpResp.Body)
	resp = &TranslationResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return
}
