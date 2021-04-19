package server

import (
	"bytes"
	"context"
	"image/color"
	"image/png"

	"github.com/issue9/identicon"
	toolpb "github.com/sbasestarter/proto-repo/gen/protorepo-tool-go"
	"github.com/sbasestarter/toolsb/internal/config"
)

type serverImpl struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) toolpb.UserServiceServer {
	return &serverImpl{
		cfg: cfg,
	}
}

func (s *serverImpl) GeneratePixelAvatar(_ context.Context, req *toolpb.GeneratePixelAvatarRequest) (
	*toolpb.GeneratePixelAvatarResponse, error) {
	back := color.NRGBA{}
	fore := color.NRGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0xff,
	}
	if req.Fore != nil {
		fore.R = uint8(req.Fore.R)
		fore.G = uint8(req.Fore.G)
		fore.B = uint8(req.Fore.B)
		fore.A = uint8(req.Fore.A)
	}
	if req.Back != nil {
		back.R = uint8(req.Back.R)
		back.G = uint8(req.Back.G)
		back.B = uint8(req.Back.B)
		back.A = uint8(req.Back.A)
	}

	img, _ := identicon.Make(int(req.Size), back, fore, []byte("192.168.1.1"))
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return &toolpb.GeneratePixelAvatarResponse{
		Status: &toolpb.ServerStatus{
			Status: toolpb.ToolStatus_TS_SUCCESS,
		},
		Data: buf.Bytes(),
	}, nil
}
