package main

import (
	"net/http/httptest"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ReportSuite struct{}

var _ = Suite(&ReportSuite{})

func (s *ReportSuite) TestReportProcess(c *C) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("lb-author", "test-author")
	req.Header.Set("lb-req-cnt", "1")

	r := make(Report)

	r.Process(req)
	c.Assert(reflect.DeepEqual(r["test-author"], []string{"1"}), Equals, true)

	req.Header.Set("lb-req-cnt", "2")
	r.Process(req)
	c.Assert(reflect.DeepEqual(r["test-author"], []string{"1", "2"}), Equals, true)

	req.Header.Set("lb-author", "test-len")
	for i := 0; i < 103; i++ {
		req.Header.Set("lb-req-cnt", "test-len")
		r.Process(req)
	}
	c.Assert(len(r["test-len"]), Equals, reportMaxLen)
}
