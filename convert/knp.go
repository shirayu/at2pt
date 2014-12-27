package convert

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ConvertKNP(infile *os.File, outfile *os.File, tokenize bool) (err error) {

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
				if tokenize {
					plainLine, myerr = ParseKNPTokenized(doc)
				} else {
					plainLine, myerr = ParseKNP(doc)
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

	return err
}

func getToken(line string) (string, error) {
	if strings.Contains(line, " 数詞 ") {
		return "<数量>", nil
	}

	//TODO how to deal with conjugation morpheme
	const norm_attr = "<正規化代表表記:"
	norm_pos := strings.Index(line, norm_attr)
	if norm_pos >= 0 {
		tail := line[norm_pos+len(norm_attr):]
		tail_pos := strings.Index(tail, ">")
		return tail[:tail_pos], nil
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
func ParseKNP(data string) (string, error) {
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
func ParseKNPTokenized(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var out bytes.Buffer
	no_next_space := true
	is_first_token_in_bunsetsu := true

	for lineid, line := range lines {
		if strings.HasPrefix(line, "# S-ID") { //sentence

		} else if line == "EOS" {
			break

		} else if strings.HasPrefix(line, "* ") { //buntsetu phrase
			is_first_token_in_bunsetsu = true

		} else if strings.HasPrefix(line, "+ ") { //basic phrase
			if !is_first_token_in_bunsetsu && isConnectTarget(&lines, lineid) {
				out.WriteString("+")
				no_next_space = true
			}

		} else { //tokens
			token, err := getToken(line)
			if err != nil {
				return lines[0], err
			}

			if token != "　" {
				if no_next_space {
					no_next_space = false
				} else if strings.Contains(line, "助数辞 ") && strings.Contains(lines[lineid-1], " 数詞 ") {
					out.WriteString("+")
				} else {
					out.WriteString(" ")
				}
				out.WriteString(token)
			}

			is_first_token_in_bunsetsu = false
		}
	}

	return out.String(), nil
}
