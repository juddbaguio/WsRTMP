package WsRTMP

func (c *WsRTMPClient) StreamToPipe(stream *[]byte) {
	c.InputStream <- stream
}
