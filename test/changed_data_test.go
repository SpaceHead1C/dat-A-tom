package test

import (
	"context"
	"errors"
	"testing"

	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/test/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ChangedDataManagerTestSuite struct {
	suite.Suite
	man  *api.ChangedDataManager
	repo *mocks.ChangedDataRepository
}

func TestChangedDataManager(t *testing.T) {
	suite.Run(t, new(ChangedDataManagerTestSuite))
}

func (s *ChangedDataManagerTestSuite) SetupTest() {
	s.man, s.repo = newTestChangedDataManager(s.T())
}

func (s *ChangedDataManagerTestSuite) TestGet() {
	resp := []domain.ChangedData{
		{ID: 1, DataType: domain.ChangedDataRefType, Key: []byte("rt")},
		{ID: 2, DataType: domain.ChangedDataProperty, Key: []byte("prop")},
		{ID: 3, DataType: domain.ChangedDataRecord, Key: []byte("rec")},
		{ID: 4, DataType: domain.ChangedDataValue, Key: []byte("val")},
	}
	ctxV, ctxVN, ctxVE := "V", "VN", "E"
	ctx := context.WithValue(context.Background(), ctxV, ctxV)
	ctxN := context.WithValue(context.Background(), ctxVN, ctxVN)
	ctxE := context.WithValue(context.Background(), ctxVE, ctxVE)
	fn := func(ctx context.Context) bool { return ctx.Value(ctxV) != nil }
	fnN := func(ctx context.Context) bool { return ctx.Value(ctxVN) != nil }
	fnE := func(ctx context.Context) bool { return ctx.Value(ctxVE) != nil }
	s.repo.
		On("GetChanges", mock.MatchedBy(fn)).Return(resp, nil).
		On("GetChanges", mock.MatchedBy(fnN)).Return(nil, nil).
		On("GetChanges", mock.MatchedBy(fnE)).Return(nil, errors.New("error"))

	type args struct {
		ctx context.Context
	}
	type testCase struct {
		name    string
		args    args
		want    []domain.ChangedData
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: ctx},
			want: resp,
		},
		{
			name: "get nil",
			args: args{ctx: ctxN},
		},
		{
			name:    "get error",
			args:    args{ctx: ctxE},
			wantErr: true,
			err:     errors.New("error"),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Get(c.args.ctx)
			if c.wantErr {
				s.Require().Error(err)
				s.EqualError(err, c.err.Error())
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *ChangedDataManagerTestSuite) TestPurge() {
	s.repo.
		On("PurgeChanges", mock.Anything, int64(1), mock.Anything).Return(nil).
		On("PurgeChanges", mock.Anything, int64(2), mock.Anything).Return(errors.New("error"))

	type args struct {
		ctx context.Context
		id  int64
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "purge",
			args: args{ctx: context.Background(), id: int64(1)},
		},
		{
			name:    "purge error",
			args:    args{ctx: context.Background(), id: int64(2)},
			wantErr: true,
			err:     errors.New("error"),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := s.man.Purge(c.args.ctx, c.args.id, nil)
			if c.wantErr {
				s.Require().Error(err)
				s.EqualError(err, c.err.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}
