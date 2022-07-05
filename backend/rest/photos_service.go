package main

import (
	protos "backend/grpc/proto/api/photos"
	photoshandlers "backend/rest/handlers/google/photos"
	gAuth "backend/rest/middleware/google/auth/OAuth"
	callingCatchables "backend/rest/middleware/google/callingCatchables/photos"

	"log"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PhotoService struct {
	logger     *log.Logger
	gpsc       protos.GooglePhotoServiceClient
	pHandler   *photoshandlers.PhotoHandler
	gAuthMware *gAuth.AuthMiddleware
	gPhotos    *callingCatchables.PhotosRpcCaller
}

func GetPhotoService() *PhotoService {
	ps := PhotoService{}
	return &ps
}

func (ps *PhotoService) InitServiceAndReturnCloseConnectionFunc() func() {
	// Main logger
	ps.logger = log.New(os.Stdout, "rest-server-photos", log.LstdFlags)

	/////// Initialize GRPC connections
	photoConn, err := grpc.Dial("grpc_backend:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	//defer photoConn.Close()
	ps.gpsc = protos.NewGooglePhotoServiceClient(photoConn)

	/////// Initialize middleware and handlers here ///////
	ps.pHandler = photoshandlers.NewPhotoHandler(ps.logger)
	ps.gAuthMware = gAuth.GetAuthMiddleware()
	ps.gPhotos = callingCatchables.NewPhotosRpcCaller(ps.logger, &ps.gpsc)

	return func() { photoConn.Close() }
}

func (ps *PhotoService) RegisterGetRoutes(getRouter *mux.Router) {

	//getRouter.HandleFunc("/", )
	getRouter.HandleFunc("/authenticate", ps.gAuthMware.Authenticate)
	getRouter.HandleFunc("/oauth-callback", ps.gAuthMware.RedirectCallback)

	//route for listing albums - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/albumsList", ps.gAuthMware.IsAuthorized(ps.gPhotos.CatchableListAlbums(ps.pHandler.ListAlbums)))

	//route for listing photos in an album - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/album/{albumId:[-_0-9A-Za-z]+}", ps.gAuthMware.IsAuthorized(ps.gPhotos.CatchablePhotosFromAlbum(ps.pHandler.ListPhotosFromAlbum)))

	getRouter.HandleFunc("/oh-no", ps.pHandler.OhNo)
}
