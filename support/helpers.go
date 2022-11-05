package support

import (
	"io/ioutil"
	"net/http"
)

type Regions int8

const (
	NA Regions = 0
	EU         = 1
	JP         = 2
	FR         = 3
	DE         = 4
)

func (r Regions) String() string {
	switch r {
	case NA:
		return "North America"
	case EU:
		return "European Union"
	case JP:
		return "Japan"
	case FR:
		return "France"
	case DE:
		return "Germany"
	}
	return ""
}

func (r Regions) URIExtensions() string {
	switch r {
	case NA:
		return "NA"
	case EU:
		return "EU"
	case JP:
		return "JP"
	case FR:
		return "FR"
	case DE:
		return "DE"
	}
	return ""
}

func GetLodestoneBaseURI(r Regions) string {
	if r == NA {
		return "https://na.finalfantasyxiv.com/"
	} else if r == EU {
		return "https://eu.finalfantasyxiv.com/"
	} else if r == JP {
		return "https://jp.finalfantasyxiv.com/"
	} else if r == FR {
		return "https://fr.finalfantasyxiv.com/"
	} else if r == DE {
		return "https://de.finalfantasyxiv.com/"
	}
	return "https://na.finalfantasyxiv.com/"
}

func GetHtmlPage(webPage string) (string, error) {
	resp, err := http.Get(webPage)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func UniqueStringSlice(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
