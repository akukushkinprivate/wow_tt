package quote

import "math/rand"

var quotes = []string{
	"You create your own opportunities.",
	"Never break your promises.",
	"You are never as stuck as you think you are.",
	"Happiness is a choice.",
	"Asking for help is a sign of strength.",
	"Replace every negative thought with a positive one.",
	"Accept what is, let go of what was, have faith in what will be.",
	"Be confident enough to encourage confidence in others.",
}

func GetRandomQuote() string {
	return quotes[rand.Int63n(int64(len(quotes)))]
}
