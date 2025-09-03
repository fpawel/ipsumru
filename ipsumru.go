package ipsumru

import (
	_ "embed"
	"math/rand"
	"strings"
	"sync"
)

// https://wortschatz.uni-leipzig.de/en/download/Russian?utm_source=chatgpt.com

//go:embed data/rus_news_2023_1M-sentences.txt
var news2023 string

//go:embed data/rus_news_2024_1M-sentences.txt
var news2024 string

type SentenceGenerator struct {
	sync.Mutex
	texts []string   // исходные тексты (2023 + 2024)
	lines []lineSpan // индексы строк
}

type lineSpan struct {
	src int // индекс, какой из texts
	beg int // начало строки
	end int // конец строки
}

func NewSentenceGenerator() *SentenceGenerator {
	g := &SentenceGenerator{
		texts: []string{news2023, news2024},
	}
	g.reload()
	return g
}

// reload строит индексы строк заново
func (g *SentenceGenerator) reload() {
	var spans []lineSpan
	for idx, txt := range g.texts {
		offset := 0
		for {
			nl := strings.IndexByte(txt[offset:], '\n')
			if nl < 0 {
				if offset < len(txt) {
					spans = append(spans, lineSpan{src: idx, beg: offset, end: len(txt)})
				}
				break
			}
			end := offset + nl
			if end > offset {
				spans = append(spans, lineSpan{src: idx, beg: offset, end: end})
			}
			offset = end + 1
		}
	}
	g.lines = spans
}

func (g *SentenceGenerator) NextSentences(sentencesCount int) string {
	var sb strings.Builder
	for i := 0; i < sentencesCount; i++ {
		if i != 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(g.NextSentence())
	}
	return sb.String()
}

// NextSentence возвращает случайное предложение без повторов,
// когда кончатся строки — перезагружает индексы
func (g *SentenceGenerator) NextSentence() string {
	g.Lock()
	defer g.Unlock()

	if len(g.lines) == 0 {
		g.reload()
	}

	// выбираем случайный lineSpan
	i := rand.Intn(len(g.lines))
	span := g.lines[i]

	// вырезаем строку из исходного текста
	line := g.texts[span.src][span.beg:span.end]

	// убираем id до табуляции
	if tab := strings.IndexByte(line, '\t'); tab >= 0 && tab+1 < len(line) {
		line = line[tab+1:]
	}

	// удаляем элемент (swap+pop)
	last := len(g.lines) - 1
	g.lines[i] = g.lines[last]
	g.lines = g.lines[:last]

	return line
}
