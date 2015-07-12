package httpd

import (
	"bufio"
	"fmt"
	"github.com/Symantec/Dominator/sub/scanner"
	"io"
	"net"
	"net/http"
	"net/rpc"
)

type HtmlWriter interface {
	WriteHtml(writer io.Writer)
}

var onlyHtmler HtmlWriter
var onlyFsh *scanner.FileSystemHistory

type Subd int

func (t *Subd) Poll(generation uint64, reply *scanner.FileSystem) error {
	fs := onlyFsh.FileSystem()
	if fs != nil && generation != onlyFsh.GenerationCount() {
		*reply = *fs
	}
	return nil
}

func StartServer(portNum uint, fsh *scanner.FileSystemHistory) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		return err
	}
	onlyHtmler = fsh
	onlyFsh = fsh
	http.HandleFunc("/", onlyHandler)
	subd := new(Subd)
	rpc.Register(subd)
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	return nil
}

func onlyHandler(w http.ResponseWriter, req *http.Request) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()
	fmt.Fprintln(writer, "<title>subd status page</title>")
	fmt.Fprintln(writer, "<body>")
	fmt.Fprintln(writer, "<center>")
	fmt.Fprintln(writer, "<h1>subd status page</h1>")
	fmt.Fprintln(writer, "</center>")
	fmt.Fprintln(writer, "<h3>")
	onlyHtmler.WriteHtml(writer)
	fmt.Fprintln(writer, "</body>")
}
