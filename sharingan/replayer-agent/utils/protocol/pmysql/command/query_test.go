package command

import (
	"bytes"
	"testing"

	"github.com/modern-go/parse"
	"github.com/stretchr/testify/require"
)

func TestDecodeResultSet(t *testing.T) {
	var testCase = []struct {
		netTraffic      []byte
		expectResultset *ResultSet
		expectErr       error
	}{
		{
			netTraffic: []byte{0x01, 0x00, 0x00, 0x01, 0x02, 0x36,
				0x00, 0x00, 0x02, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74,
				0x0a, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x0a,
				0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x04, 0x6e,
				0x61, 0x6d, 0x65, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x0c, 0x21, 0x00, 0x3c,
				0x00, 0x00, 0x00, 0xfd, 0x01, 0x10, 0x00, 0x00, 0x00, 0x34, 0x00, 0x00,
				0x03, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x0a, 0x64,
				0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x0a, 0x64, 0x65,
				0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x03, 0x61, 0x67, 0x65,
				0x03, 0x61, 0x67, 0x65, 0x0c, 0x3f, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x03,
				0x01, 0x10, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x04, 0x04, 0x64, 0x65,
				0x65, 0x6e, 0x02, 0x32, 0x34, 0x0d, 0x00, 0x00, 0x05, 0x09, 0x63, 0x61,
				0x69, 0x62, 0x69, 0x72, 0x64, 0x6d, 0x65, 0x02, 0x32, 0x33, 0x09, 0x00,
				0x00, 0x06, 0x05, 0x6a, 0x61, 0x6d, 0x65, 0x73, 0x02, 0x33, 0x33, 0x0b,
				0x00, 0x00, 0x07, 0x07, 0x52, 0x6f, 0x6e, 0x61, 0x6c, 0x64, 0x6f, 0x02,
				0x33, 0x34, 0x07, 0x00, 0x00, 0x08, 0xfe, 0x00, 0x00, 0x22, 0x00, 0x00,
				0x00},
			expectResultset: &ResultSet{
				Columns: []Columndef{
					{
						Schema:     "test",
						Table:      "department",
						OrgTable:   "department",
						ColName:    "name",
						OrgColName: "name",
						Charset:    33,
						ColLength:  60,
						Type:       0xfd,
						ExtraBytes: []byte{0x01, 0x10, 0x00, 0x00, 0x00},
					},
					{
						Schema:     "test",
						Table:      "department",
						OrgTable:   "department",
						ColName:    "age",
						OrgColName: "age",
						Charset:    63,
						ColLength:  10,
						Type:       0x03,
						ExtraBytes: []byte{0x01, 0x10, 0x00, 0x00, 0x00},
					},
				},
				DataSet: [][]string{
					[]string{"deen", "24"},
					[]string{"caibirdme", "23"},
					[]string{"james", "33"},
					[]string{"Ronaldo", "34"},
				},
			},
			expectErr: nil,
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		reader := bytes.NewReader(tc.netTraffic)
		src, err := parse.NewSource(reader, 30)
		should.NoError(err)
		actual, err := DecodeResultSet(src)
		should.Equal(tc.expectErr, err, "case #%d fail %+v", idx, actual)
		should.Equal(tc.expectResultset.String(), actual.String(), "case #%d fail %+v", idx, actual)
	}
}

func TestDecodeQueryReq(t *testing.T) {
	var testCase = []struct {
		raw      []byte
		sql      string
		shoudErr bool
	}{
		{
			raw: []byte{
				0x19, 0x00, 0x00, 0x00, 0x03, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x20,
				0x2a, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x64, 0x65, 0x70, 0x61, 0x72,
				0x74, 0x6d, 0x65, 0x6e, 0x74,
			},
			sql: "select * from department",
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		src, err := parse.NewSource(bytes.NewReader(tc.raw), 10)
		should.NoError(err)
		actual, err := DecodeQueryReq(src)
		if tc.shoudErr {
			should.Error(err)
		} else {
			stringer := &QueryBody{RawBody: tc.sql}
			should.Equal(stringer.String(), actual.String(), "case #%d fail", idx)
		}
	}
}
