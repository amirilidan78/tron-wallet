package test

import (
	"testing"
	"time"
)

// Crawler test
func TestCrawler(t *testing.T) {
	c := crawler()

	time.Sleep(time.Second * 5)

	res, err := c.ScanBlocks(10)
	if err != nil {
		t.Errorf("ScanBlocks error was incorect, got: %q, want: %q.", err, "nil")
	}
	if len(res) == 0 {
		t.Errorf("ScanBlocks res was incorect, got: %q, want: %q.", res, "not empty")
	}
}
