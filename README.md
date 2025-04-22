## mcp-log-proxy

`mcp-log-proxy` can be used to see the messages to and from a MCP client and a MCP server using a Web interface.
Currently, it only supports the STDIO interface.

## install

    go install github.com/emicklei/mcp-log-proxy@latest

### usage

`mcp-log-proxy` requires one argument `-command` that contains the full command line for starting the MCP server.

For example, to proxy traffic to the `melrose-mcp` server, the full command is:

    mcp-log-proxy -command melrose-mcp

This example assumes that both tools are available on your execution PATH.

Optionally, you can override the log file location of the proxy that captures errors in the proxy itself.

    mcp-log-proxy -command melrose-mcp -log /your/logs/mcp-log-proxy.log

When the proxy is started, messages can be viewed on `http:/localhost:5656` (use `-port` to override).

&copy; 2025, https://ernestmicklei.com. MIT License.