package spotibot

import "regexp"

func regexMatch(regEx string, text string) (map[string]string, error) {
	compRegEx, err := regexp.Compile(regEx)
	if err != nil {
		return nil, err
	}
	matches := compRegEx.FindStringSubmatch(text)
	matchMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(matches) {
			matchMap[name] = matches[i]
		}
	}
	return matchMap, nil
}
