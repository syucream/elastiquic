# elastiquic

An elastic QUIC client and benchmark tool written in golang.

## quickstart

* Prepare goquic

```
$ go get -u -d github.com/devsisters/goquic
$ cd $GOPATH/src/github.com/devsisters/goquic/
$ ./build_libs.sh
```

* Write definitions

```
$ vim definitions.json
```

* run elastiquic

```
$ go run elastiquic.go
..

Total requests: 2, successed: 2, failed: 0
```
