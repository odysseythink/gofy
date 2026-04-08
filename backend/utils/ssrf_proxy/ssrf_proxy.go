package ssrfproxy

// const (
// 	SSRFDefaultMaxRetries     = 3
// 	SSRFDefaultTimeOut        = 30
// 	SSRFDefaultConnectTimeOut = 10
// 	SSRFDefaultReadTimeOut    = 30
// 	SSRFDefaultWriteTimeOut   = 30
// 	BACKOFF_FACTOR            = 500
// )

// var (
// 	STATUS_FORCELIST = [5]int{429, 500, 502, 503, 504}
// )

// type rwTimeoutConn struct {
// 	*net.TCPConn
// 	read_timeout  time.Duration
// 	write_timeout time.Duration
// }

// func (this *rwTimeoutConn) Read(b []byte) (int, error) {
// 	if this.read_timeout > 0 {
// 		err := this.TCPConn.SetDeadline(time.Now().Add(this.read_timeout))
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return this.TCPConn.Read(b)
// }

// func (this *rwTimeoutConn) Write(b []byte) (int, error) {
// 	if this.write_timeout > 0 {
// 		err := this.TCPConn.SetDeadline(time.Now().Add(this.write_timeout))
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return this.TCPConn.Write(b)
// }

// type client struct {
// 	HTTPClient      *http.Client
// 	connect_timeout time.Duration
// 	read_timeout    time.Duration
// 	write_timeout   time.Duration
// }

// func newClient(connect_timeout, read_timeout, write_timeout time.Duration) *client {
// 	dialer := func(netw, addr string) (net.Conn, error) {
// 		conn, err := net.DialTimeout(netw, addr, connect_timeout)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &rwTimeoutConn{
// 			TCPConn:       conn.(*net.TCPConn),
// 			read_timeout:  read_timeout,
// 			write_timeout: write_timeout,
// 		}, nil
// 	}
// 	cli := &client{
// 		HTTPClient: &http.Client{
// 			Transport: &http.Transport{
// 				Proxy:                 http.ProxyFromEnvironment,
// 				Dial:                  dialer,
// 				MaxIdleConns:          100,
// 				IdleConnTimeout:       90 * time.Second,
// 				TLSHandshakeTimeout:   10 * time.Second,
// 				ExpectContinueTimeout: 1 * time.Second,
// 				ForceAttemptHTTP2:     true,
// 				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
// 			},
// 		},
// 	}

// 	return cli
// }

// func make_request(method, urlStr string, maxRetries int, opts ...func(*http.Request)) (*http.Response, error) {
// 	var proxy_url *url.URL
// 	var err error
// 	if confy.Get[string]("ssrf.proxy_all_url") != "" {
// 		proxy_url, err = url.Parse(confy.Get[string]("ssrf.proxy_all_url"))
// 		if err != nil {
// 			mlog.Errorf("Invalid ssrf.proxy_all_url(%s): %v", confy.Get[string]("ssrf.proxy_all_url"), err)
// 		}
// 	} else if confy.Get[string]("ssrf.proxy_http_url") != "" && confy.Get[string]("ssrf.proxy_https_url") != "" {
// 		http_proxy_url, err := url.Parse(confy.Get[string]("ssrf.proxy_http_url"))
// 		if err != nil {
// 			mlog.Errorf("Invalid ssrf.proxy_http_url(%s): %v", confy.Get[string]("ssrf.proxy_http_url"), err)
// 		}
// 		https_proxy_url, err := url.Parse(confy.Get[string]("ssrf.proxy_https_url"))
// 		if err != nil {
// 			mlog.Errorf("Invalid ssrf.proxy_https_url(%s): %v", confy.Get[string]("ssrf.proxy_https_url"), err)
// 		}
// 		proxyURL = httpProxyURL
// 		req.Header.Set("HTTP-Proxy", httpProxyURL.String())
// 		req.Header.Set("HTTPS-Proxy", httpsProxyURL.String())
// 	}

// 	if proxyURL != nil {
// 		client.HTTPClient.Transport = &http.Transport{
// 			Proxy: http.ProxyURL(proxyURL),
// 		}
// 	}

// 	client := newClient(
// 		time.Duration(confy.GetWithDefault[int]("ssrf.default_connect_timeout", 5))*time.Second,
// 		time.Duration(confy.GetWithDefault[int]("ssrf.default_read_timeout", 5))*time.Second,
// 		time.Duration(confy.GetWithDefault[int]("ssrf.default_write_timeout", 5))*time.Second,
// 	)

// 	req, err := http.NewRequest(method, urlStr, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, opt := range opts {
// 		opt(req)
// 	}

// 	if req.Header.Get("Allow-Redirects") != "" {
// 		if req.Header.Get("Follow-Redirects") == "" {
// 			req.Header.Set("Follow-Redirects", req.Header.Get("Allow-Redirects"))
// 		}
// 		req.Header.Del("Allow-Redirects")
// 	}

// 	req.Header.Set("Connect-TimeOut", fmt.Sprintf("%f", confy.GetWithDefault[int]("ssrf.default_connect_timeout", 5)))
// 	req.Header.Set("Read-TimeOut", fmt.Sprintf("%f", confy.GetWithDefault[int]("ssrf.default_read_timeout", 5)))
// 	req.Header.Set("Write-TimeOut", fmt.Sprintf("%f", confy.GetWithDefault[int]("ssrf.default_write_timeout", 5)))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !contains(resp.StatusCode, STATUS_FORCELIST) {
// 		return resp, nil
// 	}

// 	log.Printf("Received status code %d for URL %s which is in the force list", resp.StatusCode, urlStr)
// 	return resp, nil
// }
