package types

type Track struct {
	ID   string `xml:"id,attr"`
	Path string `xml:"path,attr,omitempty"`
}
