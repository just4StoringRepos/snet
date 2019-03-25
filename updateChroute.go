// +build ignore

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var chnrouteTmplate = template.Must(
	template.New("").Parse(
		`// Code generated by go generate; DO NOT EDIT.
// Generated at {{ .Timestamp }}
package main

var Chnroutes = []string{
	{{- range .Routes }}	
		{{ printf "%q" . }},
	{{- end }}
}
`))

func ipStartCountToCIDR(ipStart string, ipCount int) string {
	result := math.Log2(float64(ipCount))
	return fmt.Sprintf("%s/%d", ipStart, 32-int(result))
}

func genChnroute() ([]string, error) {
	var apnicGlobalFile = "apnic.txt"
	f, err := os.Open(apnicGlobalFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	result := make([]string, 0, 8000) // 8000 is a approximate number of cn cidr numbers
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "apnic|CN|ipv4") {
			record := strings.FieldsFunc(line,
				func(c rune) bool {
					return c == '|'
				})

			ipCount, err := strconv.Atoi(record[4])
			if err != nil {
				return nil, err
			}
			result = append(result, ipStartCountToCIDR(record[3], ipCount))
		}
	}
	return result, nil
}

func main() {
	routes, err := genChnroute()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("chnroutes.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	chnrouteTmplate.Execute(f, struct {
		Timestamp time.Time
		Routes    []string
	}{
		Timestamp: time.Now(),
		Routes:    routes,
	})
}
