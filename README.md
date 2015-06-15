
# at2pt: Annotated Texts to Plain Texts


## What's this

- This tool converts texts annotated by NLP tools to plain texts
- For example, you can make [word2vec](https://code.google.com/p/word2vec/) models using the tokenized texts
    - Use option ``-m 1``

## Usage
```
Usage of ./at2pt:
  -h, --help    Show this help message
  -i, --input=  Input file name. - or no designation means STDIN (-)
  -o, --output= Output file name. - or no designation means STDOUT (-)
  -m, --mode=   Mode {0:PLAIN, 1:TOKENIZED} (0)
  -s, --style=  Input file style {KNP, MeCab, CaboCha} (KNP)
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
