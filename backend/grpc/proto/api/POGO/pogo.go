package POGO

// AlbumsInfoPOGO is a struct utilized to unmarshal the JSON
// response from Google's Google Photos API for listing albums.
// The response is a JSON object with a key of "albums" and a
// value of an array JSON values of album info
// Ex. JSON Response:
// {"albums":[
//     {
//       "id": "",
//       "title": "",
//       "productUrl": "",
//       "mediaItemsCount": "",
//       "coverPhotoBaseUrl": "",
//       "coverPhotoMediaItemId", "",
//      },
//    ]
//  }
type AlbumsInfoPOGO struct {
	AlbumsInfo    []AlbumInfo `json:"albums"`
	NextPageToken string      `json:"nextPageToken"`
}

// AlbumsInfo is a struct utilized to unmarshal the JSON
// response from Google's Google Photos API for listing albums.
// The response is a JSON object with a key of "albums" and a
// value of an array JSON values of album info. This struct
// is a single instance of the array of album infos.
// Ex. JSON Response:
// {"albums":[
//     {
//       "id": "",
//       "title": "",
//       "productUrl": "",
//       "mediaItemsCount": "",
//       "coverPhotoBaseUrl": "",
//       "coverPhotoMediaItemId", "",
//      },
//    ]
//  }
type AlbumInfo struct {
	Id                    string `json:"id"`
	Title                 string `json:"title"`
	ProductUrl            string `json:"productUrl"`
	MediaItemsCount       string `json:"mediaItemsCount"`
	CoverPhotoBaseUrl     string `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemId string `json:"coverPhotoMediaItemId"`
}

// Equals checks for and asserts equality of one AlbumsInfoPOGO
// with another AlbumsInfoPOGO
func (ai *AlbumsInfoPOGO) Equals(AIP AlbumsInfoPOGO) bool {
	for idx, aip := range ai.AlbumsInfo {
		comp := AIP.AlbumsInfo[idx]

		if !aip.Equals(comp) {
			return false
		}
	}
	return true
}

// Equals checks for and asserts equality of one AlbumsInfo
// with another AlbumsInfo
func (ai *AlbumInfo) Equals(info AlbumInfo) bool {
	if ai.Id != info.Id {
		return false
	} else if ai.Title != info.Title {
		return false
	} else if ai.ProductUrl != info.ProductUrl {
		return false
	} else if ai.MediaItemsCount != info.MediaItemsCount {
		return false
	} else if ai.CoverPhotoBaseUrl != info.CoverPhotoBaseUrl {
		return false
	} else if ai.CoverPhotoMediaItemId != info.CoverPhotoMediaItemId {
		return false
	} else {
		return true
	}
}

// PhotosInfoPOGO is a struct utilized to unmarshal the JSON
// response from Google's Google Photos API for listing photos.
// from a specific album. The response is a JSON object with
// a key of "mediaItems" and a value of an array JSON values of
// photos info
// Ex. JSON Response:
// {"mediaItems":[
//     {
//       "id": "",
//       "productUrl": "",
//       "baseUrl": "",
//       "mimeType": "",
//       "mediaMetadata": "",
//       "filename", "",
//      },
//    ]
//  }
type PhotosInfoPOGO struct {
	MediaItems    []PhotosInfo `json:"mediaItems"`
	NextPageToken string       `json:"nextPageToken"`
}

// PhotosInfo is a struct utilized to unmarshal the JSON
// response from Google's Google Photos API for listing photos.
// from a specific album. The response is a JSON object with
// a key of "mediaItems" and a value of an array JSON values of
// photos info
// Ex. JSON Response:
// {"mediaItems":[
//     {
//       "id": "",
//       "productUrl": "",
//       "baseUrl": "",
//       "mimeType": "",
//       "mediaMetadata": "",
//       "filename", "",
//      },
//    ]
//  }
type PhotosInfo struct {
	Id         string `json:"id"`
	ProductUrl string `json:"productUrl"`
	BaseUrl    string `json:"baseUrl"`
	MimeType   string `json:"mimeType"`
	Filename   string `json:"filename"`
}

type CalendarListResponse struct {
	NextPageToken string         `json:"nextPageToken"`
	NextSyncToken string         `json:"nextSyncToken"`
	Items         []CalendarList `json:"items"`
}

type CalendarList struct {
	Id                   string              `json:"id"`
	Summary              string              `json:"summary"`
	Description          string              `json:"description"`
	Location             string              `json:"location"`
	Timezone             string              `json:"timezone"`
	ColorId              string              `json:"colorId"`
	BackgroundColor      string              `json:"backgroundColor"`
	ForegroundColor      string              `json:"foregroundColor"`
	Hidden               bool                `json:"hidden"`
	Selected             bool                `json:"selected"`
	AccessRole           string              `json:"accessRole"`
	DefaultReminders     []Reminders         `json:"defaultReminders"`
	NotificationSettings NotificationSetting `json:"notificationSettings"`
	Primary              bool                `json:"primary"`
	Deleted              bool                `json:"deleted"`
}

type Reminders struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

type NotificationSetting struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	Type   string `json:"type"`
	Method string `json:"method"`
}
