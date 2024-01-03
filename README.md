# Proxy

This is a simple proxy server that can be used to forward requests to a different server. It is inspired by NGINX and
Apache, but is much simpler and easier to use.

## Usage

To use the proxy, you must first create a configuration file. The configuration file is a YAML file that contains
information about the proxy server and the servers that it will forward requests to. The configuration file must be
named `config.yaml` and must be in a subdirectory named `data`. The `data` directory must be in the same directory as
working directory of the proxy server. The configuration file must have the following format:

```yaml
resources:
  - endpoint: /server1
    method: POST
    redirect: 'http://example.com/path'
  - endpoint: /server2
    method: GET
    redirect: 'http://example.com/path'
    auth:
      username: username
      password: password
```

The `resources` key is a list of resources that the proxy server will forward requests to. Each resource must have an
`endpoint`, `method`, and `redirect` key. The `endpoint` key is the endpoint that the proxy server will listen for
requests on. The `method` key is the HTTP method that the proxy server will listen for requests on. The `redirect` key
is the URL that the proxy server will forward requests to. The `auth` key is optional and is used to specify the
username and password that the proxy server will use to authenticate with the server that it is forwarding requests to.
