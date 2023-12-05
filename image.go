package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

const (
	uploadPictureEndpoint = "/pictures"
)

// PictureEncoding specifies the encoding and Content-Type
// for the image upload.
type PictureEncoding string

const (
	PictureEncodingPNG  = "png"
	PictureEncodingJPEG = "jpeg"
)

// PictureURL contains URLS to an uploaded picture as well as the
// various thumbnails provided by GroupMe.
type PictureURL struct {
	Base    string
	Preview string
	Large   string
	Avatar  string
}

// UploadPicture posts an image to the GroupMe image service. Accepts either PNG or JPEG.
// Returns URLs to the uploaded image to be used in messages or avatars.
func (c *Client) UploadPicture(ctx context.Context, img image.Image, encoding PictureEncoding) (PictureURL, error) {
	var imgBytes bytes.Buffer
	var err error
	switch encoding {
	case PictureEncodingPNG:
		err = png.Encode(&imgBytes, img)
	case PictureEncodingJPEG:
		err = jpeg.Encode(&imgBytes, img, nil)
	default:
		err = fmt.Errorf("unsupported encoding: %s", encoding)
	}

	if err != nil {
		return PictureURL{}, fmt.Errorf("failed to encode image: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.imageEndpointBase+uploadPictureEndpoint, &imgBytes)
	if err != nil {
		return PictureURL{}, err
	}
	httpReq.Header.Add("Content-Type", "image/"+string(encoding))

	URL := httpReq.URL
	query := URL.Query()
	query.Set("token", c.authorizationToken)
	URL.RawQuery = query.Encode()

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return PictureURL{}, err
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)

	var resp struct {
		Payload struct {
			URL        string `json:"url"`
			PictureURL string `json:"picture_url"`
		} `json:"payload"`
	}
	if err := decoder.Decode(&resp); err != nil {
		return PictureURL{}, err
	}

	return PictureURL{
		Base:    resp.Payload.URL,
		Preview: resp.Payload.URL + ".preview",
		Large:   resp.Payload.URL + ".large",
		Avatar:  resp.Payload.URL + ".avatar",
	}, nil
}
