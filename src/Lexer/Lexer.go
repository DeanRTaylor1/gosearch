package lexer

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/net/html"

	"github.com/tebeka/snowball"
)

type Lexer struct {
	content []rune
}

type stat struct {
	token string
	freq  int
}

func NewLexer(content string) *Lexer {
	return &Lexer{[]rune(content)}
}

func (l *Lexer) TrimLeft() {
	for len(l.content) > 0 && unicode.IsSpace(rune(l.content[0])) {
		l.content = l.content[1:]
	}
}

func (l *Lexer) Chop(n int) (token []rune) {
	token = l.content[:n]
	l.content = l.content[n:]
	return token
}

func (l *Lexer) ChopWhile(f func(rune) bool) (token []rune) {
	n := 0
	for n < len(l.content) && f(l.content[n]) {
		n += 1
	}
	return l.Chop(n)
}

func (l *Lexer) NextToken() []rune {

	l.TrimLeft()

	if len(l.content) == 0 {
		//fmt.Println("end of content")
		return nil
	}
	if unicode.IsNumber(l.content[0]) {
		return l.ChopWhile(unicode.IsNumber)
	}
	if unicode.IsLetter(l.content[0]) {
		stemmer, err := snowball.New("english")
		if err != nil {
			fmt.Println(err)
		}
		defer stemmer.Close()

		term := l.ChopWhile(func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsNumber(r)
		})

		return []rune(stemmer.Stem(strings.ToLower(string(term))))

	}
	return l.Chop(1)
}

func (l *Lexer) Next() (string, error) {

	token := l.NextToken()
	if token == nil {
		return "EOF", errors.New("no more tokens")
	}
	return (string(token)), nil
}

func ParseLinks(htmlContent string) []string {
	links := []string{}
	nodes, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(nodes)
	return links
}

/*func indexDocument(content string) map[string]int {*/
/*return*/
/*}*/

func ReadEntireXMLFile(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var content string

	d := xml.NewDecoder(f)
	for {
		t, err := d.Token()
		if err != nil {
			break
		}

		switch se := t.(type) {
		case xml.CharData:
			content += string(se)
		}
	}
	return content
}

func ReadEntireHTMLFile(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var content string

	d := html.NewTokenizer(f)
	for {
		tt := d.Next()
		switch tt {
		case html.ErrorToken:
			return content
		case html.TextToken:
			content += string(d.Text())
		}
	}
}
func MapToSortedSlice(m map[string]int) (stats []stat) {
	for k, v := range m {
		stats = append(stats, struct {
			token string
			freq  int
		}{k, v})
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].freq > stats[j].freq })

	return stats
}

func LogStats(filePath string, stats []stat, topN int) {
	fmt.Println(filePath)
	if len(stats) < topN {
		for _, v := range stats {
			fmt.Println(v.token, " => ", v.freq)
		}
	} else {
		for _, v := range stats[:topN] {
			fmt.Println(v.token, " => ", v.freq)
		}
	}
}
