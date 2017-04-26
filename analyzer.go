package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jbrukh/bayesian"
)

const (
	//Positive classifier
	Positive bayesian.Class = "Positive"
	//Negative classifier
	Negative bayesian.Class = "Negative"
	//Neutral classifier
	Neutral bayesian.Class = "Neutral"
)

//DATA_FILE dump location
const DATA_FILE = "./sentiment-data/sentiment-classifier.dmp"

//Analyzer struct
type Analyzer struct {
	classifier *bayesian.Classifier
}

//Classify will return sentiment level
func (a *Analyzer) Classify(s string) int {
	if len(s) <= 2 {
		return 0
	}
	tokens := tokenize(s)

	_, likely, _ := a.classifier.LogScores(tokens)

	sentiment := 0
	switch likely {
	case 0:
		sentiment = 1
	case 1:
		sentiment = -1
	case 2:
		sentiment = 0
	}

	return sentiment
}

//NewAnalyzer ...
func NewAnalyzer() Analyzer {
	a := Analyzer{}

	_, err := os.Stat(DATA_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			a.downloadDataSet()
		}
	}

	c, err := bayesian.NewClassifierFromFile(DATA_FILE)
	if err == nil {
		a.classifier = c
	} else {
		// Note: Nothing will be trained at this point, but we'll still have a classifier that can be trained
		a.classifier = bayesian.NewClassifier(Positive, Negative, Neutral)
	}

	return a
}

// Retrieves training data (which is much too large to keep in GitHub)
func (a *Analyzer) downloadDataSet() {
	os.Mkdir("./sentiment-data", 0777)
	out, oErr := os.Create(DATA_FILE)
	defer out.Close()
	if oErr == nil {
		r, rErr := http.Get("https://s3.amazonaws.com/socialharvest/public-data/sentiment/sentiment-classifier.dmp")
		if rErr != nil {
			log.Println(rErr.Error())
			return
		}
		defer r.Body.Close()
		if rErr == nil {
			_, nErr := io.Copy(out, r.Body)
			if nErr != nil {
				err := os.Remove(DATA_FILE)
				if err != nil {
					log.Println(err)
				}
			}
			r.Body.Close()
		} else {
			log.Println(rErr)
		}
		out.Close()
	} else {
		log.Println(oErr)
	}
}

func tokenize(s string) []string {
	tokens := []string{}

	tokenSlice := strings.Split(s, " ")
	for k, v := range tokenSlice {
		tokens = append(tokens, v)
		if len(tokenSlice)-1 > k && len(v) > 1 {
			ngram := tokenSlice[k] + " " + tokenSlice[k+1] //+ " " + tokenSlice[k+2]
			tokens = append(tokens, ngram)
		}
	}

	return tokens
}
