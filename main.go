package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Measurement struct {
	Timestamp string
	CPU       string
	Mem       string
}

func elapsedTimeMilliseconds(start, stop time.Time) int {
	return int(stop.Sub(start) / time.Millisecond)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must supply absolute path to CSV file")
	}
	if len(os.Args) < 3 {
		log.Fatal("Must supply start timestamp")
	}
	if len(os.Args) < 4 {
		log.Fatal("Must supply stop timestamp")
	}

	path := os.Args[1]
	start := os.Args[2]
	stop := os.Args[3]

	startTime, err := time.Parse(time.RFC3339Nano, start)
	if err != nil {
		log.Println(err)
	}

	//fmt.Println("Full path:", path)

	// Open file
	f, err := os.Open(path) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	// Read entire file
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f) // Error handling elided for brevity.
	f.Close()
	s := strings.Replace(string(buf.Bytes()), ";", ",", -1)

	r := csv.NewReader(strings.NewReader(s))

	measurements := make(map[int]*Measurement)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == "Timestamp" || record[0] < start || record[0] > stop {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			log.Fatal(err)
		}

		key := elapsedTimeMilliseconds(startTime, timestamp)

		if _, ok := measurements[key]; ok {
			// Second instance of same timestamp, offset by 500 milliseconds
			key += 500
		}

		measurements[key] = &Measurement{
			Timestamp: record[0],
			CPU:       record[1],
			Mem:       record[2],
		}

		//fmt.Println(record)
	}

	//result, err := json.Marshal(measurements)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Println(string(result))

	// To store the keys in slice in sorted order
	var keys []int
	for k := range measurements {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println("Time (ms),%CPU,%MEM,OriginalTimestamp")

	for _, k := range keys {
		m := measurements[k]
		fmt.Printf("%d,%s,%s,%s\n", k, m.CPU, m.Mem, m.Timestamp)
	}
}
