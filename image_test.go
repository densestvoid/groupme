package groupme

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"image/jpeg"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

//go:embed image.jpeg
var imgBytes []byte

type PictureAPISuite struct{ APISuite }

func (s *PictureAPISuite) SetupSuite() {
	s.handler = picturesTestRouter()
	s.setupSuite()
}

func (s *PictureAPISuite) TestUsersMe() {
	img, err := jpeg.Decode(bytes.NewBuffer(imgBytes))
	s.Require().NoError(err)

	picture, err := s.client.UploadPicture(context.Background(), img, PictureEncodingJPEG)
	s.Require().NoError(err)
	s.Assert().NotZero(picture)
}

func TestPicturesAPISuite(t *testing.T) {
	suite.Run(t, new(PictureAPISuite))
}

// nolint // not duplicate code
func picturesTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Me
	router.Path("/pictures").
		Methods("POST").
		Name("UploadPicture").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"payload": {
					"url": "https://test.com/100x100.jpeg.123456789",
					"picture_url": "https://test.com/100x100.jpeg.123456789",
				},
			}`)
		})

	/*// Return test router //*/
	return router
}
