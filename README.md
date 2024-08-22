# Reddit Stats Tracker

## What This Application Does

Reddit Stats Tracker is a Golang-based application that tracks real-time statistics from a chosen subreddit. It monitors and logs:

- Top posts by upvotes
- Users with the most posts

The application fetches data in near real-time and processes it concurrently to handle large volumes of posts efficiently. It also respects Reddit's API rate limits.

## Technologies Used

- **Golang**: Main programming language.
- **Reddit API**: For fetching subreddit data.
- **Logrus**: For structured logging.
- **Goroutines and Channels**: For concurrency and efficient data processing.
- **godotenv**: For managing environment variables.

## How to Use It
    
1. **Set Up Environment Variables:** Create a `.env` file in the project root with your Reddit API credentials:
    
    ```bash
    REDDIT_CLIENT_ID=your_client_id 
    REDDIT_CLIENT_SECRET=your_client_secret 
    REDDIT_USERNAME=your_username 
    REDDIT_PASSWORD=your_password
    ```
    
2. **Install Dependencies:**
    
    ```bash
    go mod tidy
    ```
    
3. **Run the Application:**
    
    ```bash
    go run main.go
    ```
    

The application will start fetching and logging statistics from your chosen subreddit.