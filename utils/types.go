package utils

type NextData struct {
	Props struct {
		PageProps struct {
			Know KnowPage `json:"know"`
		} `json:"pageProps"`
	} `json:"props"`
}

type KnowPage struct {
	Description string   `json:"description"`
	Documents   []PDFDoc `json:"documents"`
}

type PDFDoc struct {
	URL string `json:"contentUrl"`
}

type SearchResponse struct {
	Content []PDFItem `json:"content"`
}

type PDFItem struct {
	Know  KnowDetails `json:"know"`
	Score float64     `json:"score"`
}

type KnowDetails struct {
	ID     string `json:"uuid"`
	Title  string `json:"title"`
	Likes  int64  `json:"likes"`
	Thumb  string `json:"thumbnailLargeUrl"`
	Knower Knower `json:"knower"`
}

type Knower struct {
	User UserDetails `json:"user"`
}

type UserDetails struct {
	Name string `json:"name"`
}

type PDF struct {
	ID, Title, ThumbURL, Knower string
	Score                       float64
	Likes                       int64
}
