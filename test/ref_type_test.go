package test

import (
	"context"
	"testing"

	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/test/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RefTypeManagerTestSuite struct {
	suite.Suite
	man    *api.RefTypeManager
	repo   *mocks.RefTypeRepository
	broker *mocks.RefTypeBroker
}

func TestRefTypeManager(t *testing.T) {
	suite.Run(t, new(RefTypeManagerTestSuite))
}

func (s *RefTypeManagerTestSuite) SetupTest() {
	s.man, s.repo, s.broker = newTestRefTypeMockedManager(s.T())
}

func (s *RefTypeManagerTestSuite) TestAdd() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	req := domain.AddRefTypeRequest{Name: "rt"}
	s.repo.On("AddRefType", mock.Anything, req).Return(id, nil)

	type args struct {
		ctx context.Context
		req domain.AddRefTypeRequest
	}
	type testCase struct {
		name string
		args args
		want uuid.UUID
	}
	tests := []testCase{
		{
			name: "add",
			args: args{ctx: context.Background(), req: req},
			want: id,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.Add(test.args.ctx, test.args.req)
			s.Require().NoError(err)
			s.EqualValues(test.want, actual)
		})
	}
}

func (s *RefTypeManagerTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	req := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	rt := domain.RefType{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.On("UpdateRefType", mock.Anything, req).Return(&rt, nil)

	type args struct {
		ctx context.Context
		req domain.UpdRefTypeRequest
	}
	type testCase struct {
		name string
		args args
		want *domain.RefType
	}
	tests := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: &rt,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.Update(test.args.ctx, test.args.req)
			s.Require().NoError(err)
			s.EqualValues(test.want, actual)
		})
	}
}

func (s *RefTypeManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	rt := domain.RefType{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.On("GetRefType", mock.Anything, id).Return(&rt, nil)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name string
		args args
		want *domain.RefType
	}
	tests := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: &rt,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.Get(test.args.ctx, test.args.id)
			s.Require().NoError(err)
			s.EqualValues(test.want, actual)
		})
	}
}
