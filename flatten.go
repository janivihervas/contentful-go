package contentful

import "time"

const (
	linkTypeAsset = "Asset"
	linkTypeEntry = "Entry"
)

type sys struct {
	Sys struct {
		Type     string `json:"type"`
		LinkType string `json:"linkType"`
		ID       string `json:"id"`
	} `json:"sys"`
}

type itemInfo struct {
	Type        string    `json:"type"`
	ID          string    `json:"id"`
	ContentType sys       `json:"contentType"`
	Revision    int       `json:"revision"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Locale      string    `json:"locale"`
}

type item struct {
	Sys    itemInfo               `json:"sys"`
	Fields map[string]interface{} `json:"fields"`
}

type includes struct {
	Entry []item `json:"entry"`
	Asset []item `json:"asset"`
}

type searchResults struct {
	Total    int      `json:"total"`
	Skip     int      `json:"skip"`
	Limit    int      `json:"limit"`
	Items    []item   `json:"items"`
	Includes includes `json:"includes"`
}

// flattenField injects the references from "includes" object to the field therefore flattening the json.
// Example:
// response in json:
//   {
//     "items": [
//       {
//         "fields": {
//           "reference": {
//             "sys": {
//               "type": "Link",
//               "linkType": "Entry",
//               "id": "entryID"
//             }
//           }
//         }
//       }
//     ],
//     "includes": {
//       "Entry": [
//         {
//           "sys": {
//             "type": "Entry",
//             "id": "entryID"
//           },
//           "fields": {
//             "key": "value"
//           }
//         },
//       ],
//       "Asset": [...]
//     }
//}
//
// where field is
//   "reference": {
//     "sys": {
//       "type": "Link",
//       "linkType": "Entry",
//       "id": "entryID"
//     }
//   }
//
// Return value will then be
//   "reference": {
//     "key": "value"
//   }
func flattenField(response searchResults, field interface{}) (interface{}, error) {
	return nil, nil
}
