package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/issue9/identicon"
	toolpb "github.com/sbasestarter/proto-repo/gen/protorepo-tool-go"
	"github.com/sbasestarter/toolsb/internal/config"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

type generatePixelAvatarHTTPRequest struct {
	Size int
	Fore color.RGBA
	Back color.RGBA
	Text string
}

func (s *Server) GetHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/GeneratePixelAvatar", func(writer http.ResponseWriter, request *http.Request) {
		var req generatePixelAvatarHTTPRequest
		var err error
		var d []byte
		if request.Method == "POST" {
			d, err = ioutil.ReadAll(request.Body)
			if err == nil {
				err = json.Unmarshal(d, &req)
			}
			if err == nil {
				if req.Size <= 0 {
					err = errors.New("no size")
				}
			}
		} else {
			req.Size = 128
			req.Fore.R = 0xff
			req.Fore.A = 0xff
		}

		if err == nil {
			d = s.generatePixelAvatar(context.Background(), req.Size, req.Fore, req.Back, req.Text)
			writer.Header().Set("Content-Disposition", "attachment; filename=avatar.png")
			writer.Header().Set("Content-Type", "image/png")
			writer.Header().Set("Content-Length", strconv.FormatInt(int64(len(d)), 10))
			_, _ = writer.Write(d)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
		}
	})
	return r
}

func (s *Server) GeneratePixelAvatar(ctx context.Context, req *toolpb.GeneratePixelAvatarRequest) (
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

	return &toolpb.GeneratePixelAvatarResponse{
		Status: &toolpb.ServerStatus{
			Status: toolpb.ToolStatus_TS_SUCCESS,
		},
		Data: s.generatePixelAvatar(ctx, int(req.Size), fore, back, req.Text),
	}, nil
}

func (s *Server) generatePixelAvatar(_ context.Context, size int, fore, back color.Color, text string) []byte {
	img, _ := identicon.Make(size, back, fore, []byte(text))
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
