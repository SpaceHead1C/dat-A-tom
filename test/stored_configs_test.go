package test

import (
	"context"
	"errors"
	"math"
	"testing"

	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/test/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
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

type StoredConfigsManagerTestSuite struct {
	suite.Suite
	man  *api.StoredConfigsManager
	repo *mocks.StoredConfigRepository
}

func TestStoredConfigManager(t *testing.T) {
	suite.Run(t, new(StoredConfigsManagerTestSuite))
}

func (s *StoredConfigsManagerTestSuite) SetupTest() {
	s.man, s.repo = newTestStoredConfigsManager(s.T())
}

func (s *StoredConfigsManagerTestSuite) TestSet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	s.repo.
		On("SetStoredConfigDatawayTomID", mock.Anything, id).Return(nil).
		On("SetStoredConfigDatawayTomID", mock.Anything, idE).Return(errors.New("error"))

	type args struct {
		ctx   context.Context
		sc    domain.StoredConfig
		value any
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
		leftErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "set dataway tom ID",
			args: args{ctx: context.Background(), sc: domain.StoredConfigTomID, value: id},
		},
		{
			name:    "set error",
			args:    args{ctx: context.Background(), sc: domain.StoredConfigTomID, value: idE},
			wantErr: true,
			err:     errors.New("error"),
		},
		{
			name:    "set dataway tom ID with unexpected type error",
			args:    args{ctx: context.Background(), sc: domain.StoredConfigTomID, value: "hello"},
			wantErr: true,
			leftErr: true,
			err:     errors.New("unexpected type"),
		},
		{
			name:    "set stored config error",
			args:    args{ctx: context.Background(), sc: math.MaxUint},
			wantErr: true,
			leftErr: true,
			err:     errors.New("unexpected stored config"),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := s.man.Set(c.args.ctx, c.args.sc, c.args.value)
			if c.wantErr {
				s.Require().Error(err)
				if c.leftErr {
					right := len(c.err.Error())
					s.Require().LessOrEqual(right, len(err.Error()))
					s.Require().Greater(right, 0)
					s.Require().Equal(c.err.Error(), err.Error()[:right])
				} else {
					s.EqualError(err, c.err.Error())
				}
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *StoredConfigsManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	scv := domain.StoredConfigUUID{Value: id}
	scvE := domain.StoredConfigUUID{Value: uuid.Nil}
	ctx := context.WithValue(context.Background(), id, id)
	ctxE := context.WithValue(context.Background(), idE, idE)
	s.repo.
		On("GetStoredConfigDatawayTomID", mock.MatchedBy(func(ctx context.Context) bool { return ctx.Value(id) != nil })).Return(scv, nil).
		On("GetStoredConfigDatawayTomID", mock.MatchedBy(func(ctx context.Context) bool { return ctx.Value(idE) != nil })).Return(scvE, domain.ErrStoredConfigTomIDNotSet)

	type args struct {
		ctx context.Context
		sc  domain.StoredConfig
	}
	type testCase struct {
		name    string
		args    args
		want    domain.StoredConfigValue
		wantErr bool
		leftErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get dataway tom ID",
			args: args{ctx: ctx, sc: domain.StoredConfigTomID},
			want: scv,
		},
		{
			name:    "get dataway tom ID not set error",
			args:    args{ctx: ctxE, sc: domain.StoredConfigTomID},
			want:    scvE,
			wantErr: true,
			err:     domain.ErrStoredConfigTomIDNotSet,
		},
		{
			name:    "get unexpected stored config error",
			args:    args{ctx: ctxE, sc: math.MaxUint},
			wantErr: true,
			leftErr: true,
			err:     errors.New("unexpected stored config"),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Get(c.args.ctx, c.args.sc)
			if c.wantErr {
				s.Require().Error(err)
				if c.leftErr {
					s.Require().Nil(c.want)
					right := len(c.err.Error())
					s.Require().LessOrEqual(right, len(err.Error()))
					s.Require().Greater(right, 0)
					s.Require().Equal(c.err.Error(), err.Error()[:right])
				} else {
					s.Require().EqualError(err, c.err.Error())
				}
			} else {
				s.Require().NoError(err)
			}
			s.EqualValues(c.want, actual)
		})
	}
}
