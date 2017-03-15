package curscraft

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Tag struct {
	Name string
	Attr []Attribute
}

type Attribute struct {
	Key, Val string
}

type Argument struct {
	Key, Val string
}

func (t *Tag) GetAttr(attrName string) (val string, err error) {
	for _, attr := range t.Attr {
		if attr.Key == attrName {
			val = attr.Val
		}
	}
	return val, err
}

func GetURL(url string, tag string) []Tag {
	var tags []Tag // Tag slice
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err == nil {
		n, _ := html.Parse(resp.Body)
		parser(n, &tag, &tags)
	}
	return tags
}

func parser(n *html.Node, tag *string, tags *[]Tag) {
	//var t = *tag[0:1]
	var t = ""
	if n.Type == html.ElementNode {
		b, currentTag, nextTag := checkNextTag(*tag)
		switch t {
		case ".": // Class

		case "#": //ID

		default: // Tag
			t, args := checkSubArg(currentTag)
			if n.Data == t {
				test := setTag(n)

				// Sub Arguments for filtring tag
				argOk := true
				if len(args) > 0 {
					argOk = getValidAttributes(test, args)
				}

				if argOk {
					*tag = nextTag
					if !b {
						*tag = currentTag
						*tags = append(*tags, test)
					}
				}
			}
		}
	}

	for d := n.FirstChild; d != nil; d = d.NextSibling {
		parser(d, tag, tags)
	}
}

func checkSubArg(tag string) (newTag string, args []Argument) {
	var aArgs = strings.Split(tag, ".")
	newTag = aArgs[0]
	if len(aArgs) > 1 {
		for i := 1; i < len(aArgs); i++ {
			subArg := strings.Split(aArgs[i], "=")
			arg := new(Argument)
			arg.Key = subArg[0]
			arg.Val = subArg[1]
			args = append(args, *arg)
		}
	}
	return newTag, args
}

func checkNextTag(tag string) (bNextTag bool, currentTag string, nextTag string) {
	var i = strings.Index(tag, "|")
	currentTag = tag
	if i != -1 && len(tag[i+1:]) > 0 {
		currentTag = tag[0:i]
		nextTag = tag[i+1:]
		bNextTag = true
	}
	return bNextTag, currentTag, nextTag
}

func getValueAttribute(n *html.Node, attribute string) (value string) {
	for _, attr := range n.Attr {
		if attr.Key == attribute {
			value = attr.Val
		}
	}
	return value
}

func setTag(n *html.Node) (t Tag) {
	// Main Attributes
	t.Name = n.Data // Tag Type (a, input, div...)
	for _, attr := range n.Attr {
		tagAttr := Attribute{attr.Key, attr.Val}
		t.Attr = append(t.Attr, tagAttr)
	}

	// Text Attribute (If exist)
	nt := n.FirstChild
	if nt != nil && nt.Type == html.TextNode {
		tagAttr := Attribute{"text", nt.Data}
		t.Attr = append(t.Attr, tagAttr)
	}
	return t
}

func getValidAttributes(t Tag, args []Argument) bool {
	valid := make([]bool, len(args))
	for _, tagAttr := range t.Attr {
		for i, arg := range args {
			if tagAttr.Key == arg.Key && tagAttr.Val == arg.Val {
				valid[i] = true
			}
		}
	}

	var ok bool = true
	for _, v := range valid {
		if v == false {
			ok = false
		}
	}
	return ok
}
