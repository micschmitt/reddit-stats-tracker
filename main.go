package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// RedditAPI defines the methods that the Reddit client must implement
type RedditAPI interface {
	FetchPosts(subreddit string) ([]*reddit.Post, error)
}

// RedditClient handles communication with Reddit API
type RedditClient struct {
	client *reddit.Client
}

// NewRedditClient creates a new Reddit client
func NewRedditClient(id, secret, username, password string) (*RedditClient, error) {
	client, err := reddit.NewClient(
		reddit.Credentials{
			ID:       id,
			Secret:   secret,
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return nil, err
	}
	return &RedditClient{client: client}, nil
}

// FetchPosts fetches the latest posts from a subreddit
func (c *RedditClient) FetchPosts(subreddit string) ([]*reddit.Post, error) {
	posts, _, err := c.client.Subreddit.NewPosts(context.Background(), subreddit, &reddit.ListOptions{Limit: 100})
	return posts, err
}

// Stats handles tracking statistics for Reddit posts
type Stats struct {
	topPosts  []*reddit.Post
	userPosts map[string]int
	postsCh   chan *reddit.Post
	doneCh    chan bool
	client    RedditAPI
	subreddit string
	mu        sync.Mutex
}

// NewStats initializes a new Stats tracker
func NewStats(client RedditAPI, subreddit string) *Stats {
	return &Stats{
		topPosts:  []*reddit.Post{},
		userPosts: make(map[string]int),
		postsCh:   make(chan *reddit.Post),
		doneCh:    make(chan bool),
		client:    client,
		subreddit: subreddit,
	}
}

// Start begins fetching posts and updating stats
func (s *Stats) Start() {
	go s.fetchPosts()
	go s.updateStats()
}

// Stop halts the stats tracking
func (s *Stats) Stop() {
	close(s.doneCh)
	close(s.postsCh)
}

// fetchPosts continuously fetches posts from Reddit
func (s *Stats) fetchPosts() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			posts, err := s.client.FetchPosts(s.subreddit)
			if err != nil {
				logrus.Errorf("Error fetching posts: %v", err)
				continue
			}
			for _, post := range posts {
				s.postsCh <- post
			}
		case <-s.doneCh:
			return
		}
	}
}

// updateStats processes incoming posts and updates statistics
func (s *Stats) updateStats() {
	for post := range s.postsCh {
		s.mu.Lock()
		s.userPosts[post.Author]++
		s.topPosts = append(s.topPosts, post)

		sort.Slice(s.topPosts, func(i, j int) bool {
			return s.topPosts[i].Score > s.topPosts[j].Score
		})

		if len(s.topPosts) > 10 {
			s.topPosts = s.topPosts[:10]
		}

		s.logStats()
		s.mu.Unlock()
	}
}

// logStats outputs the current statistics
func (s *Stats) logStats() {
	logrus.Infof("Top 10 Posts by Upvotes:")
	for _, post := range s.topPosts {
		logrus.Infof("%s - Upvotes: %d", post.Title, post.Score)
	}

	logrus.Infof("Top Users by Post Count:")
	for user, count := range s.userPosts {
		logrus.Infof("%s - Posts: %d", user, count)
	}
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client, err := NewRedditClient(
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_CLIENT_SECRET"),
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("Failed to create Reddit client: %v", err)
	}

	stats := NewStats(client, "golang")
	stats.Start()

	// Run for 1 minute then stop
	logrus.Infof("Tracking stats for 1 minute...")
	time.Sleep(1 * time.Minute)
	stats.Stop()
}
