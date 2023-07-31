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
	cases := []testCase{
		{
			name: "Get function for tom ID",
			sc:   domain.StoredConfigTomID,
			want: funcName(s.T(), domain.StoredConfigRepository.GetStoredConfigDatawayTomID),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			out, err := c.sc.GetFunc()
			s.Require().NoError(err, "<%s>.GetFunc()", c.sc.String())
			s.Equal(c.want, funcName(s.T(), out), "<%s>.GetFunc()", c.sc.String())
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
	cases := []testCase{
		{
			name: "Scan UUID in stored configs",
			scv:  domain.StoredConfigUUID{Value: id},
			dest: &uuid.UUID{},
			want: &id,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := c.scv.ScanStoredConfigValue(c.dest)
			s.Require().NoError(err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", c.dest)
			s.EqualValues(c.want, c.dest, "domain.StoredConfigValue.ScanStoredConfigValue(%v) of type %T", c.dest, c.dest)
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
	cases := []testCase{
		{
			name: "Scan UUID in stored configs",
			scv:  domain.StoredConfigUUID{Value: uuid.MustParse("12345678-1234-1234-1234-123456789012")},
			dest: &stringV,
			err:  domain.ErrUnexpectedType,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := c.scv.ScanStoredConfigValue(c.dest)
			s.Require().Error(err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", c.dest)
			s.Require().ErrorIs(err, c.err, "domain.StoredConfigValue.ScanStoredConfigValue(%T)", c.dest)
		})
	}
}
