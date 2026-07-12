package main

import (
	"crypto/rand"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net"
	"net/http"
	"strings"
	"unicode"
)

//go:embed comments.txt
var commentsRaw string

var comments = strings.Split(strings.TrimSpace(commentsRaw), "\n")

func main() {
	http.HandleFunc("GET /", indexHandler)

	log.Println("Listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	name := capitalize(nameFromHost(r.Host))
	comment := fmt.Sprintf(randomComment(), name)
	title := name + " liegt auf dem Gleis"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="de">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<meta property="og:type" content="website">
<meta property="og:title" content="%s">
<meta property="og:description" content="%s">
<meta property="og:url" content="%s">
<meta name="twitter:card" content="summary">
<meta name="twitter:title" content="%s">
<meta name="twitter:description" content="%s">
<style>
body {
  margin: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  flex-direction: column;
  font-family: sans-serif;
  text-align: center;
}
pre { font-size: 1.25rem; }
h1 { font-size: clamp(3rem, 12vw, 8rem); margin: 1rem 0; }
p { font-size: 1.5rem; }
</style>
</head>
<body>
<pre>     o
    /|\
    / \
|=|=|=|=|=|=|=|=|=|</pre>
<h1>%s</h1>
<p>%s</p>
</body>
</html>`, escape(title), escape(title), escape(comment), escape(pageURL(r)), escape(title), escape(comment), escape(name), escape(comment))
}

func pageURL(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return scheme + "://" + r.Host + r.RequestURI
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

func nameFromHost(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if host == "" || net.ParseIP(host) != nil {
		return "Jemand"
	}

	parts := strings.Split(host, ".")
	if len(parts) < 3 || !validLabel(parts[0]) {
		return "Jemand"
	}

	return parts[0]
}

func validLabel(label string) bool {
	if label == "" || len(label) > 63 || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}

	for _, r := range label {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			continue
		}

		return false
	}

	return true
}

func randomComment() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(comments))))
	if err != nil {
		return comments[0]
	}

	return comments[n.Int64()]
}

func escape(s string) string {
	return template.HTMLEscapeString(s)
}
