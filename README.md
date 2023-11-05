# discord-niconico-comment-generator

## Description

discord のチャンネルに投稿されたメッセージを [NiCommentGenerator](https://github.com/totoraj930/NiCommentGenerator) で利用できる comment.xml として出力するためのツール。

## Build

```bash
GOOS=windows CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build
```
