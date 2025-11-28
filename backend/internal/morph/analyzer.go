package morph

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/y-plusknasy/jifree-chrome-extension/backend/internal/api"
)

type Analyzer struct {
	t *tokenizer.Tokenizer
}

func New() (*Analyzer, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tokenizer: %w", err)
	}
	return &Analyzer{t: t}, nil
}

func (a *Analyzer) Analyze(text string) ([]api.Token, error) {
	tokens := a.t.Tokenize(text)
	var result []api.Token

	for _, t := range tokens {
		if t.Class == tokenizer.DUMMY {
			continue
		}

		features := t.Features()
		reading := t.Surface // デフォルトは表層形（読みがない場合や記号など）
		pos := ""

		if len(features) > 0 {
			pos = features[0]
		}
		// IPADICでは index 7 が読み仮名（カタカナ）
		if len(features) > 7 {
			reading = features[7]
		}

		result = append(result, api.Token{
			BaseForm:     t.Surface, // フロントエンドでの置換用に表層形を返す
			Reading:      reading,
			PartOfSpeech: pos,
		})
	}

	return result, nil
}

// ExtractText はHTML文字列からテキストのみを抽出します。
// scriptタグやstyleタグの中身は除外します。
func (a *Analyzer) ExtractText(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse html: %w", err)
	}

	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return sb.String(), nil
}
