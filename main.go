package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net"
	"net/http"
	"strings"
)

var comments = []string{"Bleib ruhig liegen, %s. Deine Commits waren eh alle gottlos scheiße.", "Endlich tust du mal was Produktives für die Gesellschaft, %s. Als Bahnschwelle reichst du gerade so.", "Mach dir keine Hoffnungen, %s. Nicht mal der Zug hat Bock auf dich.", "Leg dich ruhig hin, %s. Dein ganzes Leben hat eh mehr Bugs als dein Code.", "Mach Platz für jemanden, der tatsächlich Mehrwert liefert, %s.", "Selbst die DB-Infrastruktur ist stabiler als deine mentale Gesundheit, %s.", "Mach dir keine Sorgen um deine Projekte, %s. Ohne deine Bugs laufen die eh viel besser.", "Deine Existenz hat mehr Memory Leaks als ein ungetesteter C-Parser, %s.", "Endlich machst du mal Platz im RAM, %s.", "Keine Panik, %s. Die DNS-Propagation deines Ablebens dauert eh nicht lang.", "Sogar ein PHP-Script in Production hat mehr Daseinsberechtigung als du, %s.", "Schreib noch schnell einen Abschieds-Commit, %s. Ach ne, deine Commit-Messages rafft eh keiner.", "Gleich überrollt dich der Zug, %s. Das ist wahrscheinlich der erste echte Impact, den du je hast.", "Dein GitHub-Green-Graph ist eh so tot wie deine Zukunft, %s. Bleib einfach liegen.", "Keine Sorge, %s. Ein Rollback auf deine Geburt würde eh fehlschlagen.", "Selbst eine NullPointerException hat mehr Struktur als dein Alltag, %s.", "Deine Freundin hat dich bestimmt schon auf 'Deprecated' gesetzt, %s.", "Du wurdest erfolgreich aus dem Genpool geforkt, %s.", "Mach dir nichts draus, %s. Dein Leben war eh nur ein ungefilterter Stacktrace voller Errors.", "Leg dich hin, %s. Deine CPU-Auslastung sinkt gleich dauerhaft auf 0%%.", "Nicht mal ein Garbage Collector würde dich einsammeln, %s.", "Deine gesamte Karriere war eh nur ein einziger Merge-Konflikt, %s.", "Bleib liegen, %s. Jedes Legacy-System ist einfacher zu warten als deine Ausreden.", "Selbst Windows Vista lief runder als deine Lebensplanung, %s.", "Gleich kommt der Zug, %s. Endlich mal ein Event-Handler, den du nicht blockieren kannst.", "Mach einfach die Augen zu, %s. Du bist ab jetzt offline.", "Deine Eltern wollten dich bestimmt schon damals nach dem ersten Commit deleten, %s.", "Sogar ein ungelabeltes Jira-Ticket hat mehr Priorität als du, %s.", "Der Zug hat wenigstens einen Fahrplan – im Gegensatz zu deiner Zukunft, %s.", "Gleich wirst du hart geshuttet, %s. sudo shutdown -h now."}

func main() {
	http.HandleFunc("GET /", indexHandler)

	log.Println("Listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	name := nameFromHost(r.Host)
	comment := fmt.Sprintf(randomComment(), name)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="de">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s liegt auf dem Gleis</title>
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
</html>`, escape(name), escape(name), escape(comment))
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
