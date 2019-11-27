# go-scrape


```go
    DefaultOutputPath = "video"

    //如果需要使用代理请注册代理地址
    //RegisterProxy("https://localhost:10808")
    //RegisterProxy("http://localhost:10808")
	e := RegisterProxy("socks5://localhost:10808")
	if e != nil {
		return
    }
    //创建搜刮器
	grab2 := NewGrabJavbus()
	grab3 := NewGrabJavdb()
	scrape := NewScrape(GrabOption(grab2), GrabOption(grab3), OptimizeOption(true))

    //需要查找的番号：多次或单次调用皆可
	e = scrape.Find("abp-891")
	e = scrape.Find("abp-892")

    //遍历结果
	e = scrape.Range(func(key string, content Content) error {
		t.Log("key", key, "info", content)
		return nil
	})
	//或输出到DefaultOutputPath
	e = scrape.Output()
```