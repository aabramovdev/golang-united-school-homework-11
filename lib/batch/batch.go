package batch

import (
	"sort"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

// func getBatch(n int64, pool int64) (res []user) {
// 	res = make([]user, n)
// 	var wg sync.WaitGroup
// 	var m sync.Mutex
// 	sem := make(chan struct{}, pool)
// 	var i int64
// 	for i = 0; i < n; i++ {
// 		wg.Add(1)
// 		sem <- struct{}{}
// 		go func(j int64) {
// 			defer wg.Done()
// 			m.Lock()
// 			res[j] = getOne(j)
// 			m.Unlock()
// time.Sleep(time.Nanosecond * 20)
// <-sem
// wg.Done()
// }(i)
// }
// wg.Wait()
// return
// }

func getBatch(n int64, pool int64) (res []user) {
	jobs := make(chan int64, n)
	results := make(chan user, n)
	var w int64
	for w = 0; w < pool; w++ {
		go worker(w, jobs, results)
	}
	var idx int64
	for idx = 0; idx < n; idx++ {
		jobs <- idx
	}
	close(jobs)

	return sortUsers(chanToUsers(n, results))
}

func worker(id int64, jobs <-chan int64, results chan<- user) {
	for j := range jobs {
		// fmt.Println("Worker  ID:", id)
		results <- getOne(j)
	}
}

func chanToUsers(n int64, results <-chan user) (users []user) {
	var idx int64
	for idx = 0; idx < n; idx++ {
		users = append(users, <-results)
	}
	return
}

func sortUsers(users []user) []user {
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users
}
