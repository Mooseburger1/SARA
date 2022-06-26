package photos

import (
	"backend/grpc/proto/api/POGO"
	photos "backend/grpc/proto/api/photos"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

// MOCK 1
var mockAlbumInfo1 = POGO.AlbumInfo{
	Id:                    "123",
	Title:                 "abc",
	ProductUrl:            "xyz",
	MediaItemsCount:       "456",
	CoverPhotoBaseUrl:     "abc.com",
	CoverPhotoMediaItemId: "xyz.gov"}
var count, _ = strconv.Atoi(mockAlbumInfo1.MediaItemsCount)
var mockProtoAlbumInfo1 = &photos.AlbumInfo{
	Id:                    mockAlbumInfo1.Id,
	Title:                 mockAlbumInfo1.Title,
	ProductUrl:            mockAlbumInfo1.ProductUrl,
	MediaItemsCount:       int32(count),
	CoverPhotoBaseUrl:     mockAlbumInfo1.CoverPhotoBaseUrl,
	CoverPhotoMediaItemId: mockAlbumInfo1.CoverPhotoMediaItemId,
}
var mockPOGO1 = POGO.AlbumsInfoPOGO{
	AlbumsInfo: []POGO.AlbumInfo{mockAlbumInfo1, mockAlbumInfo2, mockAlbumInfo3}}
var mockPROTO1 = photos.AlbumsInfo{
	AlbumsInfo: []*photos.AlbumInfo{mockProtoAlbumInfo1, mockProtoAlbumInfo2, mockProtoAlbumInfo3},
}

// MOCK 2
var mockAlbumInfo2 = POGO.AlbumInfo{
	Id:                "123",
	Title:             "abc",
	ProductUrl:        "xyz",
	CoverPhotoBaseUrl: "abc.com"}
var mockProtoAlbumInfo2 = &photos.AlbumInfo{
	Id:                mockAlbumInfo2.Id,
	Title:             mockAlbumInfo2.Title,
	ProductUrl:        mockAlbumInfo2.ProductUrl,
	CoverPhotoBaseUrl: mockAlbumInfo2.CoverPhotoBaseUrl,
}
var mockPOGO2 = POGO.AlbumsInfoPOGO{
	AlbumsInfo: []POGO.AlbumInfo{mockAlbumInfo2}}

var mockPROTO2 = photos.AlbumsInfo{
	AlbumsInfo: []*photos.AlbumInfo{mockProtoAlbumInfo2},
}

// MOCK 3
var mockAlbumInfo3 = POGO.AlbumInfo{}
var mockProtoAlbumInfo3 = &photos.AlbumInfo{}
var mockPOGO3 = POGO.AlbumsInfoPOGO{
	AlbumsInfo: []POGO.AlbumInfo{mockAlbumInfo3}}
var mockPROTO3 = photos.AlbumsInfo{
	AlbumsInfo: []*photos.AlbumInfo{mockProtoAlbumInfo3},
}

// MOCK 4
var mockPOGO4 = POGO.AlbumsInfoPOGO{
	AlbumsInfo: []POGO.AlbumInfo{}}
var mockPROTO4 = photos.AlbumsInfo{}

// MOCK 1
var mockListAlbumResponseString = `{
	"albums": [
	  {
		"id": "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"title": "Sims Family",
		"productUrl": "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"mediaItemsCount": "16",
		"coverPhotoBaseUrl": "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		"coverPhotoMediaItemId": "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw"
	  }
	]
  }`

var mockListAlbumResponse1 = ioutil.NopCloser(strings.NewReader(mockListAlbumResponseString))

// MOCK 2
var mockListAlbumResponseString2 = `{
	"albums": [
	  {
		"id": "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"title": "Sims Family",
		"productUrl": "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"mediaItemsCount": "16",
		"coverPhotoBaseUrl": "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		"coverPhotoMediaItemId": "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw"
	  },
	  {
		"id": "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"title": "Sims Family",
		"productUrl": "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		"mediaItemsCount": "16",
		"coverPhotoBaseUrl": "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		"coverPhotoMediaItemId": "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw"
	  }
	]
  }`

var mockListAlbumResponse2 = ioutil.NopCloser(strings.NewReader(mockListAlbumResponseString2))

// MOCK 3
var mockListAlbumResponseString3 = `{
	"albums": [
	  {
		"id": "",
		"title": "",
		"productUrl": "",
		"mediaItemsCount": "",
		"coverPhotoBaseUrl": "",
		"coverPhotoMediaItemId": ""
	  }
	]
  }`

var mockListAlbumResponse3 = ioutil.NopCloser(strings.NewReader(mockListAlbumResponseString3))

func TestAlbumListDecoder(t *testing.T) {
	mockAI := POGO.AlbumInfo{
		Id:                    "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		Title:                 "Sims Family",
		ProductUrl:            "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		MediaItemsCount:       "16",
		CoverPhotoBaseUrl:     "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		CoverPhotoMediaItemId: "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw",
	}

	expectedAI := POGO.AlbumsInfoPOGO{
		AlbumsInfo: []POGO.AlbumInfo{mockAI},
	}

	AI := albumListDecoder(mockListAlbumResponse1)
	if !expectedAI.Equals(AI) {
		t.Fatalf("Expected: %v\nGot: %v", expectedAI, AI)
	}

	mockAI21 := POGO.AlbumInfo{
		Id:                    "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		Title:                 "Sims Family",
		ProductUrl:            "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		MediaItemsCount:       "16",
		CoverPhotoBaseUrl:     "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		CoverPhotoMediaItemId: "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw",
	}

	mockAI22 := POGO.AlbumInfo{
		Id:                    "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		Title:                 "Sims Family",
		ProductUrl:            "https://photos.google.com/lr/album/AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd",
		MediaItemsCount:       "16",
		CoverPhotoBaseUrl:     "https://lh3.googleusercontent.com/lr/AFBm1_aPsjiK0dNoBX5fLs4CPz5KMbfPTVoDnGmX4zzo2slUx5TLjCxSfo1rRcYTGd16mg4osivBvpy5yaimbM9XTAHJgkuSktTyjjOI9BbTh1_zR6ZXB1_gbe0dbn-ib6XXVEuZiM14zymF3ZJSHJxztTJm43aP-byCZarNpz6g6e4dC7IV7Xjny7i5YQ6uDbcZJqT9_QjeMY2VBoxMsZWXQtchn97RgZJwZz7Puv9INNrEdHLOlwYYGFRn9KovWAbocuHQxoAsWOVCQzjFo3Eas0ZUPvNUX7UV7yj6GHA6f6tyRgTFGpZ9a3ZzGmAy2AwUG9c8B61AjW46I-Mmc-OrHXc6EwwCZInx4ZqVvaUvGRpcVTvwtmrP7dfYnptKOmM2E01PQsUsyQkYi7NLCL45fXfNT6UhJtZktSxL4YI8mJZVOsP2VpJ0Byq3_9p-QTmdhnKss6FlItOYedS3RwFZsFSAxf0Ahz5TJrv7Gw_5X3wNsIp0Wp2wzAPzvIbW9iFAERABn5izo1bglpSKXTuCviDVAwsF3DXv9R-6nmHFdAq6oz770ynl38RJjDcDu5IHi0ab738wWy7MWw1kwRVqUHF7POlNY3ARfVZdqFLYv-bSSfjkKySxYMns0qAtMbpyWMOv_hUxS0-naTJpS8kwKg5Q7_tBVsWoTRxi5c7qpUEVSToRMlyTY0volbpKCUUOjsJh0TIn6lY-PWg7L8HR8i-4nJ3SWsJLsBhT2uiPw1rP3UaxY5jtqGK16hQlu08yurbwF4rKSqabd9QXVMxMYTq5Iza6umSp_zzHejWrXsq4RFcyo6sv52A4d-3zqPVAzgreCw",
		CoverPhotoMediaItemId: "AAmUB7WNP3mdI5CphiXKMgse0qdLuz6Xvl-qlfwr18IrL9hm326YAlBIHefxOus7g7HMP0heSYJ1tlpDPRJNRelZddl7fbRnaw",
	}

	expectedAI2 := POGO.AlbumsInfoPOGO{
		AlbumsInfo: []POGO.AlbumInfo{mockAI21, mockAI22},
	}

	AI2 := albumListDecoder(mockListAlbumResponse2)
	if !expectedAI2.Equals(AI2) {
		t.Fatalf("Expected: %v\nGot: %v", expectedAI2, AI2)
	}

	mockAI3 := POGO.AlbumInfo{
		Id:                    "",
		Title:                 "",
		ProductUrl:            "",
		MediaItemsCount:       "",
		CoverPhotoBaseUrl:     "",
		CoverPhotoMediaItemId: "",
	}

	expectedAI3 := POGO.AlbumsInfoPOGO{
		AlbumsInfo: []POGO.AlbumInfo{mockAI3},
	}

	AI3 := albumListDecoder(mockListAlbumResponse3)
	if !expectedAI3.Equals(AI3) {
		t.Fatalf("Expected: %v\nGot: %v", expectedAI3, AI3)
	}
}

func TestAlbumsPogo2Proto(t *testing.T) {

	albuminfo1 := albumsPogo2Proto(&mockPOGO1)

	for i, ai := range albuminfo1.AlbumsInfo {
		expected := mockPROTO1.AlbumsInfo[i]
		if ai.Id != expected.Id {
			t.Fatalf("Expected: %v\nGot: %v", expected.Id, ai.Id)
		} else if ai.Title != expected.Title {
			t.Fatalf("Expected: %v\nGot: %v", expected.Title, ai.Title)
		} else if ai.ProductUrl != expected.ProductUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.ProductUrl, ai.ProductUrl)
		} else if ai.CoverPhotoBaseUrl != expected.CoverPhotoBaseUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoBaseUrl, ai.CoverPhotoBaseUrl)
		} else if ai.CoverPhotoMediaItemId != expected.CoverPhotoMediaItemId {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoMediaItemId, ai.CoverPhotoMediaItemId)
		}
	}

	albuminfo2 := albumsPogo2Proto(&mockPOGO2)

	for i, ai := range albuminfo2.AlbumsInfo {
		expected := mockPROTO2.AlbumsInfo[i]
		if ai.Id != expected.Id {
			t.Fatalf("Expected: %v\nGot: %v", expected.Id, ai.Id)
		} else if ai.Title != expected.Title {
			t.Fatalf("Expected: %v\nGot: %v", expected.Title, ai.Title)
		} else if ai.ProductUrl != expected.ProductUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.ProductUrl, ai.ProductUrl)
		} else if ai.CoverPhotoBaseUrl != expected.CoverPhotoBaseUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoBaseUrl, ai.CoverPhotoBaseUrl)
		} else if ai.CoverPhotoMediaItemId != expected.CoverPhotoMediaItemId {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoMediaItemId, ai.CoverPhotoMediaItemId)
		}
	}

	albuminfo3 := albumsPogo2Proto(&mockPOGO3)

	for i, ai := range albuminfo3.AlbumsInfo {
		expected := mockPROTO3.AlbumsInfo[i]
		if ai.Id != expected.Id {
			t.Fatalf("Expected: %v\nGot: %v", expected.Id, ai.Id)
		} else if ai.Title != expected.Title {
			t.Fatalf("Expected: %v\nGot: %v", expected.Title, ai.Title)
		} else if ai.ProductUrl != expected.ProductUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.ProductUrl, ai.ProductUrl)
		} else if ai.CoverPhotoBaseUrl != expected.CoverPhotoBaseUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoBaseUrl, ai.CoverPhotoBaseUrl)
		} else if ai.CoverPhotoMediaItemId != expected.CoverPhotoMediaItemId {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoMediaItemId, ai.CoverPhotoMediaItemId)
		}
	}

	albuminfo4 := albumsPogo2Proto(&mockPOGO4)

	for i, ai := range albuminfo4.AlbumsInfo {
		expected := mockPROTO4.AlbumsInfo[i]
		if ai.Id != expected.Id {
			t.Fatalf("Expected: %v\nGot: %v", expected.Id, ai.Id)
		} else if ai.Title != expected.Title {
			t.Fatalf("Expected: %v\nGot: %v", expected.Title, ai.Title)
		} else if ai.ProductUrl != expected.ProductUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.ProductUrl, ai.ProductUrl)
		} else if ai.CoverPhotoBaseUrl != expected.CoverPhotoBaseUrl {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoBaseUrl, ai.CoverPhotoBaseUrl)
		} else if ai.CoverPhotoMediaItemId != expected.CoverPhotoMediaItemId {
			t.Fatalf("Expected: %v\nGot: %v", expected.CoverPhotoMediaItemId, ai.CoverPhotoMediaItemId)
		}
	}

}
