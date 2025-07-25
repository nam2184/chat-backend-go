package middleware

import (
  "regexp"
)

var (
  contentSelectorRegexJSON *regexp.Regexp
  contentSelectorRegexAPI *regexp.Regexp 
) 

func init() {
  contentSelectorRegexJSON = regexp.MustCompile(`(?i)^application\/json`)
	contentSelectorRegexAPI = regexp.MustCompile(`(?i)^application\/dem\.api\+json$`)
}
