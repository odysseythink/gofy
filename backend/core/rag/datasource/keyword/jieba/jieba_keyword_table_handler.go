package jieba

import (
	"regexp"
	"slices"
	"strings"
	"sync"

	gojieba "github.com/odysseythink/gofy/backend/utils/jieba"
	jiebautils "github.com/odysseythink/gofy/backend/utils/jieba"
)

var (
	jiebaInstance *jiebautils.Jieba
	once          sync.Once
)

type JiebaKeywordTableHandler struct {
	jieba *jiebautils.Jieba
}

// NewJiebaKeywordTableHandler 初始化（保证只加载一次词典）
func NewJiebaKeywordTableHandler() *JiebaKeywordTableHandler {
	once.Do(func() {
		jiebaInstance = gojieba.NewJieba()
		jiebaInstance.ClearStopWordDict()
		for _, v := range STOPWORDS {
			jiebaInstance.AddStopWord(v)
		}
	})
	return &JiebaKeywordTableHandler{jieba: jiebaInstance}
}
func (h *JiebaKeywordTableHandler) ExtractKeywords(text string, maxKeywordsPerChunk int) []string {
	tags := h.jieba.ExtractWithWeight(text, maxKeywordsPerChunk)
	keywords := make([]string, len(tags))
	for _, tw := range tags {
		keywords = append(keywords, tw.Word)
	}
	return h._expand_tokens_with_subtokens(keywords)
}

func (h *JiebaKeywordTableHandler) _expand_tokens_with_subtokens(tokens []string) []string {
	result := make([]string, len(tokens)*2)
	re := regexp.MustCompile(`\w+`)

	for _, tok := range tokens {
		result = append(result, tok)
		sub := re.FindAllString(tok, -1)
		if len(sub) > 1 {
			for _, s := range sub {
				s = strings.ToLower(s)
				if !slices.Contains(STOPWORDS, s) {
					result = append(result, s)
				}
			}
		}
	}
	return result
}

// def extract_keywords(self, text: str, max_keywords_per_chunk: Optional[int] = 10) -> set[str]:
//     """Extract keywords with JIEBA tfidf."""
//     import jieba.analyse  # type: ignore

//     keywords = jieba.analyse.extract_tags(
//         sentence=text,
//         topK=max_keywords_per_chunk,
//     )
//     # jieba.analyse.extract_tags returns list[Any] when withFlag is False by default.
//     keywords = cast(list[str], keywords)

//     return set(self._expand_tokens_with_subtokens(set(keywords)))

// def _expand_tokens_with_subtokens(self, tokens: set[str]) -> set[str]:
//     """Get subtokens from a list of tokens., filtering for stopwords."""
//     from core.rag.datasource.keyword.jieba.stopwords import STOPWORDS

//     results = set()
//     for token in tokens:
//         results.add(token)
//         sub_tokens = re.findall(r"\w+", token)
//         if len(sub_tokens) > 1:
//             results.update({w for w in sub_tokens if w not in list(STOPWORDS)})

//     return results
