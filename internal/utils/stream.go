package utils

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sugawarayuuta/sonnet"
)

type toEncoder[E any] struct {
	encoder *sonnet.Encoder
	resp    *echo.Response
	count   int
}

func InitEncoder[E any](resp *echo.Response) toEncoder[E] {
	resp.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	resp.WriteHeader(http.StatusOK)

	encoder := sonnet.NewEncoder(resp)

	resp.Writer.Write([]byte("{"))
	resp.Writer.Write([]byte(`"result":`))
	resp.Writer.Write([]byte("["))

	return toEncoder[E]{
		encoder: encoder,
		resp:    resp,
	}
}

func (e *toEncoder[E]) Iterator(m E) (err error) {
	if e.count != 0 {
		e.resp.Writer.Write([]byte(","))
	}

	if err = e.encoder.Encode(m); err != nil {
		return
	}

	e.resp.Flush()

	e.count++

	return nil
}

func (e *toEncoder[E]) EndStream() (err error) {
	if _, err = e.resp.Writer.Write([]byte("]")); err != nil {
		return
	}

	e.resp.Writer.Write([]byte(","))
	e.resp.Writer.Write([]byte(`"count":`))
	e.resp.Writer.Write([]byte(strconv.Itoa(e.count)))

	if _, err = e.resp.Writer.Write([]byte("}")); err != nil {
		return
	}

	return nil
}
