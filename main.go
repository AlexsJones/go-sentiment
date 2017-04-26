package main // package main
import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dimiro1/banner"
)

const b string = `
{{ .AnsiColor.Blue }} _____ ____        ____  _____ _      _____  _  _      _____ _      _____
{{ .AnsiColor.Blue }}/  __//  _ \      / ___\/  __// \  /|/__ __\/ \/ \__/|/  __// \  /|/__ __\
{{ .AnsiColor.Blue }}| |  _| / \|_____ |    \|  \  | |\ ||  / \  | || |\/|||  \  | |\ ||  / \
{{ .AnsiColor.Blue }}| |_//| \_/|\____\\___ ||  /_ | | \||  | |  | || |  |||  /_ | | \||  | |
{{ .AnsiColor.Blue }}\____\\____/      \____/\____\\_/  \|  \_/  \_/\_/  \|\____\\_/  \|  \_/
{{ .AnsiColor.Default }}
`

func usage() {
	fmt.Println("Displays general sentiment around tweets that are positive/negative (default setting prints both)")
	fmt.Println("go-sentiment <KEYWORD> [positive|negative]")
	os.Exit(1)
}
func main() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))
	if len(os.Args) < 2 {
		usage()
	}
	var condition string

	if len(os.Args) > 2 {
		condition = os.Args[2]
	}
	clienv := os.Getenv("CLIENT") //LaFliJ9xgAghWIawfhyq46pBK
	if clienv == "" {
		fmt.Println("Missing client key")
		os.Exit(1)
	}
	clikeyenv := os.Getenv("CLIENT_KEY") //kZj4ajOml9qedXOFZYMUqfhRqZW9Y4Sk4mRJCz8TLYotfC5waj
	if clikeyenv == "" {
		fmt.Println("Missing CLIENT_KEY")
		os.Exit(1)
	}
	tokenv := os.Getenv("API") //ugshTlKTFd2dQ2fKf4TyF4qZE2Bw2W5tnNlEZqt
	if tokenv == "" {
		fmt.Println("Missing API")
		os.Exit(1)
	}
	tokkeyenv := os.Getenv("API_KEY") //5O5h22hoG1qY4S3dHSXOZLFcVw8NTjdxLney7Dmk1dtW2
	if tokkeyenv == "" {
		fmt.Println("Missing API_KEY")
		os.Exit(1)
	}
	log.Println(clienv)
	log.Println(clikeyenv)
	log.Println(tokenv)
	log.Println(tokkeyenv)
	config := oauth1.NewConfig(clienv, clikeyenv)
	token := oauth1.NewToken(tokenv, tokkeyenv)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	params := &twitter.StreamFilterParams{
		Track:         []string{os.Args[1]},
		StallWarnings: twitter.Bool(true),
	}
	stream, _ := client.Streams.Filter(params)
	demux := twitter.NewSwitchDemux()

	a := NewAnalyzer()
	a.classifier.Learn([]string{"happy", "love", "fantastic"}, Positive)
	a.classifier.Learn([]string{"attacking", "attack", "racist", "bigot", "false news", ""}, Negative)
	a.classifier.Learn([]string{"tax", "weather", "news"}, Neutral)

	demux.Tweet = func(tweet *twitter.Tweet) {
		r := a.Classify(tweet.Text)

		switch condition {
		case "positive":
			if r == 1 {
				fmt.Println("------------Found positive sentiment------------------")
				fmt.Println(tweet.Text)
				fmt.Println("------------------------------------------------------")
			}
		case "negative":
			if r == -1 {
				fmt.Println("------------Found negative sentiment------------------")
				fmt.Println(tweet.Text)
				fmt.Println("------------------------------------------------------")
			}

		}
	}
	log.Println("Starting...")
	demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	log.Println(<-ch)

	stream.Stop()
}
