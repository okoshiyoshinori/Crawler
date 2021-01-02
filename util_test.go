package crawler

import (
	"net/url"
	"testing"
)



func TestToAbsUrl(m *testing.T) {
  base,_ := url.Parse("http://google.co.jp")
  path := "../ookoshi/ooko"
  result,_ := ToAbsUrl(base,path)
  m.Log(result == "http://google.co.jp/ookoshi/ooko")
}

