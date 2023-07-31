package test

import (
	"context"
	"testing"
	"time"

	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"datatom/test/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PropertyManagerTestSuite struct {
	suite.Suite
	man    *api.PropertyManager
	repo   *mocks.PropertyRepository
	broker *mocks.PropertyBroker
}

func TestPropertyManager(t *testing.T) {
	suite.Run(t, new(PropertyManagerTestSuite))
}

func (s *PropertyManagerTestSuite) SetupTest() {
	s.man, s.repo, s.broker = newTestPropertyMockedManager(s.T())
}

func (s *PropertyManagerTestSuite) TestAdd() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	req := domain.AddPropertyRequest{Name: "prop"}
	s.repo.On("AddProperty", mock.Anything, req).Return(id, nil)

	type args struct {
		ctx context.Context
		req domain.AddPropertyRequest
	}
	type testCase struct {
		name string
		args args
		want uuid.UUID
	}
	cases := []testCase{
		{
			name: "add",
			args: args{ctx: context.Background(), req: req},
			want: id,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Add(c.args.ctx, c.args.req)
			s.Require().NoError(err)
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *PropertyManagerTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	req := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	prop := domain.Property{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.On("UpdateProperty", mock.Anything, req).Return(&prop, nil)

	type args struct {
		ctx context.Context
		req domain.UpdPropertyRequest
	}
	type testCase struct {
		name string
		args args
		want *domain.Property
	}
	cases := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: &prop,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Update(c.args.ctx, c.args.req)
			s.Require().NoError(err)
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *PropertyManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	rt := domain.Property{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.On("GetProperty", mock.Anything, id).Return(&rt, nil)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name string
		args args
		want *domain.Property
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: &rt,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Get(c.args.ctx, c.args.id)
			s.Require().NoError(err)
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *PropertyManagerTestSuite) TestGetByKey() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	prop := domain.Property{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.On("GetProperty", mock.Anything, id).Return(&prop, nil)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Property
		wantErr bool
	}
	cases := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(`{"id":"12345678-1234-1234-1234-123456789012"}`)},
			want: &prop,
		},
		{
			name:    "get by key as invalid UUID",
			args:    args{ctx: context.Background(), key: []byte(`{"id":"uuid"}`)},
			wantErr: true,
		},
		{
			name:    "get by key without ID",
			args:    args{ctx: context.Background(), key: []byte(`{"some":"thing"}`)},
			wantErr: true,
		},
		{
			name:    "get by key as invalid JSON",
			args:    args{ctx: context.Background(), key: []byte("key")},
			wantErr: true,
		},
		{
			name:    "get by nil key",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.GetByKey(c.args.ctx, c.args.key)
			if c.wantErr {
				s.Require().Error(err)
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *PropertyManagerTestSuite) TestGetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	ps := domain.PropertySentState{
		ID:     id,
		Sum:    "hash",
		SentAt: time.Now().UTC(),
	}
	s.repo.On("GetPropertySentStateForUpdate", mock.Anything, id, mock.Anything).Return(&ps, nil)

	type args struct {
		ctx context.Context
		id  uuid.UUID
		tx  db.Transaction
	}
	type testCase struct {
		name string
		args args
		want *domain.PropertySentState
	}
	cases := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), id: id},
			want: &ps,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.GetSentState(c.args.ctx, c.args.id, c.args.tx)
			s.Require().NoError(err)
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *PropertyManagerTestSuite) TestSetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	sentAt := time.Now()
	req := domain.PropertySentState{ID: id, Sum: "hash", SentAt: sentAt}
	resp := domain.PropertySentState{ID: id, Sum: "hash", SentAt: sentAt}
	s.repo.On("SetSentProperty", mock.Anything, req, mock.Anything).Return(&resp, nil)

	type args struct {
		ctx context.Context
		req domain.PropertySentState
		tx  db.Transaction
	}
	type testCase struct {
		name string
		args args
		want *domain.PropertySentState
	}
	cases := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: &resp,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.SetSentState(c.args.ctx, c.args.req, c.args.tx)
			s.Require().NoError(err)
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *PropertyManagerTestSuite) TestSend() {
	req := domain.SendPropertyRequest{
		Property: domain.Property{
			ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			Name:        "name",
			Description: "description",
			Sum:         "hash",
			ChangeAt:    time.Now().UTC(),
		},
		TomID:       uuid.MustParse("88888888-4444-4444-4444-cccccccccccc"),
		Exchange:    "exhange",
		RoutingKeys: []string{"routing.key"},
	}
	s.broker.On("SendProperty", mock.Anything, req).Return(nil)

	type args struct {
		ctx context.Context
		req domain.SendPropertyRequest
	}
	type testCase struct {
		name string
		args args
	}
	cases := []testCase{
		{
			name: "send",
			args: args{ctx: context.Background(), req: req},
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := s.man.Send(c.args.ctx, c.args.req)
			s.Require().NoError(err)
		})
	}
}

func (s *PropertyManagerTestSuite) TestGetSender() {
	req := domain.SendPropertyRequest{
		Property: domain.Property{
			ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			Name:        "name",
			Description: "description",
			Sum:         "hash",
			ChangeAt:    time.Now().UTC(),
		},
		TomID:       uuid.MustParse("88888888-4444-4444-4444-cccccccccccc"),
		Exchange:    "exhange",
		RoutingKeys: []string{"routing.key"},
	}

	type args struct {
		req domain.SendPropertyRequest
	}
	type testCase struct {
		name string
		args args
		want *api.Sender
	}
	cases := []testCase{
		{
			name: "get sender",
			args: args{req: req},
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual := s.man.GetSender(c.args.req)
			s.Implements(c.want, actual)
		})
	}
}
