package weaviate

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
)

func hashUUID(s string) string {
	sum := sha1.Sum([]byte("gofy-weaviate:" + s))
	h := hex.EncodeToString(sum[:16])
	return fmt.Sprintf("%s-%s-%s-%s-%s", h[0:8], h[8:12], h[12:16], h[16:20], h[20:32])
}

func sanitizeClass(s string) string {
	var b strings.Builder
	for i, r := range s {
		switch {
		case i == 0 && r >= 'a' && r <= 'z':
			b.WriteRune(r - 'a' + 'A')
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" || (out[0] >= '0' && out[0] <= '9') {
		out = "C" + out
	}
	return out
}

func floatSlice(v []float32) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, x := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", x)
	}
	b.WriteByte(']')
	return b.String()
}

func toJSONString(m map[string]any) string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func jsonUnmarshal(s string, out any) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), out)
}

func buildGraphQLWhere(opts ragentities.SearchOptions) string {
	var operands []string
	for k, val := range opts.Filter {
		operands = append(operands, fmt.Sprintf(`{path:["%s"],operator:Equal,valueText:%q}`, k, fmt.Sprint(val)))
	}
	if len(opts.DocumentIDs) > 0 {
		quoted := make([]string, len(opts.DocumentIDs))
		for i, id := range opts.DocumentIDs {
			quoted[i] = fmt.Sprintf("%q", id)
		}
		operands = append(operands, fmt.Sprintf(`{path:["document_id"],operator:ContainsAny,valueTextArray:[%s]}`, strings.Join(quoted, ",")))
	}
	if len(operands) == 0 {
		return ""
	}
	if len(operands) == 1 {
		return "where:" + operands[0]
	}
	return "where:{operator:And,operands:[" + strings.Join(operands, ",") + "]}"
}
