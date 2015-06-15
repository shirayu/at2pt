package convert

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestConvertKNP(t *testing.T) {
	tests := []struct {
		input string
		gold  string
		mode  Mode
	}{
		{
			input: "knp_test/input.knp",
			gold:  "knp_test/plain_gold.txt",
			mode:  PLAIN,
		},
		{
			input: "knp_test/input.knp",
			gold:  "knp_test/tokenized_gold.txt",
			mode:  TOKENIZED,
		},
	}

	for _, test := range tests {
		inf, err := os.Open(test.input)
		if err != nil {
			t.Errorf("Error: %v", err)
			return
		}
		defer inf.Close()

		//make output
		//cf: http://stackoverflow.com/questions/10473800/
		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err = ConvertKNP(inf, os.Stdout, test.mode)
		if err != nil {
			t.Errorf("Error: %s", err)
			return
		}

		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		// back to normal state
		w.Close()
		os.Stdout = old // restoring the real stdout
		gotStdout := <-outC

		//get golds
		goldf, err := os.Open(test.gold)
		gold_reader := bufio.NewReader(goldf)
		if err != nil {
			t.Errorf("Error when open the gold file: %v", err)
			return
		}

		//Check
		for _, line := range strings.Split(gotStdout, "\n") {
			_gold_line, _, err := gold_reader.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				t.Errorf("Error when getting gold line: %s", err)
				return
			}

			gold_line := string(_gold_line)
			if line != gold_line {
				t.Errorf("got\n[%s]\nbut want\n[%s]\n", line, gold_line)
			}
		}
	}
}
