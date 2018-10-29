# Contributing

## Introduction

Thank you for your interest in deeplgobindings. We all together make deeplgobindings a good API implementation for go.

These guidelines help you communicating and contributing in a efficient and stress-free manner. Or... we hope so.
Anyway, please keep in mind, that we are just humans and do something wrong. So keep respectful and watch the [code of conduct](https://github.com/PineiroHosting/deeplgobindings/blob/master/CODE_OF_CONDUCT.md). Thank you!

## Requirements

* Before creating a pull request, remember that there are more people like you who would like to help. An idea gets better when a lot of people bring in their suggestions. Go and create an issue where you may discuss feature requests.
* Ensure, that your code is cross-platform compatible.

**Working on your first Pull Request?** You can learn how from this *free* series [How to Contribute to an Open Source Project on GitHub](https://egghead.io/series/how-to-contribute-to-an-open-source-project-on-github) 

## Issues
#### Suggestion of a feature or an enhancement
When suggesting a feature, your issue should at least contain the following information:
* short explanation
* required/suggested dependencies
* size of feature
* importance of this issue

#### Report of errors or bugs
Oh dear, development is like being the detective in a crime movie where you are also the murderer. In this way, Filipe Fortes described debugging. But the movie should end... or this part, another will start soon.

When reporting errors or bugs, your issue should at least contain the following information:
* go version
* operating system information (like version, architecture, type of os, ...)
* error log
* any additional information of your environment which could be important for this issue. Keep in mind, that normally issues are discussed after they were opened and you may have to provide further information.

## Pull requests

In your pull request, you should:
* describe the content, you changed and reference to a issue this request deals with, if there is one
* list dependencies you have added to the project - run `dep ensure` if you added something (see [godep](https://github.com/golang/dep/) for more information)
* reformat your files, just use `go fmt .` for this

## Code conventions
This application follows the general Golang code convention as well as its commentary conventions. We use simple tabs (`\t`) as the tab intent.

## Questions?
We are glad to help, just write an issue.