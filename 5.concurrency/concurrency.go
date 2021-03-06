package main

import (
    "fmt"
    "time"
    "golang.org/x/tour/tree"
    "sync"
)

func main() {
	coroutines()

    channels()
    channelsDemo()
    bufferedChannels()

    rangeAndClose()

    selectStatement()
    selectDefault()

    equivalentBinaryTreesExercise()

    mutex()

    webCrawlerExercise()
}

/*
A goroutine is a lightweight thread managed by the Go runtime.
go f(x, y, z)
starts a new goroutine running
f(x, y, z)
The evaluation of f, x, y, and z happens in the current goroutine and the execution of f happens in the new goroutine.
Goroutines run in the same address space, so access to shared memory must be synchronized. The sync package provides useful primitives, although you won't need them much in Go as there are other primitives. (See the next slide.)
*/
func coroutines() {
    fmt.Println("Coroutines----------------")
    go say("world")
	say("hello")
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

/*
Channels are a typed conduit through which you can send and receive values with the channel operator, <-.

ch <- v    // Send v to channel ch.
v := <-ch  // Receive from ch, and
           // assign value to v.
(The data flows in the direction of the arrow.)

Like maps and slices, channels must be created before use:

ch := make(chan int)
By default, sends and receives block until the other side is ready. This allows goroutines to synchronize without explicit locks or condition variables.

The example code sums the numbers in a slice, distributing the work between two goroutines. Once both goroutines have completed their computation, it calculates the final result.
*/
func sum(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
    }
    c <- sum // send sum to c
}

func channels() {
    fmt.Println("Channels----------------")
    s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}

// channels demo

func work(miliseconds int64, channel chan int64) {
    fmt.Println("blocking ", miliseconds)
    sleepMS := time.Duration(miliseconds)*time.Millisecond
    time.Sleep(sleepMS)
    fmt.Println("unblocking ", miliseconds)
    channel <- miliseconds
}

func channelsDemo() {
    fmt.Println("channelsDemo-----------")
    fmt.Println(time.Now())
    channelA := make(chan int64)
    channelB := make(chan int64)

    go work(1000, channelA)
    go work(2000, channelB)

    resultA, resultB := <- channelA, <- channelB    
    fmt.Println("results ", resultA, resultB)
    fmt.Println(time.Now())
}

/*
Buffered Channels
Channels can be buffered. Provide the buffer length as the second argument to make to initialize a buffered channel:

ch := make(chan int, 100)
Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.

Modify the example to overfill the buffer and see what happens.
*/
func bufferedChannels() {
    fmt.Println("BufferedChannels----------------")
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    val1, val2 := <-ch, <-ch
    fmt.Println(val1, val2)
}

/*
Range and Close
A sender can close a channel to indicate that no more values will be sent. Receivers can test whether a channel has been closed by assigning a second parameter to the receive expression: after

v, ok := <-ch
ok is false if there are no more values to receive and the channel is closed.

The loop for i := range c receives values from the channel repeatedly until it is closed.
Note: Only the sender should close a channel, never the receiver. Sending on a closed channel will cause a panic.
Another note: Channels aren't like files; you don't usually need to close them. Closing is only necessary when the receiver must be told there are no more values coming, such as to terminate a range loop.
*/

func fibonacci(n int, c chan int) {
    time.Sleep(1000 * time.Millisecond)
	x, y := 0, 1
	for i := 0; i < n; i++ {        
		c <- x
		x, y = y, x+y
        time.Sleep(100 * time.Millisecond)
	}
	close(c) // close the channel
}

func rangeAndClose() {
    fmt.Println("rangeAndClose------------")
    c := make(chan int, 10)
	go fibonacci(cap(c), c)
    fmt.Println("range listening channel")
	for i := range c { // range receives values until channel is closed
		fmt.Println(i)
	}
    fmt.Println("closed listening channel")
}

/*
The select statement lets a goroutine wait on multiple communication operations.
A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
*/

func fibonacci2(c, quit chan int) {
	x, y := 0, 1
	for {
		select { // blocks until it is signaled
		case c <- x:    
			x, y = y, x+y
		// case quit <- 0: // invoked only, invoker doesn't send a value, we have a fallback
		// 	fmt.Println("quit")
		// 	return
        case <-quit: // only listening, not getting the value
			fmt.Println("quit")
			return
		}
	}
}

func selectStatement() {
    fmt.Println("selectStatement------------")
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c) // doesn't send a val. only triggers
		}
		//<-quit // invoke only, no send value
        quit <- 0 //. sending a value
	}()
	fibonacci2(c, quit)
}

/*
The default case in a select is run if no other case is ready.

Use a default case to try a send or receive without blocking:

select {
case i := <-c:
    // use i
default:
    // receiving from c would block
}
*/
func selectDefault() {
    fmt.Println("selectDefault------------")
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default: // runs multiple times until a case is ready
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

/*
Exercise: Equivalent Binary Trees
There can be many different binary trees with the same sequence of values stored in it. For example, here are two binary trees storing the sequence 1, 1, 2, 3, 5, 8, 13.
ie1:
          3
    1          8
  1   2     5    13
ie2:
        8
    3       13
  1   5
1   2  
ie3: (-((1 (2)) 3 (4))- 5 -((6) 7 ((8) 9))-) 10
                    10
            5       
    3               7
1        4       6       8
    2                       9
A function to check whether two binary trees store the same sequence is quite complex in most languages. We'll use Go's concurrency and channels to write a simple solution.
This example uses the tree package, which defines the type:
type Tree struct {
    Left  *Tree
    Value int
    Right *Tree
}
Exercise: Equivalent Binary Trees
1. Implement the Walk function.
2. Test the Walk function.
The function tree.New(k) constructs a randomly-structured (but always sorted) binary tree holding the values k, 2k, 3k, ..., 10k.
Create a new channel ch and kick off the walker:
go Walk(tree.New(1), ch)
Then read and print 10 values from the channel. It should be the numbers 1, 2, 3, ..., 10.
3. Implement the Same function using Walk to determine whether t1 and t2 store the same values.
4. Test the Same function.
Same(tree.New(1), tree.New(1)) should return true, and Same(tree.New(1), tree.New(2)) should return false.
The documentation for Tree can be found here. https://godoc.org/golang.org/x/tour/tree#Tree
*/

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    if t.Left == nil && t.Right == nil {
        ch <- t.Value
    } else {        
        if t.Left != nil {
            Walk(t.Left, ch)
            ch <- t.Value
        } else if t.Left == nil {
            ch <- t.Value
        }     
        if t.Right != nil {
            Walk(t.Right, ch)
        }
    }
}

func Equals[T comparable](a, b []T) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    var c1 = make(chan int)
    var c2 = make(chan int)
    var a1 = make([]int, 0)
    var a2 = make([]int, 0)
    go Walk(t1, c1)
    go Walk(t2, c2)
    for i := 0; i < 10; i++ {
        a1 = append(a1, <-c1)
        a2 = append(a2, <-c2)
    }
    return Equals(a1, a2)
}

func equivalentBinaryTreesExercise() {
    fmt.Println("equivalentBinaryTrees------------")

    fmt.Println("TestWalk")
    var t1 *tree.Tree = tree.New(1)
    fmt.Println("T1: ", t1)
    c := make(chan int)
    go Walk(t1, c)  
    for i := 0; i < 10; i++ {
        fmt.Printf("%d, ", <-c)
    }
    
    fmt.Println("")
    fmt.Println("TestSame")
    var t2 *tree.Tree = tree.New(2)
    fmt.Println("T2: ", t2)
    var t1Equals = Same(t1, t1)
    fmt.Println("T(1) and T(1) are equal?", t1Equals)
    var t1t2Equals = Same(t1, t2)
    fmt.Println("T(1) and T(2) are equal?", t1t2Equals)
}

/*
We've seen how channels are great for communication among goroutines.
But what if we don't need communication? What if we just want to make sure only one goroutine can access a variable at a time to avoid conflicts?
This concept is called mutual exclusion, and the conventional name for the data structure that provides it is mutex.
Go's standard library provides mutual exclusion with sync.Mutex and its two methods:
Lock
Unlock
We can define a block of code to be executed in mutual exclusion by surrounding it with a call to Lock and Unlock as shown on the Inc method.
We can also use defer to ensure the mutex will be unlocked as in the Value method.
*/

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v[key]
}

func mutex() {
    fmt.Println("mutex------------")
    c := SafeCounter{v: make(map[string]int)} // why is not needed to initialize the mutex? mu sync.Mutex is ready to use immediately upon allocation or just declaration(all internal stuff is initialized with zeros). explanation: https://go.dev/doc/effective_go#data
	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}

/*
In this exercise you'll use Go's concurrency features to parallelize a web crawler.
Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.
Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!
*/
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

type Crawled struct {
    mu sync.Mutex
    crawled map[string]bool
}

func (c *Crawled) Add(url string ) {
    c.mu.Lock()
    c.crawled[url] = true 
    c.mu.Unlock()
}

func (c *Crawled) Contains(url string) bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    _, exists := c.crawled[url]
    return exists
}

func (c *Crawled) Empty() bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    return len(c.crawled) == 0
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, crawled Crawled, wg *sync.WaitGroup) {
    if !crawled.Empty() {
        defer wg.Done()
    }    
	if depth <= 0 {        
		return
	}
    if crawled.Contains(url) {
        return
    }
	body, urls, err := fetcher.Fetch(url)    
    crawled.Add(url)
	if err != nil {
		fmt.Println(err)
		return
	}    
	fmt.Printf("found: %s %q %v\n", url, body, urls)
	for _, u := range urls {
        wg.Add(1)
		go Crawl(u, depth-1, fetcher, crawled, wg)
	}
	return
}

func webCrawlerExercise() {
    fmt.Println("webCrawlerExercise------------")
    // fetcher is a populated fakeFetcher.
    var fetcher = fakeFetcher{
        "https://golang.org/": &fakeResult{
            "The Go Programming Language",
            []string{
                "https://golang.org/pkg/",
                "https://golang.org/cmd/",
            },
        },
        "https://golang.org/pkg/": &fakeResult{
            "Packages",
            []string{
                "https://golang.org/",
                "https://golang.org/cmd/",
                "https://golang.org/pkg/fmt/",
                "https://golang.org/pkg/os/",
            },
        },
        "https://golang.org/pkg/fmt/": &fakeResult{
            "Package fmt",
            []string{
                "https://golang.org/",
                "https://golang.org/pkg/",
            },
        },
        "https://golang.org/pkg/os/": &fakeResult{
            "Package os",
            []string{
                "https://golang.org/",
                "https://golang.org/pkg/",
            },
        },
    }
    var wg sync.WaitGroup
    crawled := Crawled{crawled: make(map[string]bool)}    
    Crawl("https://golang.org/", 4, fetcher, crawled, &wg)
    wg.Wait()
}