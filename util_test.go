package crawler

import (
	"net/url"
	"testing"
)

func TestToAbsUrl(m *testing.T) {
  base,_ := url.Parse("https://www.omecci.jp")
  path := "/chance/market/okutama_premium.html"
  result,_ := ToAbsUrl(base,path)
  m.Log(result == "https://www.omecci.jp/chance/market/okutama_premium.html")
}

