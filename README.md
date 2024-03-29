# deeplgobindings ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg) [![Unit test](https://github.com/PineiroHosting/deeplgobindings/actions/workflows/go.yml/badge.svg)](https://github.com/PineiroHosting/deeplgobindings/actions/workflows/go.yml) ![GoDoc](https://godoc.org/github.com/PineiroHosting/deeplgobindings?status.svg) ![Go Report Card](https://goreportcard.com/badge/github.com/PineiroHosting/deeplgobindings)

Low-level bindings for version 2 of the DeepL translation API - https://www.deepl.com/api.html

## Description

### About DeepL
> DeepL is a deep learning company that develops AI systems for languages. The company, based in Cologne, Germany, was founded in 2009 as Linguee, and introduced the first internet search engine for translations. Linguee has answered over 10 billion queries from more than 1 billion users.
>
> In the summer of 2017, the company introduced DeepL Translator, a free machine translation system that produces translations of unprecedented quality.

Source: [www.deepl.com](https://www.deepl.com/publisher.html)

### This project

deeplgobindings allows automated interaction with DeepL translation's API in Golang.

## Features

The following list contains all features which are/should be supported (in the future):
- [x] translate Function (*/v2/translate*)
- [x] usage Function (*/v2/usage*)
- [x] document translate Function (*/v2/document*)
- [x] support POST and GET request methods (including file upload with multipart)
- [ ] implement DeepL API's limitation rules

## Usage

For Documentation see [https://godoc.org/github.com/PineiroHosting/deeplgobindings](https://godoc.org/github.com/PineiroHosting/deeplgobindings).

If you are interested in some examples, see the examples directory in this repository.

## Contribution

Feel free to contribute and help this project to grow. You can also just suggest features/enhancements.
