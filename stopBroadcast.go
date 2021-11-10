package WsRTMP

import "fmt"

func (c *WsRTMPClient) StopBroadcast() error {
	err := c.ffmpeg.Process.Kill()

	if err != nil {
		return fmt.Errorf("failed to kill FFmpeg process: %w", err)
	}

	return nil
}
