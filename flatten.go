package contentful

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	linkTypeAsset = "Asset"
	linkTypeEntry = "Entry"
	linkType      = "Link"
)

type sys struct {
	Type     string `json:"type"`
	LinkType string `json:"linkType"`
	ID       string `json:"id"`
}

type itemInfo struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	ContentType struct {
		Sys sys `json:"sys"`
	} `json:"contentType"`
	Revision  int       `json:"revision"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Locale    string    `json:"locale"`
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

func flattenItems(includes includes, items []item) ([]map[string]interface{}, error) {
	flattenedItems := make([]map[string]interface{}, len(items))
	for i, item := range items {
		flattenedItem, err := flattenItem(includes, item)
		if err != nil {
			return flattenedItems, err
		}

		flattenedItems[i] = flattenedItem
	}

	return flattenedItems, nil
}

func flattenItem(includes includes, item item) (map[string]interface{}, error) {
	flattenedFields := make(map[string]interface{}, len(item.Fields))

	for key, field := range item.Fields {
		flattenedField, err := flattenField(includes, field)
		if err != nil {
			return flattenedFields, err
		}
		flattenedFields[key] = flattenedField
	}

	flattenedFields["contentful_id"] = item.Sys.ID
	flattenedFields["contentful_contentType"] = item.Sys.ContentType.Sys.ID
	flattenedFields["contentful_revision"] = item.Sys.Revision
	flattenedFields["contentful_createdAt"] = item.Sys.CreatedAt
	flattenedFields["contentful_updatedAt"] = item.Sys.UpdatedAt
	flattenedFields["contentful_locale"] = item.Sys.Locale

	return flattenedFields, nil
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
//   }
//
// where includes object is
//   "includes": {
//     "Entry": [
//       {
//         "sys": {
//           "type": "Entry",
//           "id": "entryID"
//         },
//         "fields": {
//           "key": "value"
//         }
//       },
//     ],
//     "Asset": [...]
//   }
//
// and where field is
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
func flattenField(includes includes, field interface{}) (interface{}, error) {
	switch t := field.(type) {
	// Either multiple references or values, flatten each individually
	case []interface{}:
		flattenedFields := make([]interface{}, len(t))
		for i, v := range t {
			flattenedField, err := flattenField(includes, v)
			if err != nil {
				return field, err
			}
			flattenedFields[i] = flattenedField
		}
		return flattenedFields, nil

	// Reference or an object
	case map[string]interface{}:
		// Reference
		if sys, ok := parseToSys(t["sys"]); ok {
			return fetchReference(includes, sys)
		}

		// Field is not a reference but an object. Flatten like as if it were an item in search result.
		flattenedItem, err := flattenItem(includes, item{Fields: t})
		if err != nil {
			return field, err
		}
		return flattenedItem, nil

	// Plain value
	default:
		return field, nil
	}
}

func parseToSys(field interface{}) (sys, bool) {
	mapSys, ok := field.(map[string]interface{})
	if !ok {
		return sys{}, false
	}

	id, ok := mapSys["id"].(string)
	if !ok {
		return sys{}, false
	}
	linkTypeValue, ok := mapSys["linkType"].(string)
	if !ok {
		return sys{}, false
	}
	typeValue, ok := mapSys["type"].(string)
	if !ok {
		return sys{}, false
	}

	sys := sys{
		ID:       id,
		LinkType: linkTypeValue,
		Type:     typeValue,
	}

	return sys, sys.ID != "" && (sys.LinkType == linkTypeAsset || sys.LinkType == linkTypeEntry) && sys.Type == linkType
}

func fetchReference(includes includes, sys sys) (interface{}, error) {
	var (
		references []item
		found      bool
		item       item
	)

	if sys.LinkType == linkTypeEntry {
		references = includes.Entry
	} else if sys.LinkType == linkTypeAsset {
		references = includes.Asset
	} else {
		return struct{}{}, fmt.Errorf("link type is not %s or %s, but instead %s", linkTypeEntry, linkTypeAsset, sys.LinkType)
	}

	for _, ref := range references {
		if ref.Sys.ID == sys.ID && ref.Sys.Type == sys.LinkType {
			found = true
			item = ref
			break
		}
	}

	if !found {
		// TODO: try to fetch data separately
		refString := "Could not convert to string"
		bytes, err := json.MarshalIndent(references, "", "  ")
		if err == nil {
			refString = string(bytes)
		}
		return struct{}{}, fmt.Errorf("could not find a reference with type %s and with id %s.\nReferences:\n%s", sys.LinkType, sys.ID, refString)
	}

	return flattenField(includes, item.Fields)
}
