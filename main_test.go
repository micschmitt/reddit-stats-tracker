package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// MockRedditClient mocks the RedditClient for unit testing
type MockRedditClient struct {
	mock.Mock
}

func (m *MockRedditClient) FetchPosts(subreddit string) ([]*reddit.Post, error) {
	args := m.Called(subreddit)
	return args.Get(0).([]*reddit.Post), args.Error(1)
}

func TestFetchPosts(t *testing.T) {
	mockClient := new(MockRedditClient)
	mockClient.On("FetchPosts", "golang").Return([]*reddit.Post{
		{Title: "Post 1", Score: 100},
		{Title: "Post 2", Score: 50},
	}, nil)

	posts, err := mockClient.FetchPosts("golang")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, 100, posts[0].Score)
}

func TestUpdateStats(t *testing.T) {
	mockClient := new(MockRedditClient)
	mockClient.On("FetchPosts", "golang").Return([]*reddit.Post{
		{Title: "Post 1", Author: "user1", Score: 100},
		{Title: "Post 2", Author: "user2", Score: 50},
	}, nil)

	stats := NewStats(mockClient, "golang")

	// Start the stats processing in a goroutine
	go func() {
		stats.postsCh <- &reddit.Post{Title: "Post 1", Author: "user1", Score: 100}
		stats.postsCh <- &reddit.Post{Title: "Post 2", Author: "user2", Score: 50}
		close(stats.postsCh)
	}()

	// Update stats and ensure processing stops when the channel is closed
	stats.updateStats()

	// Assertions
	assert.Equal(t, 2, len(stats.topPosts))
	assert.Equal(t, "user1", stats.topPosts[0].Author)
	assert.Equal(t, 100, stats.topPosts[0].Score)
	assert.Equal(t, 1, stats.userPosts["user1"])
	assert.Equal(t, 1, stats.userPosts["user2"])
}
