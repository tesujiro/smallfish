package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	cases := []struct {
		url    string
		err    error
		status int
	}{
		//{url: "", err: nil, status: 200},
		{url: "/", err: nil, status: 200},
		{url: "/consumer/@123,456", err: nil, status: 200},
	}

	for _, c := range cases {
		r, err := http.NewRequest("GET", c.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		Router().ServeHTTP(w, r)

		if status := w.Code; status != c.status {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, c.status)
		}
		fmt.Printf("w.Body=%v\n", w.Body.String())

	}
}

/*
func main() {
	err, conf := config()
	if err != nil {
		log.Fatal(err)
		return
	}
	var cancel, timeoutCancel context.CancelFunc
	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)
	if conf.timeout > 0 {
		ctx, timeoutCancel = context.WithTimeout(ctx, time.Second*time.Duration(conf.timeout))
		defer timeoutCancel()
	}
	wg := &sync.WaitGroup{}

	// tester goroutine start
	got := make(chan struct{}, conf.threads)
	for i := 0; i < conf.threads; i++ {
		wg.Add(1)
		go func() {
			conf.tester.run(ctx, got)
			wg.Done()
		}()
	}

	// wait goroutine start
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	// tick
	var tick <-chan time.Time
	if conf.tick > 0 {
		tick = time.NewTicker(time.Second * time.Duration(conf.tick)).C
	}

	// signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// select all the events
	var counter int
L:
	for {
		select {
		case <-got:
			counter += 1
		case <-tick:
			log.Printf("tick. current %d Requests returned OK.\n", counter)
		case <-done:
			//log.Printf("Waitgroup done:\n")
			break L
		case <-ctx.Done():
			//log.Printf("Context done:\n")
			break L
		case s := <-sig:
			log.Printf("Got signal:%d\n", s)
			cancel()
		}
	}
	log.Printf("Finished. %d Requests returned OK.\n", counter)

}

type conf struct {
	threads int
	tick    int
	timeout int
	tester  *tester
}

type tester struct {
	client    *http.Client
	url       string
	loop      int
	min       int
	max       int
	keepalive bool
	debug     bool
}

func config() (error, *conf) {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 50
	http.DefaultClient.Timeout = 0
	//client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	t := tester{
		client: http.DefaultClient,
	}
	c := conf{
		tester: &t,
	}
	flag.IntVar(&c.threads, "thread", 10, "threads")
	flag.IntVar(&c.tick, "tick", 0, "tick in seconds")
	flag.IntVar(&c.timeout, "timeout", 0, "timeout in seconds")
	flag.StringVar(&t.url, "url", "http://127.0.0.1:80", "request url")
	flag.IntVar(&t.loop, "loop", 0, "loop limit")
	flag.IntVar(&t.min, "min", 0, "min sleep timer in msec")
	flag.IntVar(&t.max, "max", 100, "max sleep timer in msec")
	flag.BoolVar(&t.keepalive, "keepalive", false, "keep alive Tcp connections")
	flag.BoolVar(&t.debug, "debug", false, "debug")
	flag.Parse()
	if t.min > t.max {
		err := fmt.Errorf("Error: min > max")
		return err, &c
	} else if c.threads < 0 || c.tick < 0 || c.timeout < 0 || t.min < 0 || t.max < 0 || t.loop < 0 {
		err := fmt.Errorf("Error: negative number")
		return err, &c
	}
	return nil, &c
}

func (t *tester) run(ctx context.Context, got chan<- struct{}) {
	for i := 0; t.loop <= 0 || i < t.loop; i++ {
		if err := t.get(ctx); err == nil {
			got <- struct{}{}
		} else {
			log.Println(err)
		}
	}
}

func (t *tester) get(ctx context.Context) error {
	// set Query parameter "timer=n"
	values := url.Values{}
	if t.min == t.max {
		values.Add("timer", strconv.Itoa(t.min))
	} else {
		values.Add("timer", strconv.Itoa(t.min+rand.Intn(t.max-t.min)))
	}

	req, err := http.NewRequest("GET", t.url+"/", nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = values.Encode()

	req.WithContext(ctx)
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if t.keepalive || t.debug {
		t.dump(resp)
	}
	return nil
}

func (t *tester) dump(resp *http.Response) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if t.debug {
		log.Print(string(b))
	}
}
*/
