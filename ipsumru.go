package ipsumru

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SentenceGenerator struct {
	sync.Mutex
	news2023 string
	news2024 string
	texts    []string   // исходные тексты (2023 + 2024)
	lines    []lineSpan // индексы строк
}

type lineSpan struct {
	src int // индекс, какой из texts
	beg int // начало строки
	end int // конец строки
}

func NewSentenceGenerator() (*SentenceGenerator, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home dir: %v", err)
	}
	fileName2023 := filepath.Join(homeDir, nameNews2023)
	fileName2024 := filepath.Join(homeDir, nameNews2024)

	if err = ensureFile(url2023, fileName2023); err != nil {
		return nil, err
	}
	if err = ensureFile(url2024, fileName2024); err != nil {
		return nil, err
	}

	news2023content, err := os.ReadFile(fileName2023)
	if err != nil {
		return nil, fmt.Errorf("could not read news2023 file: %v", err)
	}

	news2024content, err := os.ReadFile(fileName2023)
	if err != nil {
		return nil, fmt.Errorf("could not read news2024 file: %v", err)
	}

	g := &SentenceGenerator{
		texts: []string{string(news2023content), string(news2024content)},
	}
	b, err := os.ReadFile(fileName2023)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fileName2023, err)
	}
	g.news2023 = string(b)
	b, err = os.ReadFile(fileName2024)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fileName2024, err)
	}
	g.news2023 = string(b)

	g.reload()
	return g, nil
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
