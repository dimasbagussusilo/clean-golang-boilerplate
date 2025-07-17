package helper

import (
	"fmt"
	"strings"
)

type (
	// DOMElement ...
	DOMElement struct {
		name      string
		style     map[string]interface{}
		childs    []interface{}
		text      string
		SingleTag bool
	}
)

// InsertChildBefore
//
//	insert child on first of element
func (dom *DOMElement) InsertChildBefore(node DOMElement) {
	child := make([]interface{}, 0)

	child = append(child, node)
	child = append(child, dom.childs...)

	dom.childs = child
}

// InsertChild
//
//	insert child on last of element
func (dom *DOMElement) InsertChild(node DOMElement) {
	dom.childs = append(dom.childs, node)
}

// SetStyle
//
//	add style on element
func (dom *DOMElement) SetStyle(name string, value interface{}) {
	switch len(dom.style) {
	case 0:
		mp := make(map[string]interface{})
		mp[name] = value
		dom.style = mp
	default:
		dom.style[name] = value
	}
}

// AddText
//
//	add text value on element
func (dom *DOMElement) AddText(value string) {
	dom.text += value
}

// SetName
//
//	Set name element or type of element for generate
func (dom *DOMElement) SetName(name string) {
	dom.name = name
}

// OuterHTML
//
//	set outer html with string of element
func (dom DOMElement) OuterHTML() string {
	var (
		html, childstr, style string

		styles = make([]string, 0)
	)

	if len(dom.name) > 0 {
		if len(dom.style) > 0 {
			for key, val := range dom.style {
				styles = append(styles, fmt.Sprintf("%s:%v", key, val))
			}

			style = strings.Join(styles, ";")
			style = fmt.Sprintf("style=\"%s\"", style)
		}

		if len(dom.childs) > 0 {
			for _, child := range dom.childs {
				if el, ok := child.(DOMElement); ok {
					childstr += el.OuterHTML()
				}
			}
		}

		elName := strings.ToLower(dom.name)
		html = fmt.Sprintf("</%s>", elName)
		if !dom.SingleTag {
			html = fmt.Sprintf("<%s %s>%v%v</%s>", elName, style, childstr, dom.text, elName)
		}
	}

	return html
}

// RedBoldText
//
//	Create new red bold text element
func (dom DOMElement) RedBoldText(text string) DOMElement {
	dom.SetName("B")
	dom.AddText(text)
	dom.SetStyle("color", "red")

	return dom
}

// BoldText
//
//	Create new bold text of element
func (dom DOMElement) BoldText(text string) DOMElement {
	dom.SetName("B")
	dom.AddText(text)

	return dom
}

// BoldTextScanf ...
func (dom DOMElement) BoldTextScanf(text string, format ...interface{}) DOMElement {
	dom.SetName("B")
	dom.AddText(fmt.Sprintf(text, format...))

	return dom
}

// RedBoldTextScanf ...
func (dom DOMElement) RedBoldTextScanf(text string, format ...interface{}) DOMElement {
	dom.SetName("B")
	dom.AddText(fmt.Sprintf(text, format...))
	dom.SetStyle("color", "red")

	return dom
}

func GenerateErrorMessage(num int64, errs ...string) string {
	var (
		message                        string
		div, span, tred, tbold, br, ol DOMElement
	)

	div.SetName("DIV")
	span.SetName("SPAN")
	br.SetName("BR")
	ol.SetName("OL")
	br.SingleTag = true

	tred = tred.RedBoldText("Data Tidak Valid!")
	tbold = tbold.BoldText(fmt.Sprintf("%d", num))
	span.AddText(fmt.Sprintf(
		"%s, Row %s, tidak dapat diproses:",
		tred.OuterHTML(),
		tbold.OuterHTML(),
	))

	div.InsertChild(span)
	div.InsertChild(br)

	for _, er := range errs {
		var li DOMElement
		li.SetName("LI")
		li.AddText(er)

		ol.InsertChild(li)
	}

	div.InsertChild(ol)
	message = div.OuterHTML()

	return message
}
