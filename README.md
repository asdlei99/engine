# Engine

[![Build Status](https://travis-ci.org/osrtss/engine.svg)](https://travis-ci.org/osrtss/engine) [![Coverage Status](https://coveralls.io/repos/github/osrtss/engine/badge.svg?branch=master)](https://coveralls.io/github/osrtss/engine?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/osrtss/engine)](https://goreportcard.com/report/github.com/osrtss/engine) [![GoDoc](https://godoc.org/github.com/osrtss/engine?status.svg)](https://godoc.org/github.com/osrtss/engine)

A generic transport agnostic stream layer

It aims to be a generic stream based network framework.

## Arch

```
+-----------+
|   Codec   | pluggable: json, protobuf, tlv, rtp, flv
+-----------+
|   Stream  | io.ReadWriter
+-----------+
| Transport | pluggable: tcp, udp, unix, http, rtmp, rtsp
+-----------+
```

* Codec: contains avcodec and avformat.
* Transport: contains tcp/udp layer and up layer like http, rtmp.

## Features
- [ ] Support "net/context".
- [ ] Add sesssion support.
- [ ] Add codec support.
- [ ] Add transport support.

## RFCs
* 
