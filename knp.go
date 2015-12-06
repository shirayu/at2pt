package at2pt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

//ConvertKNP processes texts in Cabocha format
func ConvertKNP(infile *os.File, outfile *os.File, mode Mode) (err error) {

	//file open
	reader := bufio.NewReader(infile)

	var buffer bytes.Buffer
	for {
		myline, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF || bytes.HasPrefix(myline, []byte("# ")) {
			doc := buffer.String()
			if len(doc) != 0 {
				var plainLine string
				var myerr error
				switch mode {
				case TOKENIZED:
					plainLine, myerr = GetTokensFromKNP(doc)
				default:
					plainLine, myerr = GetPlainTextsFromKNP(doc)
				}

				if myerr != nil {
					fmt.Fprintf(os.Stderr, "%s [%s]\n", myerr, plainLine)
				} else if len(plainLine) != 0 {
					fmt.Fprintln(outfile, plainLine)
				}

				if err == io.EOF {
					return nil
				}
			}

			buffer.Reset()
		}
		buffer.Write(myline)
		buffer.WriteString("\n")
	}
}

func getToken(line string) (string, error) {
	if strings.Contains(line, " 数詞 ") {
		return "<数量>", nil
	}

	//TODO how to deal with conjugation morpheme
	const normAttr = "<正規化代表表記:"
	normPos := strings.Index(line, normAttr)
	if normPos >= 0 {
		tail := line[normPos+len(normAttr):]
		tailPos := strings.Index(tail, ">")
		return tail[:tailPos], nil
	}

	i := strings.Index(line, " ")
	if i >= 0 {
		return line[:i], nil
	}

	return "", errors.New("Format error")
}

func isConnectTarget(_lines *[]string, lineid int) bool {
	lines := *_lines
	line := lines[lineid]
	if strings.Contains(line, "<一文字漢字>") {
		if strings.Contains(lines[lineid+1], "<文節主辞>") {
			return true
		}
	}
	return false
}

//GetPlainTextsFromKNP returns plain texts from KNP format texts
func GetPlainTextsFromKNP(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var out bytes.Buffer

	for _, line := range lines {
		if strings.HasPrefix(line, "# S-ID") { //sentence
		} else if line == "EOS" {
			break
		} else if strings.HasPrefix(line, "* ") { //buntsetu phrase
		} else if strings.HasPrefix(line, "+ ") { //basic phrase
		} else { //tokens
			i := strings.Index(line, " ")
			if i >= 0 {
				out.WriteString(line[:i])
			}
		}
	}

	return out.String(), nil
}

//GetTokensFromKNP returns tokenized words from KNP format texts
func GetTokensFromKNP(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var out bytes.Buffer
	noNextSpace := true
	isFirstTokenInBunsetsu := true

	for lineid, line := range lines {
		if strings.HasPrefix(line, "# S-ID") { //sentence

		} else if line == "EOS" {
			break

		} else if strings.HasPrefix(line, "* ") { //buntsetu phrase
			isFirstTokenInBunsetsu = true

		} else if strings.HasPrefix(line, "+ ") { //basic phrase
			if !isFirstTokenInBunsetsu && isConnectTarget(&lines, lineid) {
				out.WriteString("+")
				noNextSpace = true
			}

		} else { //tokens
			token, err := getToken(line)
			if err != nil {
				return lines[0], err
			}

			if token != "　" {
				if noNextSpace {
					noNextSpace = false
				} else if strings.Contains(line, "助数辞 ") && strings.Contains(lines[lineid-1], " 数詞 ") {
					out.WriteString("+")
				} else {
					out.WriteString(" ")
				}
				out.WriteString(token)
			}

			isFirstTokenInBunsetsu = false
		}
	}

	return out.String(), nil
}
