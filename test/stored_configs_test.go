package test

import (
	"datatom/internal/domain"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type StoredConfigTestSuite struct {
	suite.Suite
}

func TestStoredConfig(t *testing.T) {
	suite.Run(t, new(StoredConfigTestSuite))
}

func (s *StoredConfigTestSuite) TestStoredConfigGetFunc() {
	type testCase struct {
		name string
		sc   domain.StoredConfig
		want string
	}
	tests := []testCase{
		{
			name: "Get function for tom ID",
			sc:   domain.StoredConfigTomID,
			want: funcName(s.T(), domain.StoredConfigRepository.GetStoredConfigDatawayTomID),
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			out, err := test.sc.GetFunc()
			s.Require().NoError(err, "<%s>.GetFunc()", test.sc.String())
			s.Equal(test.want, funcName(s.T(), out), "<%s>.GetFunc()", test.sc.String())
		})
	}
	s.Run("Unknown stored config", func() {
		sc := domain.StoredConfig(math.MaxUint)
		out, err := sc.GetFunc()
		s.Require().Error(err, "<%s>.GetFunc()", sc.String())
		s.Nil(out, "<%s>.GetFunc()", sc.String())
	})
}

func (s *StoredConfigTestSuite) TestScanStoredConfigValue() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	type testCase struct {
		name string
		scv  domain.StoredConfigValue
		dest any
		want any
	}
	tests := []testCase{
		{
			name: "Scan UUID in stored configs",
			scv:  domain.StoredConfigUUID{Value: id},
			dest: &uuid.UUID{},
			want: &id,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.scv.ScanStoredConfigValue(test.dest)
			s.Require().NoError(err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", test.dest)
			s.EqualValues(test.want, test.dest, "domain.StoredConfigValue.ScanStoredConfigValue(%v) of type %T", test.dest, test.dest)
		})
	}
}

func (s *StoredConfigTestSuite) TestScanStoredConfigValueError() {
	var stringV string
	type testCase struct {
		name string
		scv  domain.StoredConfigValue
		dest any
		err  error
	}
	tests := []testCase{
		{
			name: "Scan UUID in stored configs",
			scv:  domain.StoredConfigUUID{Value: uuid.MustParse("12345678-1234-1234-1234-123456789012")},
			dest: &stringV,
			err:  domain.ErrUnexpectedType,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.scv.ScanStoredConfigValue(test.dest)
			s.Require().Error(err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", test.dest)
			s.Require().ErrorIs(err, test.err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", test.dest)
		})
	}
}
