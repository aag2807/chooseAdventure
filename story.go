package makeAdventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Make adventure</title>
</head>

<body>
	<div style="width:80%; max-width:800px; margin:auto; margin-top:40px; margin-bottom: 40px; background:#fffcf6; box-shadow: 0 10px 6px -6px #777">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
		<p style="color:blue;">{{.}}</p>
		{{end}}
		<ul style="border-top: 1px dotted #ccc; padding: 10px 0 0 0; --webkit-padding-start:0;list-style-type:none;">
			{{range .Options}}
			<li style="padding-top:10px"><a href="/{{.Chapter}}"> {{.Text}}</a> </li>
			{{end}}
		</ul>
	</div>
</body>

</html>
`

type HandlerOption func(h *Handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *Handler) {
		h.t = t
	}
}

func NewHandler(s Story, opts ...HandlerOption) Handler {
	h := Handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := h.pathFn(req)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(res, chapter)
		if err != nil {
			log.Println(err)
			http.Error(res, "something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(res, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func defaultPathFn(req *http.Request) string {
	path := strings.TrimSpace(req.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Handler struct {
	s      Story
	t      *template.Template
	pathFn func(req *http.Request) string
}

type Story map[string]Chapter
