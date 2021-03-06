package main

import (
	"crypto/md5"
	"crypto/rc4"
	"flag"
	"io"
	"log"
	"net"
	"sync"
)

var (
	target = flag.String("target", "8.8.8.8:80", "target host")
)

type CryptoWriter struct {
	w      io.Writer
	cipher *rc4.Cipher
}

func NewCryptoWriter(w io.Writer, key string) io.Writer {
	md5sum := md5.Sum([]byte(key))
	cipher, err := rc4.NewCipher(md5sum[:])
	if err != nil {
		panic(err)
	}
	return &CryptoWriter{
		w:      w,
		cipher: cipher,
	}
}

func (w *CryptoWriter) Write(b []byte) (int, error) {
	buf := make([]byte, len(b))
	w.cipher.XORKeyStream(buf, b)
	return w.w.Write(buf)
}

type CryptoReader struct {
	r      io.Reader
	cipher *rc4.Cipher
}

func NewCryptoReader(r io.Reader, key string) io.Reader {
	md5sum := md5.Sum([]byte(key))
	cipher, err := rc4.NewCipher(md5sum[:])
	if err != nil {
		panic(err)
	}
	return &CryptoReader{
		r:      r,
		cipher: cipher,
	}
}

func (r *CryptoReader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	buf := b[:n]
	r.cipher.XORKeyStream(buf, buf)
	return n, err

}

func handleConn(conn net.Conn) {
	var remote net.Conn
	remote, err := net.Dial("tcp", *target)
	if err != nil {
		log.Println(err)
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		w := NewCryptoWriter(remote, "123456")
		defer wg.Done()
		io.Copy(w, conn)
		remote.Close()
	}()
	go func() {
		r := NewCryptoReader(remote, "123456")
		defer wg.Done()
		io.Copy(conn, r)
		conn.Close()
	}()
	wg.Wait()
}

func main() {
	listen, err := net.Listen("tcp", ":8021")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
