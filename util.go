package crawler

import "net/url"

func ToAbsUrl(base *url.URL,rurl string) (string,error) {
  t,err := url.Parse(rurl)
  if err != nil {
    return "", err
  }
  return base.ResolveReference(t).String(),nil

}
