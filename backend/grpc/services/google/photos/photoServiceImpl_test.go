package photos

import (
	"backend/grpc/proto/api/POGO"
	photos "backend/grpc/proto/api/photos"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
)

func generateMockAlbumInfoPOGO(id string, title string, productUrl string, mediaItemsCount string, coverPhotoBaseUrl string, coverPhotoMediaItemId string) POGO.AlbumInfo {
	return POGO.AlbumInfo{
		Id:                    id,
		Title:                 title,
		ProductUrl:            productUrl,
		MediaItemsCount:       mediaItemsCount,
		CoverPhotoBaseUrl:     coverPhotoBaseUrl,
		CoverPhotoMediaItemId: coverPhotoMediaItemId}
}

func generateMockAlbumInfoProto(id string, title string, productUrl string, mediaItemsCount int32, coverPhotoBaseUrl string, coverPhotoMediaItemId string) *photos.AlbumInfo {
	return &photos.AlbumInfo{
		Id:                    id,
		Title:                 title,
		ProductUrl:            productUrl,
		MediaItemsCount:       mediaItemsCount,
		CoverPhotoBaseUrl:     coverPhotoBaseUrl,
		CoverPhotoMediaItemId: coverPhotoMediaItemId,
	}
}

func generateMockAlbumsInfoPOGO(ai ...POGO.AlbumInfo) POGO.AlbumsInfoPOGO {
	var infoList []POGO.AlbumInfo

	for _, info := range ai {
		infoList = append(infoList, info)
	}
	return POGO.AlbumsInfoPOGO{
		AlbumsInfo: infoList,
	}
}

func generateMockAlbumsInfoProto(ai ...*photos.AlbumInfo) photos.AlbumsInfo {
	var infoList []*photos.AlbumInfo

	for _, info := range ai {
		infoList = append(infoList, info)
	}

	return photos.AlbumsInfo{
		AlbumsInfo: infoList,
	}
}

func generateMockAlbumInfoString(id string, title string, productUrl string, mediaItemsCount string, coverPhotoBaseUrl string, coverPhotoMeidaItemId string) string {
	return fmt.Sprintf(`{
			"id": "%s",
			"title": "%s",
			"productUrl": "%s",
			"mediaItemsCount": "%s",
			"coverPhotoBaseUrl": "%s",
			"coverPhotoMediaItemId": "%s"}`, id, title, productUrl, mediaItemsCount, coverPhotoBaseUrl, coverPhotoMeidaItemId)
}

func generateMockStringResponse(infos ...string) io.ReadCloser {
	result := ""
	for _, info := range infos {
		result += info + ","
	}

	albums := fmt.Sprintf(`{"albums": [%s] }`, result[:len(result)-1])

	return ioutil.NopCloser(strings.NewReader(albums))
}

func TestAlbumListDecoder(t *testing.T) {

	// Single album test
	expected := generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"))
	info := generateMockAlbumInfoString("123", "abc", "xyz", "456", "abc.com", "xyz.com")
	mockResponse := generateMockStringResponse(info)

	got := albumListDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

	// Multiple albums tests
	expected = generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"),
		generateMockAlbumInfoPOGO("456", "def", "422", "12", "grrr.com", "sparta.com"))

	info2 := generateMockAlbumInfoString("456", "def", "422", "12", "grrr.com", "sparta.com")
	mockResponse = generateMockStringResponse(info, info2)

	got = albumListDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

	// Multiple albums tests with empty
	expected = generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"),
		generateMockAlbumInfoPOGO("456", "def", "422", "12", "grrr.com", "sparta.com"),
		generateMockAlbumInfoPOGO("", "", "", "", "", ""))

	info3 := generateMockAlbumInfoString("", "", "", "", "", "")
	mockResponse = generateMockStringResponse(info, info2, info3)

	got = albumListDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

}

func TestAlbumsPogo2Proto(t *testing.T) {

	// Single Album
	expected := generateMockAlbumsInfoProto(generateMockAlbumInfoProto("123", "abc", "xyz", 456, "abc.com", "xyz.com"))

	mockPOGO := generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"))
	got := albumsPogo2Proto(&mockPOGO)

	if !proto.Equal(&expected, got) {
		t.Fatalf("Expected: \n%v, Got: \n%v", &expected, got)
	}

	// Multiple Albums
	expected = generateMockAlbumsInfoProto(generateMockAlbumInfoProto("123", "abc", "xyz", 456, "abc.com", "xyz.com"),
		generateMockAlbumInfoProto("123", "abc", "xyz", 456, "abc.com", "xyz.com"))

	mockPOGO = generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"),
		generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"))

	got = albumsPogo2Proto(&mockPOGO)

	if !proto.Equal(&expected, got) {
		t.Fatalf("Expected: \n%v, Got: \n%v", &expected, got)
	}

	// Multiple albums tests with empty
	expected = generateMockAlbumsInfoProto(generateMockAlbumInfoProto("123", "abc", "xyz", 456, "abc.com", "xyz.com"),
		generateMockAlbumInfoProto("123", "abc", "xyz", 456, "abc.com", "xyz.com"),
		generateMockAlbumInfoProto("", "", "", 0, "", ""))

	mockPOGO = generateMockAlbumsInfoPOGO(generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"),
		generateMockAlbumInfoPOGO("123", "abc", "xyz", "456", "abc.com", "xyz.com"),
		generateMockAlbumInfoPOGO("", "", "", "0", "", ""))

	got = albumsPogo2Proto(&mockPOGO)

	if !proto.Equal(&expected, got) {
		t.Fatalf("Expected: \n%v, Got: \n%v", &expected, got)
	}
}

func generateMockPhotoInfoPOGO(id string, productUrl string, baseUrl string, mimeType string, filename string) POGO.PhotosInfo {
	return POGO.PhotosInfo{Id: id,
		ProductUrl: productUrl,
		BaseUrl:    baseUrl,
		MimeType:   mimeType,
		Filename:   filename}
}

func generateMockPhotosInfoPOGO(pi ...POGO.PhotosInfo) POGO.PhotosInfoPOGO {
	var infoList []POGO.PhotosInfo

	for _, info := range pi {
		infoList = append(infoList, info)
	}
	return POGO.PhotosInfoPOGO{
		MediaItems: infoList,
	}
}

func generateMockPhotosInfoString(id string, productUrl string, baseUrl string, mimeType string, fileName string) string {
	return fmt.Sprintf(`{
			"id": "%s",
			"productUrl": "%s",
			"baseUrl": "%s",
			"mimeType": "%s",
			"fileName": "%s"}`, id, productUrl, baseUrl, mimeType, fileName)
}

func generateFromPhotosMockStringResponse(infos ...string) io.ReadCloser {
	result := ""
	for _, info := range infos {
		result += info + ","
	}

	photos := fmt.Sprintf(`{"mediaItems": [%s] }`, result[:len(result)-1])

	return ioutil.NopCloser(strings.NewReader(photos))
}

func TestPhotosFromAlbumDecoder(t *testing.T) {

	// Single photo test
	expected := generateMockPhotosInfoPOGO(generateMockPhotoInfoPOGO("123", "abc", "xyz", "456", "abc.com"))

	info := generateMockPhotosInfoString("123", "abc", "xyz", "456", "abc.com")
	mockResponse := generateFromPhotosMockStringResponse(info)

	got := photosFromAlbumDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

	// Multiple photo test
	expected = generateMockPhotosInfoPOGO(generateMockPhotoInfoPOGO("123", "abc", "xyz", "456", "abc.com"),
		generateMockPhotoInfoPOGO("123", "abc", "xyz", "456", "abc.com"))

	info2 := generateMockPhotosInfoString("123", "abc", "xyz", "456", "abc.com")
	mockResponse = generateFromPhotosMockStringResponse(info, info2)

	got = photosFromAlbumDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

	// Multiple photo test with empty
	expected = generateMockPhotosInfoPOGO(generateMockPhotoInfoPOGO("123", "abc", "xyz", "456", "abc.com"),
		generateMockPhotoInfoPOGO("123", "abc", "xyz", "456", "abc.com"),
		generateMockPhotoInfoPOGO("", "", "", "", ""))

	info3 := generateMockPhotosInfoString("", "", "", "", "")
	mockResponse = generateFromPhotosMockStringResponse(info, info2, info3)

	got = photosFromAlbumDecoder(mockResponse)

	if !expected.Equals(got) {
		t.Fatalf("Expected: %v, Got: %v", expected, got)
	}

}
