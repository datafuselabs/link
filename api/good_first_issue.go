package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/google/go-github/v44/github"
)

func GoodFirstIssue(w http.ResponseWriter, r *http.Request) {
	client := github.NewClient(nil)

	ctx := context.Background()

	repos := []string{"databend", "openraft", "opendal"}

	wg := &sync.WaitGroup{}
	var issues []*github.Issue
	lock := sync.Mutex{}

	for _, v := range repos {
		wg.Add(1)

		go func(repo string) {
			is, _, err := client.Issues.ListByRepo(ctx, "datafuselabs", repo, &github.IssueListByRepoOptions{Labels: []string{"good first issue"}})
			if err != nil {
				log.Fatalf("ListByOrg: %s", err)
			}

			lock.Lock()
			defer lock.Unlock()
			issues = append(issues, is...)
		}(v)
	}
	wg.Wait()

	index := rand.Intn(len(issues))
	url := *issues[index].URL

	w.Header().Add("Location", url)
	w.WriteHeader(302)
	w.Write(nil)
}
