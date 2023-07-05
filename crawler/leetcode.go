package crawler

import (
	"log"
	"sync"
)

type Crawler struct {
}

func NewCrawler(logger *log.Logger) *Crawler {
	return &Crawler{}
}

func (cr *Crawler) StartCrawling(wg *sync.WaitGroup) {
	defer wg.Done()

}

func (cr *Crawler) Shutdown() {

}
