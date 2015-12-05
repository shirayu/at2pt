package at2pt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ConvertCaboCha(infile *os.File, outfile *os.File, mode Mode) (err error) {

	//file open
	r := bufio.NewReader(infile)

	scanner := bufio.NewScanner(r)
	first_token := true
	for scanner.Scan() {
		line := scanner.Text()

		//         if line == "EOS" {
		if line == "EOS" || line == "EOS\t" {
			fmt.Fprintf(outfile, "\n")
			first_token = true
		} else if strings.HasPrefix(line, "* ") {
		} else {
			if mode == TOKENIZED {
				if first_token {
					first_token = false
				} else {
					fmt.Fprintf(outfile, " ")
				}
			}
			i := strings.Index(line, "\t")
			if i >= 0 {
				fmt.Fprintf(outfile, line[:i])
			}
		}
	}

	if err == nil || err == io.EOF {
		return nil
	}
	return scanner.Err()
}
