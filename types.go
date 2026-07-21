package mdparse

type ContentType string

var (
	TypeEmpty      ContentType = ""
	TypeHeader     ContentType = "header"
	TypeCode       ContentType = "code"
	TypeContent    ContentType = "content"
	TypeList       ContentType = "list"
	TypeBlockquote ContentType = "blockquote"
	TypeImage      ContentType = "image"
	TypeHR         ContentType = "hr"
)

type Paragraph struct {
	Header      string
	HeaderLevel int
	Body        string
	Type        ContentType
	Order       int
}
