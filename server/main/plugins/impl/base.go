package impl

type BasePluginClient struct{

}

func(cli *BasePluginClient) _request(
        method  string,
        path  string,
        headers  map[string]any,
        data: bytes | dict | str | None = None,
        params  map[string]any,
        files  map[string]any,
        stream bool,
    ) -> requests.Response:
        """
        Make a request to the plugin daemon inner API.
        """
        url = plugin_daemon_inner_api_baseurl / path
        headers = headers or {}
        headers["X-Api-Key"] = dify_config.PLUGIN_DAEMON_KEY
        headers["Accept-Encoding"] = "gzip, deflate, br"

        if headers.get("Content-Type") == "application/json" and isinstance(data, dict):
            data = json.dumps(data)

        try:
            response = requests.request(
                method=method, url=str(url), headers=headers, data=data, params=params, stream=stream, files=files
            )
        except requests.exceptions.ConnectionError:
            logger.exception("Request to Plugin Daemon Service failed")
            raise PluginDaemonInnerError(code=-500, message="Request to Plugin Daemon Service failed")

        return response