package rss

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/google/uuid"
)

// startScraping 启动定时抓取RSS源的任务
// 参数：
//   - db: 数据库查询接口
//   - concurrency: 并发数量
//   - timeBetweenRequest: 请求间隔时间
func StartScraping(
	query *db.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := query.GetNextFeedsToFetch(
			context.Background(),
			int64(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, query, feed)
		}
		wg.Wait()
	}
}

// scrapeFeed 抓取单个feed的内容并保存到数据库
// 参数：
//   - wg: WaitGroup指针，用于同步
//   - db: 数据库查询接口
//   - feed: 要抓取的feed信息
func scrapeFeed(wg *sync.WaitGroup, query *db.Queries, feed db.Feed) {
	defer wg.Done()
	_, err := query.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToRSSFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed from %s: %v\n", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		// 解析发布时间
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing date %s: %v\n", item.PubDate, err)
			continue
		}

		// 创建新的文章记录
		_, err = query.CreatePost(
			context.Background(),
			db.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Url:         item.Link,
				Description: description,
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			if strings.Contains(err.Error(), "重复键违反唯一约束") {
				continue
			}
			log.Printf("Error creating post: %v\n", err)
			continue
		}
	}
	log.Printf("==> 👀 Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Items))
}