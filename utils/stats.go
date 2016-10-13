// WIP //

package stats

import (
	"fmt"
	"time"
)
/*
import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const Version string = "0.0.4"

type Round struct {
	Client *http.Client
	Reqs   []*http.Request
	Resp   []*http.Response
	Input  string
	Sent   []time.Time
	Recv   []time.Time
	Trip   time.Duration
}

type Config struct {
	Workers int
	Reps    int
	Quiet   bool
	Verbose bool
	Infile  *os.File
}

type Payload struct {
	Method  string
	Headers map[string]string
	Body    string
	Url     string
}

func worker(jobs <-chan *Round, resp chan<- *Round) {
	for round := range jobs {
		for _, req := range round.Reqs {
			round.Sent = append(round.Sent, time.Now())

			res, err := round.Client.Do(req)
			check(err)

			round.Recv = append(round.Recv, time.Now())
			round.Resp = append(round.Resp, res)
		}
		resp <- round
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// new main to have a chance to clean up upon failure
func main() {
	_main()
}

// init
func _main() {
	start := time.Now()

	config := &Config{
		Workers: 1,
		Reps:    1,
		Quiet:   false,
		Verbose: true,
	}

	payloads := make([]*Payload, 1)
	payloads[0] = &Payload{
		Method:  "GET",
		Headers: make(map[string]string),
	}

	args := os.Args[1:]
	for x := 0; x < len(args); x++ {
		switch args[x] {
		case "-v", "--version":
			printVer()
			os.Exit(0)
		case "--verbose=false":
			config.Verbose = false
		case "-q", "--quiet":
			config.Quiet = true
			config.Verbose = false
		case "-w", "--workers":
			i, err := strconv.Atoi(args[x+1])
			check(err)
			config.Workers = i
			x++
		case "-r", "--reps":
			i, err := strconv.Atoi(args[x+1])
			check(err)
			config.Reps = i
			x++
		case "-f", "--file":
			f, err := os.Open(args[x+1])
			check(err)
			config.Infile = f
			defer config.Infile.Close()
			payloads = parseConfig(config.Infile, payloads)
			x++
		case "-X":
			payloads[0].Method = args[x+1]
			x++
		case "-d":
			payloads[0].Body = args[x+1]
			x++
		case "-H":
			h := strings.Split(args[x+1], ":")
			payloads[0].Headers[strings.TrimSpace(h[0])] = strings.TrimSpace(h[1])
			x++
		case "-h", "--help":
			printHelp()
		default:
			if strings.HasPrefix(args[x], "http") && config.Infile == nil {
				payloads[0].Url = args[x]
			} else {
				printHelp()
			}
		}
	}

	// channels
	jobs := make(chan *Round, config.Workers)
	resp := make(chan *Round, config.Workers)

	// spawn worker pool
	spawnPool(jobs, resp, config.Workers)

	client := createHttpClient()
	rounds := make([]*Round, config.Workers)
	for x := 0; x < len(rounds); x++ {
		rounds[x] = &Round{
			Client: client,
			Reqs:   make([]*http.Request, 0),
			Resp:   make([]*http.Response, 0),
			Sent:   make([]time.Time, 0),
			Recv:   make([]time.Time, 0),
			Trip:   0,
		}
		// create each payload for each round
		for _, payload := range payloads {
			for y := 0; y < config.Reps; y++ {
				rounds[x].Reqs = append(rounds[x].Reqs, craftRequest(payload))
			}
		}
	}
	spawnJobs(jobs, rounds)
	rounds = printResponses(resp, config.Workers, config.Verbose)

	// clean up and print report
	close(jobs)
	report(config.Quiet, rounds, start)
}

func printVer() {
	fmt.Printf("Autocannon v%v\n\n", Version)
}

func printHelp() {
	printVer()
	fmt.Println("Usage: ac [options]")
	fmt.Println("       ac [options] [url]")
	fmt.Println("       ac -w 2 -r 2 https://example.com")
	fmt.Println("Options:")
	fmt.Println("-v, --version      VERSION, prints the autocannon version and exits.")
	fmt.Println("--verbose          VERBOSE, prints all requests and responses. This is the default behavior.")
	fmt.Println("-q, --quiet        QUIET, supresses all output and logs the final report.")
	fmt.Println("-w, --workers      WORKERS, sets concurrency or amount of simulated users.")
	fmt.Println("-r, --reps         REPS, sets the amount of requests to be sent per simulated user.")
	fmt.Println("-h, --help         HELP, prints this help message.")
	fmt.Println("-f, --file         FILE, specify a file to parse for options.")
	fmt.Println("-X                 METHOD, request method to use (GET, POST, PUT, DELETE).")
	fmt.Println("-d                 BODY, data to be sent in the body of the request.")
	fmt.Println("-H                 HEADER, sets a key-value pair to be sent in the header of the request.")
	os.Exit(0)
}

func replaceRand(s string) string {
	data := make([]byte, 20)
	var random string

	if _, err := rand.Read(data); err == nil {
		random = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	s = strings.Replace(s, "<rand>", random[:16], -1)

	return s
}

func parseConfig(infile *os.File, payloads []*Payload) []*Payload {
	// read line of file and do shit with it
	x := 0
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		opts := strings.Split(scanner.Text(), " ")

		p := &Payload{
			Method:  "GET",
			Headers: make(map[string]string),
			Body:    "",
			Url:     opts[len(opts)-1],
		}

		for y := 0; y < len(opts); y++ {
			switch opts[y] {
			case "-X":
				p.Method = opts[y+1]
				y++
			case "-d":
				p.Body = opts[y+1]
				y++
			case "-H":
				h := strings.Split(opts[y+1], ":")
				p.Headers[strings.TrimSpace(h[0])] = strings.TrimSpace(h[1])
				y++
			default:
				if strings.HasPrefix(opts[y], "http") {
					p.Url = opts[y]
				}
			}
		}

		if x > cap(payloads)-1 {
			payloads = append(payloads, p)
		} else {
			payloads[x] = p
		}
		x++
	}
	err := scanner.Err()
	check(err)
	return payloads
}

func createHttpClient() *http.Client {
	tls := &tls.Config{
		ClientAuth:         tls.VerifyClientCertIfGiven,
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{
		Proxy:                 nil,
		Dial:                  nil,
		DialTLS:               nil,
		TLSClientConfig:       tls,
		TLSHandshakeTimeout:   0,
		DisableKeepAlives:     true,
		DisableCompression:    false,
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 0,
		TLSNextProto:          nil,
	}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	return client
}

func craftRequest(payload *Payload) *http.Request {
	url := replaceRand(payload.Url)
	body := replaceRand(payload.Body)
	req, _ := http.NewRequest(payload.Method, url, strings.NewReader(body))
	for k, v := range payload.Headers {
		req.Header.Add(k, v)
	}
	return req
}

func spawnPool(jobs chan *Round, resp chan *Round, workers int) {
	for w := 0; w < workers; w++ {
		go worker(jobs, resp)
	}
}

func spawnJobs(jobs chan *Round, rounds []*Round) {
	for _, round := range rounds {
		jobs <- round
	}
}

func printResponses(resp chan *Round, workers int, verbose bool) []*Round {
	var round *Round
	rounds := make([]*Round, 0)

	for x := 0; x < workers; x++ {
		round = <-resp
		if verbose {
			for y := 0; y < len(round.Resp); y++ {
				round.Trip += round.Recv[y].Sub(round.Sent[y])
				fmt.Printf("%8v %v %8.2f secs: %8v bytes ==> %v %v\n",
					round.Resp[y].Proto,
					round.Resp[y].StatusCode,
					round.Recv[y].Sub(round.Sent[y]).Seconds(),
					round.Reqs[y].ContentLength,
					round.Reqs[y].Method,
					round.Reqs[y].URL)
			}
		}
		rounds = append(rounds, round)
	}
	return rounds
}
*/

// Report generates statistics and formats the output when given a Round struct and the start time of the program.
func Report(quiet bool, rounds []*Round, start time.Time) {
	if !quiet {
		var count, transferred, concurrency, workingTime, longest, shortest, availability float64
		count = 0
		longest = 0
		shortest = 999
		workers := int(len(rounds))
		transactions := workers * len(rounds[0].Reqs)
		for x := 0; x < workers; x++ {
			for y := 0; y < len(rounds[x].Reqs); y++ {
				if rounds[x].Resp[y].StatusCode == 200 {
					count++
				}
				if rounds[x].Recv[y].Sub(rounds[x].Sent[y]).Seconds() > longest {
					longest = rounds[x].Recv[y].Sub(rounds[x].Sent[y]).Seconds()
				}
				if rounds[x].Recv[y].Sub(rounds[x].Sent[y]).Seconds() < shortest {
					shortest = rounds[x].Recv[y].Sub(rounds[x].Sent[y]).Seconds()
				}
				transferred += float64(rounds[x].Reqs[y].ContentLength)
			}
			workingTime += rounds[x].Trip.Seconds()
		}
		availability = count / float64(transactions) * 100
		totalTime := time.Now().Sub(start).Seconds()
		transferred = transferred / 1000000
		avgTime := workingTime / float64(transactions)
		rate := float64(transactions) / totalTime
		throughput := transferred / totalTime
		concurrency = workingTime / totalTime

		fmt.Println()
		fmt.Printf("Transactions: %14v\n", transactions)
		fmt.Printf("Availability: %14.0f %%\n", availability)
		fmt.Printf("Elapsed time: %14.2f secs\n", totalTime)
		fmt.Printf("Data transferred: %10.2f MB\n", transferred)
		fmt.Printf("Response time: %13.2f secs\n", avgTime)
		fmt.Printf("Transaction rate: %10.2f trans/sec\n", rate)
		fmt.Printf("Throughput: %16.2f MB/sec\n", throughput)
		fmt.Printf("Concurrency: %15.2f\n", concurrency)
		fmt.Printf("Longest Transaction: %7.2f\n", longest)
		fmt.Printf("Shortest Transaction: %6.2f\n", shortest)
	}
}
