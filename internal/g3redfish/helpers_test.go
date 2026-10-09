package g3redfish

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMembers(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected []string
		err      string
	}{
		{
			name: "valid members",
			data: []interface{}{
				map[string]interface{}{
					"@odata.id": "/redfish/v1/Systems/1",
					"Name":      "System 1",
				},
			},
			expected: []string{"/redfish/v1/Systems/1"},
		},
		{
			name: "members is not an array",
			data: map[string]interface{}{"@odata.id": "/redfish/v1/Systems/1"},
			err:  "collection members must be an array",
		},
		{
			name: "member is not an object",
			data: []interface{}{float64(1)},
			err:  "collection member 0 must be an object",
		},
		{
			name: "missing id",
			data: []interface{}{map[string]interface{}{"Name": "System 1"}},
			err:  "collection member 0 is missing @odata.id",
		},
		{
			name: "id is not a string",
			data: []interface{}{map[string]interface{}{"@odata.id": float64(1)}},
			err:  "collection member 0 @odata.id must be a string",
		},
		{
			name: "id is empty",
			data: []interface{}{map[string]interface{}{"@odata.id": "  "}},
			err:  "collection member 0 @odata.id must not be empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getMembers(test.data)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestMalformedCollectionResponseReturnsError(t *testing.T) {
	var response map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(`{"Members":[1]}`), &response))

	members, err := getValueByFieldName(response, "Members")
	require.NoError(t, err)

	_, err = getMembers(members)
	require.ErrorContains(t, err, "collection member 0 must be an object")
}

func TestGetValueByFieldName(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		path     string
		expected interface{}
		err      string
	}{
		{
			name: "case insensitive nested path",
			data: map[string]interface{}{
				"Status": map[string]interface{}{"Health": "OK"},
			},
			path:     "status.health",
			expected: "OK",
		},
		{
			name: "primitive at intermediate path",
			data: map[string]interface{}{"Status": float64(1)},
			path: "Status.Health",
			err:  "Status is not an object",
		},
		{
			name: "array at intermediate path",
			data: map[string]interface{}{"Status": []interface{}{}},
			path: "Status.Health",
			err:  "Status is not an object",
		},
		{
			name: "primitive root",
			data: float64(1),
			path: "Status",
			err:  "root is not an object",
		},
		{
			name: "missing field",
			data: map[string]interface{}{},
			path: "Status",
			err:  "failed to get Status",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getValueByFieldName(test.data, test.path)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}
