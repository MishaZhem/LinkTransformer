package main

import "LinkTransformer/internal/ports/httpgin"

func main() {
	httpServer := httpgin.NewHTTPServer(":18080")
	httpServer.Listen()
}
