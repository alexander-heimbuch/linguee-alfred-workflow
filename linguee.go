package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	alfred "github.com/pascalw/go-alfred"
)

type Request struct {
	URL   string
	Lang  string
	Query string
}

type Translation struct {
	Meaning string
	Phrase  []string
}

func main() {
	base := "http://www.linguee.de/%s/search?qe=%s&source=auto"

	args := params(base, os.Args)

	selection := filterDocument(transformToDocument(request(args)))

	response := alfred.NewResponse()

	for index, translation := range parseTranslations(selection) {
		response.AddItem(&alfred.AlfredResponseItem{
			Valid:    true,
			Uid:      strconv.Itoa(index),
			Title:    translation.Meaning,
			Subtitle: strings.Join(translation.Phrase, ", "),
			Arg:      fmt.Sprintf("http://www.linguee.de/%s/search?source=auto&query=%s", args.Lang, url.QueryEscape(translation.Meaning)),
		})
	}

	response.Print()
}

func params(base string, args []string) Request {
	if len(args) == 1 {
		log.Fatal("Missing Arguments")
	}

	query := args[1]
	lang := "deutsch-englisch"

	if len(args) > 2 {
		lang = args[2]
	}

	return Request{base, lang, url.QueryEscape(query)}
}

func request(req Request) *iconv.Reader {
	// Http request
	res, err := http.Get(fmt.Sprintf(req.URL, req.Lang, req.Query))
	if err != nil {
		log.Fatal(err)
	}

	// charset conversion
	body, err := iconv.NewReader(res.Body, "iso-8859-15", "utf-8")
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func transformToDocument(reader *iconv.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func filterDocument(doc *goquery.Document) *goquery.Selection {
	doc.Find(".wordtype, .sep, .grammar_info").ReplaceWithHtml("")
	return doc.Find(".autocompletion_item")
}

func parseTranslations(elements *goquery.Selection) (results []Translation) {
	elements.Each(func(index int, element *goquery.Selection) {
		results = append(results, Translation{parseMeaning(element), parsePhrase(element)})
	})

	return
}

func parseMeaning(selection *goquery.Selection) string {
	return strings.TrimSpace(selection.Find(".main_item").Text())
}

func parsePhrase(selection *goquery.Selection) (result []string) {
	selection.Find(".translation_item").Each(func(index int, meaning *goquery.Selection) {
		result = append(result, strings.TrimSpace(meaning.Text()))
	})

	return
}
