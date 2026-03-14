package rss

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validXML = `<?xml version="1.0"?>
<rss>
  <channel>
    <title>Test Feed</title>
    <link>http://example.com</link>
    <description>A test feed</description>
    <item>
      <title>Post One</title>
      <link>http://example.com/1</link>
      <description>First post</description>
      <pubDate>Mon, 01 Jan 2024 00:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

func Test_FetchFeed_ParsesFeedCorrectly(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validXML)
	}))
	defer server.Close()

	// Act
	feed, err := FetchFeed(context.Background(), server.URL)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if feed.Channel.Title != "Test Feed" {
		t.Fatalf("expected title 'Test Feed', got %q", feed.Channel.Title)
	}
	if len(feed.Channel.Item) != 1 {
		t.Fatalf("expected 1 item, got %d", len(feed.Channel.Item))
	}
	if feed.Channel.Item[0].Title != "Post One" {
		t.Fatalf("expected item title 'Post One', got %q", feed.Channel.Item[0].Title)
	}
}

func Test_FetchFeed_SetsUserAgentHeader(t *testing.T) {
	// Arrange
	var gotAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAgent = r.Header.Get("User-Agent")
		fmt.Fprint(w, validXML)
	}))
	defer server.Close()

	// Act
	FetchFeed(context.Background(), server.URL)

	// Assert
	if gotAgent != "gator" {
		t.Fatalf("expected User-Agent 'gator', got %q", gotAgent)
	}
}

func Test_FetchFeed_InvalidURL(t *testing.T) {
	// Arrange
	// (nothing)

	// Act
	_, err := FetchFeed(context.Background(), "://invalid-url")

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func Test_FetchFeed_ServerError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Act
	_, err := FetchFeed(context.Background(), server.URL)

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid XML body, got nil")
	}
}

func Test_FetchFeed_InvalidXML(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not xml at all")
	}))
	defer server.Close()

	// Act
	_, err := FetchFeed(context.Background(), server.URL)

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}
