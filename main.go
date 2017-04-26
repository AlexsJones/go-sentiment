package main // package main
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("requires search parameter...")
		os.Exit(1)
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
		log.Println(tweet.Text)
		r := a.Classify(tweet.Text)
		switch r {
		case 1:
			fmt.Println("------------Found positive sentiment------------------")
			fmt.Println(tweet.Text)
			fmt.Println("------------------------------------------------------")
		case -1:
			fmt.Println("------------Found negative sentiment------------------")
			fmt.Println(tweet.Text)
			fmt.Println("------------------------------------------------------")
		}
	}
	log.Println("Starting...")
	demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	log.Println(<-ch)

	stream.Stop()
}
