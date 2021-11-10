package WsRTMP

import (
	"fmt"
	"os"
)

func (c *WsRTMPClient) StartBroadcast() error {
	c.ffmpeg.Stderr = os.Stderr

	ffmpegInput, err := c.ffmpeg.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get input pipe for FFmpeg process: %w", err)
	}

	err = c.ffmpeg.Start()
	if err != nil {
		ffmpegInput.Close()
		return fmt.Errorf("failed to start FFmpeg process: %w", err)
	}

	go func() {
		defer ffmpegInput.Close()
		defer c.ffmpeg.Wait()

		for {
			stream := <-c.InputStream

			_, err := ffmpegInput.Write(*stream)
			if err != nil {
				break
			}
		}
	}()

	return nil
}
