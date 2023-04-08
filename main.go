package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <url>")
		os.Exit(1)
	}
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL %s: %v\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		os.Exit(1)
	}

	fetched, err := countLines(bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting lines: %v\n", err)
		os.Exit(1)
	}

	doc, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HTML content: %v\n", err)
		os.Exit(1)
	}

	issues := make(map[string]int)
	checkLabel(doc, &issues)
	checkAltAttribute(doc, &issues)
	checkButtonRole(doc, &issues)
	checkForm(doc, &issues)
	checkHiddenInput(doc, &issues)
	checkImageInput(doc, &issues)
	checkSubmitInput(doc, &issues)
	checkTableHeader(doc, &issues)
	checkTableScope(doc, "row", &issues)
	checkList(doc, &issues)
	checkHeading(doc, &issues)
	checkTitle(doc, &issues)
	checkFormLabel(doc, &issues)
	checkListItems(doc, &issues)
	checkListLabel(doc, &issues)
	checkListIndentation(doc, 0, &issues)
	checkFormSubmit(doc, &issues)
	checkFormMethod(doc, &issues)
	checkFormEnctype(doc, &issues)
	checkFormValidity(doc, &issues)
	checkArias(doc, &issues)
	checkImgAlt(doc, &issues)
	checkAnchorHref(doc, &issues)

	fmt.Printf("Accessibility report for URL: %s\n", url)
	totalIssues := 0
	if len(issues) > 0 {
		for issue, count := range issues {
			fmt.Printf("- %s (x%d)\n", issue, count)
			totalIssues += count
		}
	} else {
		fmt.Println("No accessibility issues found.")
	}

	fmt.Printf("Total issues: %d\n", totalIssues)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting lines: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total lines fetched: %d\n", fetched)
}
func countLines(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)

	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func checkArias(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "img":
			if hasAriaAttribute(node, "label") && !hasAltAttribute(node) {
				(*issues)["Image element with aria-label attribute but without an alt attribute."]++
			}
		case "button":
			if hasAriaAttribute(node, "controls") && !hasAriaAttribute(node, "label") {
				(*issues)["Button element with aria-controls attribute but without an aria-label attribute."]++
			}
		case "input":
			if hasAriaAttribute(node, "placeholder") && !hasLabel(node) {
				(*issues)["Input element with aria-placeholder attribute but without a corresponding label element."]++
			}
		case "select":
			if hasAriaAttribute(node, "label") && !hasLabel(node) {
				(*issues)["Select element with aria-label attribute but without a corresponding label element."]++
			}
		case "textarea":
			if hasAriaAttribute(node, "placeholder") && !hasLabel(node) {
				(*issues)["Textarea element with aria-placeholder attribute but without a corresponding label element."]++
			}
		case "table":
			if hasAriaAttribute(node, "label") && !hasCaption(node) {
				(*issues)["Table element with aria-label attribute but without a corresponding caption element."]++
			}
		case "audio", "video":
			if hasAriaAttribute(node, "description") && !hasAriaAttribute(node, "label") {
				(*issues)[fmt.Sprintf("%s element with aria-description attribute but without an aria-label attribute.", node.Data)]++
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkArias(child, issues)
	}
}

func hasCaption(node *html.Node) bool {
	if node.Type == html.ElementNode && node.Data == "caption" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasCaption(child) {
			return true
		}
	}

	return false
}

func hasAriaAttribute(node *html.Node, attribute string) bool {
	for _, attr := range node.Attr {
		if strings.HasPrefix(attr.Key, "aria-") && attr.Key[5:] == attribute {
			return true
		}
	}
	return false
}

func hasAltAttribute(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "alt" {
			return true
		}
	}
	return false
}

func checkFormSubmit(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" {
		var hasSubmit bool
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "input" {
				for _, attr := range child.Attr {
					if attr.Key == "type" && attr.Val == "submit" {
						hasSubmit = true
					}
				}
			}
		}
		if !hasSubmit {
			(*issues)["Form element without a submit button."]++
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkFormSubmit(child, issues)
	}
}

func checkFormMethod(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" {
		var hasMethod bool
		for _, attr := range node.Attr {
			if attr.Key == "method" {
				hasMethod = true
			}
		}
		if !hasMethod {
			(*issues)["Form element without a method attribute."]++
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkFormMethod(child, issues)
	}
}

func checkFormEnctype(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" {
		var hasEnctype bool
		for _, attr := range node.Attr {
			if attr.Key == "enctype" {
				hasEnctype = true
			}
		}
		if !hasEnctype {
			(*issues)["Form element without an enctype attribute."]++
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkFormEnctype(child, issues)
	}
}

func checkFormValidity(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" {
		if !hasFormValidation(node) {
			(*issues)["Form element without proper validation attributes."]++
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkFormValidity(child, issues)
	}
}

func hasFormValidation(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "novalidate" {
			return false
		}
		if attr.Key == "validate" {
			return true
		}
	}
	return false
}

func checkInputType(node *html.Node, inputType string, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "input" && hasInputType(node, inputType) && !hasLabel(node) {
		(*issues)[fmt.Sprintf("Input element with type '%s' without a corresponding label element.", inputType)]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkInputType(child, inputType, issues)
	}
}

func checkTableHeader(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "th" && !hasTableHeader(node.Parent) {
		(*issues)["Table header cell without a corresponding th or td parent element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkTableHeader(child, issues)
	}
}

func checkTableScope(node *html.Node, scope string, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "th" && !hasTableScope(node, scope) {
		(*issues)[fmt.Sprintf("Table header cell without a '%s' scope attribute.", scope)]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkTableScope(child, scope, issues)
	}
}

func checkList(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && (node.Data == "ul" || node.Data == "ol") && !hasListItem(node) {
		(*issues)[fmt.Sprintf("%s element without any corresponding list item elements.", node.Data)]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkList(child, issues)
	}
}

func checkHeading(node *html.Node, issues *map[string]int) {
	if hasHeading(node) && !hasTitle(node.Parent) {
		(*issues)["Heading element without a corresponding title element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkHeading(child, issues)
	}
}

func checkTitle(node *html.Node, issues *map[string]int) {
	if hasTitle(node) && node.Parent != nil && (node.Parent.Type != html.ElementNode || node.Parent.Data != "head") {
		(*issues)["Title element not inside a head element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkTitle(child, issues)
	}
}

func checkLabel(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && (node.Data == "input" || node.Data == "select" || node.Data == "textarea") && !hasLabel(node) {
		(*issues)[fmt.Sprintf("%s element without a corresponding label element.", node.Data)]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkLabel(child, issues)
	}
}

func checkAltAttribute(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "img" && !hasAltAttribute(node) {
		(*issues)["Image element without an alt attribute."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkAltAttribute(child, issues)
	}
}

func checkButtonRole(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && hasButtonRole(node) && !hasLabel(node) {
		(*issues)["Button element without a corresponding label element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkButtonRole(child, issues)
	}
}

func checkForm(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" && !hasLabel(node) {
		(*issues)["Form element without any corresponding label elements."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkForm(child, issues)
	}
}
func checkHiddenInput(node *html.Node, issues *map[string]int) {
	checkInputType(node, "hidden", issues)
}

func checkImageInput(node *html.Node, issues *map[string]int) {
	checkInputType(node, "image", issues)
}

func checkSubmitInput(node *html.Node, issues *map[string]int) {
	checkInputType(node, "submit", issues)
}
func checkImgAlt(node *html.Node, issues *map[string]int) {
	if !hasAltAttribute(node) {
		line := getLine(node)
		(*issues)[fmt.Sprintf("Image element without an alt attribute. Line: %s", line)]++
	}
}

func checkAnchorHref(node *html.Node, issues *map[string]int) {
	if !hasHrefAttribute(node) {
		line := getLine(node)
		(*issues)[fmt.Sprintf("Anchor element without an href attribute. Line: %s", line)]++
	}
}
func hasHrefAttribute(node *html.Node) bool {
	for _, a := range node.Attr {
		if a.Key == "href" {
			return true
		}
	}
	return false
}

func getLine(node *html.Node) string {
	if node != nil {
		if node.Type == html.ElementNode || node.Type == html.CommentNode {
			if node.Parent != nil {
				line := getLine(node.Parent)
				if line != "" {
					return line
				}
			}
			for _, attr := range node.Attr {
				if attr.Key == "line" {
					return fmt.Sprintf("%s:%s", node.Data, attr.Val)
				}
			}
		}
		return getLine(node.PrevSibling)
	}
	return ""
}

func hasLabel(node *html.Node) bool {
	if node.Data == "label" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasLabel(child) {
			return true
		}
	}

	return false
}

func hasButtonRole(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "role" && attr.Val == "button" {
			return true
		}
	}
	return false
}

func hasInputType(node *html.Node, inputType string) bool {
	if node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == inputType {
				return true
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasInputType(child, inputType) {
			return true
		}
	}

	return false
}

func hasTableHeader(node *html.Node) bool {
	if node.Data == "th" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasTableHeader(child) {
			return true
		}
	}

	return false
}

func hasTableScope(node *html.Node, scope string) bool {
	if node.Data == "th" {
		for _, attr := range node.Attr {
			if attr.Key == "scope" && attr.Val == scope {
				return true
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasTableScope(child, scope) {
			return true
		}
	}

	return false
}

func hasList(node *html.Node) bool {
	if node.Data == "ul" || node.Data == "ol" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasList(child) {
			return true
		}
	}

	return false
}

func hasListItem(node *html.Node) bool {
	if node.Data == "li" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasListItem(child) {
			return true
		}
	}

	return false
}

func checkFormLabel(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "form" && !hasLabel(node) {
		(*issues)["Form element without a corresponding label element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkFormLabel(child, issues)
	}
}

func checkListItems(node *html.Node, issues *map[string]int) {
	if node.Type == html.ElementNode && node.Data == "li" && !hasList(node.Parent) {
		(*issues)["List item element without a corresponding ol or ul parent element."]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkListItems(child, issues)
	}
}

func checkListLabel(node *html.Node, issues *map[string]int) {
	if (node.Data == "ul" || node.Data == "ol") && !hasLabel(node) {
		(*issues)[fmt.Sprintf("%s element without a corresponding label element.", node.Data)]++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkListLabel(child, issues)
	}
}

func checkListIndentation(node *html.Node, depth int, issues *map[string]int) {
	if node.Type == html.ElementNode && (node.Data == "ul" || node.Data == "ol") {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "li" {
				for i := 0; i < depth; i++ {
					if child.PrevSibling == nil || child.PrevSibling.Data != "li" {
						(*issues)[fmt.Sprintf("%s element with inconsistent indentation.", node.Data)]++
						break
					}
				}
				checkListIndentation(child, depth+1, issues)
			}
		}
	}
}

func hasHeading(node *html.Node) bool {
	if node.Type == html.ElementNode && node.Data[0] == 'h' && node.Data[1] >= '1' && node.Data[1] <= '6' {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasHeading(child) {
			return true
		}
	}

	return false
}

func hasTitle(node *html.Node) bool {
	if node == nil {
		return false
	}
	if node.Type == html.ElementNode && node.Data == "title" {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasTitle(child) {
			return true
		}
	}

	return false
}
