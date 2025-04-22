// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Declares a global variable service to store a filter value (which service's logs to show).
var service string

func init() {
	// Add a flag to filter by service.d
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {
	// parses the command-line flags and stores them.
	flag.Parse()
	// strings.Builder is used to build strings efficiently.
	var b strings.Builder

	service := strings.ToLower(service)
	// Sets up a scanner to read lines from standard input (like piped log output).
	scanner := bufio.NewScanner(os.Stdin)
	// Starts reading lines one by one. Each line should be a JSON log.
	for scanner.Scan() {
		s := scanner.Text()
		// Tries to parse each log line as a JSON object into a map[string]any.
		m := make(map[string]any)
		// This line inserts to the map key value pair of the JSON object
		err := json.Unmarshal([]byte(s), &m)
		// If JSON parsing fails:
		// Print the raw log line only if no service filter was provided.
		// Skip to the next line otherwise.
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// If a service filter was provided, check.
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}

		// I like always having a traceid present in the logs.
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// {"time":"2023-06-01T17:21:11.13704718Z","level":"INFO","msg":"startup","service":"SALES-API","GOMAXPROCS":1}

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["ts"],
			m["level"],
			traceID,
			m["caller"],
			m["msg"],
		))

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			switch k {
			case "service", "time", "file", "level", "trace_id", "msg":
				continue
			}

			// It's nice to see the key[value] in this format
			// especially since map ordering is random.
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, Trims off the final trailing ": "
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
