package WsRTMP

import "os/exec"

type WsRTMPClient struct {
	ffmpeg      *exec.Cmd
	InputStream chan *[]byte
}

type RTMPConf struct {
}

var ffmpegArgs = []string{
	"-i", "-",
	"-vcodec", "copy",
	"-f", "flv",
}

func New(rtmpDest string) *WsRTMPClient {
	streamPipe := make(chan *[]byte)
	cmd := exec.Command("ffmpeg",
		append(ffmpegArgs, rtmpDest)...)

	return &WsRTMPClient{
		ffmpeg:      cmd,
		InputStream: streamPipe,
	}
}
