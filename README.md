
# at2pt: Annotated Texts to Plain Texts


## What's this

- This tool converts texts annotated by NLP tools to plain texts
- For example, you can make [word2vec](https://code.google.com/p/word2vec/) models using the tokenized texts
    - Use option ``--tokenize``

## Usage
```
Usage of ./at2pt:
  -i, -input :  Input file name. - or no designation means STDIN.
  -o, -output: Output file name. - or no designation means STDOUT.
  -t, --tokenize  Enable tokenzize mode (default:false)
  -s, --style=    Input file style {KNP, MeCab, CaboCha} (default:KNP)
```

More options can be seen with ``at2pt -h``

## INSTALL

```
go get github.com/shirayu/at2pt/at2pt
```

## Requirements
- You should use UTF-8 texts for input


## Licence

- GNU General Public License version 3
- (C) Yuta Hayashibe 2014
