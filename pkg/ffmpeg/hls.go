package ffmpeg

import "context"

type HLSCfg struct {
	SegmentSec    int
	ListSize      int
	SingleVariant bool
}

func TranscodeToHLS(ctx context.Context, inURI, outDir string, cfg HLSCfg) (m3u8Key string, err error) {
	// TODO: invoke ffmpeg to produce m3u8+ts; upload to S3 and return key
	return "preview/demo/index.m3u8", nil
}
