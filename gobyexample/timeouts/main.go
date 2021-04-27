// See https://gobyexample.com/timeouts

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Quote string

type Api interface {
	RandomQuote() (Quote, error)
}

type QuoteService interface {
	Api() Api
	SendQuote(ch chan<- Quote)
}

type RandomQuoteService struct {}

func (rqs* RandomQuoteService) Api() Api {
	return api
}

func (rqs *RandomQuoteService) SendQuote(ch chan<- Quote) {
	defer close(ch)
	q, err := rqs.Api().RandomQuote()
	if err != nil {
		fmt.Println("Error!", err)
		return
	}
	ch <- q
}

func GenerateQuotes(qs QuoteService, timeout time.Duration, retryOnTimeout int) {
	if retryOnTimeout <= 0 {
		fmt.Println("No more retries. Exiting...")
		return
	}

	ch := make(chan Quote, 1)
	go qs.SendQuote(ch)

	select {
	case q := <-ch:
		fmt.Println("Success! The quote is: ", q)
	case <-time.After(timeout):
		fmt.Printf("Quote request timed-out. Attempting to retry (%v left)\n", retryOnTimeout)
		GenerateQuotes(qs, timeout, retryOnTimeout - 1)
	}
}

func main () {
	rqs := RandomQuoteService {}
	GenerateQuotes(&rqs, time.Second * 3, 2)
}

type FakeApi []Quote

func (fa FakeApi) RandomQuote() (Quote, error) {
	// Simulate delay
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randomDelay := time.Duration(r.Intn(5)) * time.Second 
	time.Sleep(randomDelay)
	
	quoteIndex := r.Intn(len(fa))
	return fa[quoteIndex], nil
}

var api = FakeApi{
	"That brain of mine is something more than merely mortal; as time will show.",
	"If you can’t give me poetry, can’t you give me poetical science?",
	"I never am really satisfied that I understand anything; because, understand it well as I may, my comprehension can only be an infinitesimal fraction of all I want to understand about the many connections and relations which occur to me, how the matter in question was first thought of or arrived at…",
	"Religion to me is science and science is religion.",
	"Your best and wisest refuge from all troubles is in your science.",
}
