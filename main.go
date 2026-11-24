package main

import (
  "bytes"
  "context"
  "encoding/binary"
  "encoding/hex"
  "io"
  "net"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"

  "github.com/coder/websocket"
)

var uuidBytes []byte
func parseUUID(s string) ([]byte, error) {
  s = strings.ReplaceAll(s, "-", "")
  b, err := hex.DecodeString(s)
  if err != nil {
    return nil, err
  }
  if len(b) != 16 {
    return nil, io.ErrUnexpectedEOF
  }
  return b, nil
}

func handleWS(w http.ResponseWriter, r *http.Request) {
  c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
  if err != nil {
    return
  }
  defer c.CloseNow()
  conn := websocket.NetConn(context.Background(), c, websocket.MessageBinary)
  defer conn.Close()

  one := make([]byte, 1)

if _, err := io.ReadFull(conn, one); err != nil || one[0] != 0 {
    return
  }
  id := make([]byte, 16)
  if _, err := io.ReadFull(conn, id); err != nil {
    return
  }
  if !bytes.Equal(id, uuidBytes) {
    return
  }
  if _, err := io.ReadFull(conn, one); err != nil {
    return
  }
if one[0] > 0 {
    if _, err := io.CopyN(io.Discard, conn, int64(one[0])); err != nil {
      return
    }
  }
  if _, err := io.ReadFull(conn, one); err != nil || one[0] != 1 {
    return
  }
  pb := make([]byte, 2)
  if _, err := io.ReadFull(conn, pb); err != nil {
    return
  }

port := binary.BigEndian.Uint16(pb)
  if _, err := io.ReadFull(conn, one); err != nil {
    return
  }
  var host string
  switch one[0] {
  case 1:
    b := make([]byte, 4)
    if _, err := io.ReadFull(conn, b); err != nil {
      return
    }
    host = net.IP(b).String()
  case 2:

if _, err := io.ReadFull(conn, one); err != nil {
      return
    }
    b := make([]byte, int(one[0]))
    if _, err := io.ReadFull(conn, b); err != nil {
      return
    }
    host = string(b)
  case 3:
    b := make([]byte, 16)
    if _, err := io.ReadFull(conn, b); err != nil {
      return
    }
    host = net.IP(b).String()
  default:
    return
  }
  upstream, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), 10*time.Second)
  if err != nil {
    return
  }
  defer upstream.Close()
  if _, err := conn.Write([]byte{0, 0}); err != nil {
    return
  }
  go io.Copy(upstream, conn)
  io.Copy(conn, upstream)
}

func main() {
  u := os.Getenv("VLESS_UUID")
  if u == "" {
    panic("VLESS_UUID not set")
  }
  var err error
  uuidBytes, err = parseUUID(u)
  if err != nil {
    panic("bad VLESS_UUID")
  }
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

mux := http.NewServeMux()
  mux.HandleFunc("/raycdn", handleWS)
  srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: mux}
  if err := srv.ListenAndServe(); err != nil {
    panic(err)
  }
}







