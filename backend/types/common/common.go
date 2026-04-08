package common

type I18nObject struct {
	/*
	   Model class for i18n object.
	*/

	ZhHans string `json:"zh_Hans" yaml:"zh_Hans"`
	EnUS   string `json:"en_US" yaml:"en_US"`
	PtBR   string `json:"pt_BR" yaml:"pt_BR"`
	JaJP   string `json:"ja_JP" yaml:"ja_JP"`
}

func (o *I18nObject) ToDict() map[string]any {
	return map[string]any{"zh_Hans": o.ZhHans, "en_US": o.EnUS, "pt_BR": o.PtBR, "ja_JP": o.JaJP}
}
