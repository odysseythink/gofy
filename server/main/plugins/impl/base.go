package impl

import "github.com/spf13/viper"

type BasePluginClient struct{

}

func _request[T []byte|map[string]any|string](
        method  string,
        path  string,
        headers  map[string]any,
        data T,
        params  map[string]any,
        files  map[string]any,
        stream bool,
    ) *http.Response{
        url := viper.GetStringWithDefault("plugin.daemon_url", "http://localhost:5002") +"/"+ path

			reqBody, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		panic(httpexceptions.NewPluginDaemonInnerError("Failed to create request"))
	}
	req.Header = headers
	req.Header["X-Api-Key"] = viper.GetStringWithDefault("plugin.plugin-api-key", "")
    req.Header["Accept-Encoding"] = "gzip, deflate, br"
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		mlog.Errorf("dump request failed:%v", err)
	} else {
		fmt.Printf("REQUEST:\n%s", string(reqDump))
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(httpexceptions.NewPluginDaemonInnerError("Request to Plugin Daemon Service failed"))
	}

        headers = headers or {}


        if headers.get("Content-Type") == "application/json" and isinstance(data, dict):
            data = json.dumps(data)
}
        try:
            response = requests.request(
                method=method, url=str(url), headers=headers, data=data, params=params, stream=stream, files=files
            )
        except requests.exceptions.ConnectionError:
            logger.exception("Request to Plugin Daemon Service failed")
            raise PluginDaemonInnerError(code=-500, message="Request to Plugin Daemon Service failed")

        return response
		}