package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ClippingType string

const (
	Highlight ClippingType = "Highlight"
	Note      ClippingType = "Note"
	Marker    ClippingType = "Marker"
	Unknown   ClippingType = "Unknown"
)

type Clipping struct {
	clippingType ClippingType
	content      string
	start        int
	end          int
}

type Parser struct {
	source    string
	target    string
	Clippings map[string][]*Clipping
}

func NewParser(source, target string) *Parser {
	return &Parser{source: source, target: target, Clippings: make(map[string][]*Clipping)}
}

func (p *Parser) Parse() error {
	file, err := os.Open(p.source)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for {

		book := ""
		if scanner.Scan() {
			book = scanner.Text()
			if strings.HasPrefix(book, "==") {
				continue
			}
		} else {
			break
		}

		clip := &Clipping{}

		scanner.Scan()
		description := scanner.Text()

		clippingType, err := p.getClippingType(description)
		if err != nil {
			clip.clippingType = Unknown
		} else {
			clip.clippingType = clippingType
		}

		start, end, err := p.getPositions(description)
		if err != nil {
			continue
		}

		content := ""

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ReplaceAll(line, "\r", "")
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "==") {
				break
			}

			content += line + "\n"
		}

		if len(book) >= 3 && book[0] == 0xEF && book[1] == 0xBB && book[2] == 0xBF {
			book = book[3:]
		}

		clip.content = content
		clip.start = start
		clip.end = end

		book = strings.TrimRight(book, "\r\n")
		p.Clippings[book] = append(p.Clippings[book], clip)

	}

	return nil
}

// func (p *Parser) Print() {
// 	for book, clippings := range p.Clippings {
// 		fmt.Println(book)
// 		fmt.Println("clip count:", len(clippings))
// 		// for _, clip := range clippings {
// 		// 	fmt.Printf("Type: %v\n", clip.clippingType)
// 		// 	fmt.Printf("Content: %v\n", clip.content)
// 		// 	fmt.Println()
// 		// }
// 	}
// }

func (p *Parser) getPositions(str string) (int, int, error) {

	reg := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := reg.FindStringSubmatch(str)

	if len(matches) == 0 {
		reg := regexp.MustCompile(`(\d+)`)
		matches = reg.FindStringSubmatch(str)
		if len(matches) != 2 {
			return 0, 0, fmt.Errorf("invalid positions: %v", str)
		}
		pos, _ := strconv.Atoi(matches[0])
		return pos, pos, nil
	}

	if len(matches) == 4 {
		pos, err := strconv.Atoi(matches[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid position: %v", matches[1])
		}
		return pos, pos, nil
	}

	start, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start position: %v", matches[1])
	}

	end, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end position: %v", matches[2])
	}

	return start, end, nil

}

func (p *Parser) getClippingType(str string) (ClippingType, error) {

	reg := regexp.MustCompile(`^- (La subrayado|La nota|El marcador)`)
	matches := reg.FindStringSubmatch(str)
	if len(matches) > 0 {
		if matches[1] == "La subrayado" {
			return Highlight, nil
		} else if matches[1] == "La nota" {
			return Note, nil
		} else if matches[1] == "El marcador" {
			return Marker, nil
		}
	}
	return "", fmt.Errorf("invalid clipping type: %v", matches)
}
