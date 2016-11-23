# engine

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
