package handler

import (
	"bytes"
	"fmt"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/myrest/httpx"
	"github.com/tiptok/gocomm/xtime"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

const slowThreshold = time.Millisecond * 500

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.Now()
		lrw := LoggedResponseWriter{
			w:    w,
			r:    r,
			code: http.StatusOK,
		}

		var dup io.ReadCloser
		r.Body, dup = DupReadCloser(r.Body)
		next.ServeHTTP(&lrw, r)
		r.Body = dup
		logBrief(r, lrw, timer)
	})
}

func DetailedLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.Now()
		var buf bytes.Buffer
		lrw := newDetailLoggedResponseWriter(&LoggedResponseWriter{
			w:    w,
			r:    r,
			code: http.StatusOK,
		}, &buf)

		var dup io.ReadCloser
		r.Body, dup = DupReadCloser(r.Body)
		next.ServeHTTP(lrw, r)
		r.Body = dup
		logDetails(r, lrw, timer)
	})
}

type LoggedResponseWriter struct {
	w    http.ResponseWriter
	r    *http.Request
	code int
	buf  *bytes.Buffer
}

func (w *LoggedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *LoggedResponseWriter) Write(bytes []byte) (int, error) {
	//w.buf.Write(bytes)
	return w.w.Write(bytes)
}

func (w *LoggedResponseWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
	w.code = code
}

type DetailLoggedResponseWriter struct {
	writer *LoggedResponseWriter
	buf    *bytes.Buffer
}

func newDetailLoggedResponseWriter(writer *LoggedResponseWriter, buf *bytes.Buffer) *DetailLoggedResponseWriter {
	return &DetailLoggedResponseWriter{
		writer: writer,
		buf:    buf,
	}
}

func (w *DetailLoggedResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *DetailLoggedResponseWriter) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *DetailLoggedResponseWriter) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

func dumpRequest(r *http.Request) string {
	reqContent, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	} else {
		return string(reqContent)
	}
}

func logBrief(r *http.Request, response LoggedResponseWriter, start time.Time) {
	var buf bytes.Buffer
	duration := time.Since(start)
	buf.WriteString(fmt.Sprintf("[HTTP] %d | %s | %s | %s | %s",
		response.code, r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent(), xtime.ReprOfDuration(duration)))
	if duration > slowThreshold {
		log.Warn(fmt.Sprintf("[HTTP] %d | %s | %s | %s | slowcall(%s)",
			response.code, r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent(), xtime.ReprOfDuration(duration)))
	}

	ok := isOkResponse(response.code)
	if !ok {
		buf.WriteString(fmt.Sprintf("\n%s\n", dumpRequest(r)))
	}
	//respBuf := response.buf.Bytes()
	//if len(respBuf) > 0 {
	//	buf.WriteString(fmt.Sprintf("<= %s", respBuf))
	//}
	if ok {
		log.Info(buf.String())
	} else {
		log.Error(buf.String())
	}
}

func logDetails(r *http.Request, response *DetailLoggedResponseWriter, start time.Time) {
	var buf bytes.Buffer
	duration := time.Since(start)
	buf.WriteString(fmt.Sprintf("[HTTP] %d | %s | %s\n=> %s\n",
		response.writer.code, r.RemoteAddr, xtime.ReprOfDuration(duration), dumpRequest(r)))
	if duration > slowThreshold {
		log.Warn(fmt.Sprintf("[HTTP] %d | %s | slowcall(%s)\n=> %s\n",
			response.writer.code, r.RemoteAddr, xtime.ReprOfDuration(duration), dumpRequest(r)))
	}

	respBuf := response.buf.Bytes()
	if len(respBuf) > 0 {
		buf.WriteString(fmt.Sprintf("<= %s", respBuf))
	}

	log.Info(buf.String())
}

func isOkResponse(code int) bool {
	// not server error
	return code < http.StatusInternalServerError
	//return true
}

// The first returned reader needs to be read first, because the content
// read from it will be written to the underlying buffer of the second reader.
func DupReadCloser(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	return ioutil.NopCloser(tee), ioutil.NopCloser(&buf)
}
