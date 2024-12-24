package main

import (
    "html/template"
    "net/http"
    "flag"
    "log"
    "fmt"
)

type PageData struct {
    Latest string
    History []string
    BaseUrl string
}

var pageData PageData

func handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        handleIndex(w, r)
        return
    }
    http.NotFound(w, r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Clippy</title>
</head>
<body>
    <form method="POST" action="{{ .BaseUrl }}/update">
        <input type="text" name="clip" />
        <button>save</button>
    </form>
    <p><b>Latest: </b>{{ .Latest }}</p>
    <ul>
        {{ range .History }}
            <li>{{ . }}</li>
        {{ end }}
    </ul>
</body>
</html>
`
    tmpl, err := template.New("index").Parse(html)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, pageData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    return
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
    clip := r.FormValue("clip")
    if len(clip) != 0 {
        if len(pageData.Latest) != 0 {
            pageData.History = append([]string{ pageData.Latest }, pageData.History...)
        }
        pageData.Latest = clip
    }
    http.Redirect(w, r, pageData.BaseUrl, http.StatusFound)
}

func handleGetLatest(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, pageData.Latest)
}

func main() {
    port := flag.String("port", "8090", "port")
    baseUrl := flag.String("base-url", "", "base url")
    flag.Parse()

    if *baseUrl == "/" {
        *baseUrl = ""
    }

    pageData.BaseUrl = *baseUrl

    if pageData.BaseUrl == "" {
        http.HandleFunc("GET /", handleRoot)
    } else {
        http.HandleFunc("GET " + pageData.BaseUrl, handleIndex)
    }
    http.HandleFunc("POST " + pageData.BaseUrl + "/update", handleUpdate)
    http.HandleFunc("GET " + pageData.BaseUrl + "/latest", handleGetLatest)

    log.Println("starting server on port " + *port)
    http.ListenAndServe(":" + *port, nil)
}
