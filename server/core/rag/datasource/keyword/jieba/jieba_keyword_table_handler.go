package jieba

import (
	"regexp"
	"slices"
	"strings"
	"sync"

	gojieba "mlib.com/gofy/server/utils/jieba"
	jiebautils "mlib.com/gofy/server/utils/jieba"
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
func (h *JiebaKeywordTableHandler) ExtractKeywords(text string, maxKeywordsPerChunk int) map[string]struct{} {
	tags := h.jieba.ExtractWithWeight(text, maxKeywordsPerChunk)
	keywords := make(map[string]struct{}, len(tags))
	for _, tw := range tags {
		keywords[tw.Word] = struct{}{}
	}
	return h._expand_tokens_with_subtokens(keywords)
}

func (h *JiebaKeywordTableHandler) _expand_tokens_with_subtokens(tokens map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{}, len(tokens)*2)
	re := regexp.MustCompile(`\w+`)

	for tok := range tokens {
		result[tok] = struct{}{}
		sub := re.FindAllString(tok, -1)
		if len(sub) > 1 {
			for _, s := range sub {
				s = strings.ToLower(s)
				if !slices.Contains(STOPWORDS, s) {
					result[s] = struct{}{}
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
