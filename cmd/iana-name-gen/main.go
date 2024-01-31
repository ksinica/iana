package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	urlFlag     = flag.String("url", "", "")
	outputFlag  = flag.String("output", "", "")
	prefixFlag  = flag.String("prefix", "", "")
	packageFlag = flag.String("package", "", "Go package name")
)

func init() {
	flag.Parse()
}

const goTmpl = `// Code generated from {{ .URL }}; DO NOT EDIT.

package {{.Package}}
{{$context := .}}
const (
{{- range .Values}}	
	{{ Title $context.Prefix }}{{ ToName . }} = "{{if $context.Prefix}}{{ $context.Prefix }}/{{end}}{{.}}"
{{- end}}
)
`

var Title = cases.Title(language.AmericanEnglish, cases.NoLower).String

var tmpl = template.Must(
	template.New("go").
		Funcs(template.FuncMap{
			"ToName": ToName,
			"Title":  Title,
		}).
		Parse(goTmpl),
)

var blacklist = map[string]struct{}{
	"*": {},
}

func ToName(name string) string {
	s := strings.FieldsFunc(name, func(r rune) bool {
		return r == '+' || r == '-' || r == '.'
	})

	for i := range s {
		s[i] = Title(s[i])
	}

	return strings.Join(s, "")
}

func parseCSV(r io.Reader) (ret []string, err error) {
	s := bufio.NewScanner(r)
	for i := 0; s.Scan(); i++ {
		if i == 0 {
			continue
		}

		if s := strings.Split(s.Text(), ","); len(s) > 0 {
			if s = strings.Fields(s[0]); len(s) > 0 {
				if _, ok := blacklist[s[0]]; !ok {
					ret = append(ret, s[0])
				}
			}
		}
	}
	err = s.Err()

	slices.Sort(ret)
	slices.Compact(ret)
	return
}

func drainAndClose(resp *http.Response) {
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func fetch(ctx context.Context, url *url.URL) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		drainAndClose(resp)
		return nil, fmt.Errorf(
			"%d: %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
	}
	return resp, nil
}

func run() error {
	url, err := url.Parse(*urlFlag)
	if err != nil {
		return err
	}

	fmt.Println("Fetching", url.String())

	ctx, cf := context.WithTimeout(context.Background(), time.Second*15)
	defer cf()

	resp, err := fetch(ctx, url)
	if err != nil {
		return err
	}
	defer drainAndClose(resp)

	vals, err := parseCSV(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(*outputFlag)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, struct {
		URL     string
		Package string
		Prefix  string
		Values  []string
	}{
		URL:     *urlFlag,
		Package: *packageFlag,
		Prefix:  *prefixFlag,
		Values:  vals,
	})
}

func printUsage() {
	const usage = `Names generator from IANA registry's CSV files.
 
Usage:
    %s --url=<url> --output=<file> --package=<package> [--prefix=<prefix>]

Options:
    --url=<url>          URL of the CSV document.
    --output=<file>      Output file.
    --package=<package>  Go package name.
    --prefix=<prefix>    Prefix.
`
	fmt.Fprintf(os.Stdout, usage, os.Args[0])
}

func printError(args ...any) {
	fmt.Fprintln(os.Stderr, []any{"Error:", args}...)
}

func main() {
	if len(*urlFlag) == 0 || len(*outputFlag) == 0 || len(*packageFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	if err := run(); err != nil {
		printError(err)
		os.Exit(1)
	}
	os.Exit(0)
}
