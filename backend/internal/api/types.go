package api

// AnalyzeRequest は解析リクエストのボディを表します
type AnalyzeRequest struct {
	HTML      string `json:"html"`
	Selection string `json:"selection"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Auth      Auth   `json:"auth"`
}

// Auth は簡易認証情報を表します
type Auth struct {
	SharedSecret string `json:"shared_secret"`
	UserID       string `json:"user_id"`
}

// AnalyzeResponse は解析結果のレスポンスを表します
type AnalyzeResponse struct {
	Tokens []Token `json:"tokens"`
}

// Token は形態素分析された単語と読み仮名のペアを表します
type Token struct {
	BaseForm     string `json:"base_form"`
	Reading      string `json:"reading"`
	PartOfSpeech string `json:"part_of_speech"`
}
