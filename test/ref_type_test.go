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

func (s *RefTypeManagerTestSuite) TestGetByKey() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	rt := domain.RefType{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.On("GetRefType", mock.Anything, id).Return(&rt, nil)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefType
		wantErr bool
	}
	tests := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(`{"id":"12345678-1234-1234-1234-123456789012"}`)},
			want: &rt,
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
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.GetByKey(test.args.ctx, test.args.key)
			if test.wantErr {
				s.Require().Error(err)
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(test.want, actual)
			}
		})
	}
}

func (s *RefTypeManagerTestSuite) TestGetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	rts := domain.RefTypeSentState{
		ID:     id,
		Sum:    "hash",
		SentAt: time.Now().UTC(),
	}
	s.repo.On("GetRefTypeSentStateForUpdate", mock.Anything, id, mock.Anything).Return(&rts, nil)

	type args struct {
		ctx context.Context
		id  uuid.UUID
		tx  db.Transaction
	}
	type testCase struct {
		name string
		args args
		want *domain.RefTypeSentState
	}
	tests := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), id: id},
			want: &rts,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.GetSentState(test.args.ctx, test.args.id, test.args.tx)
			s.Require().NoError(err)
			s.EqualValues(test.want, actual)
		})
	}
}

func (s *RefTypeManagerTestSuite) TestSetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	sentAt := time.Now()
	req := domain.RefTypeSentState{ID: id, Sum: "hash", SentAt: sentAt}
	resp := domain.RefTypeSentState{ID: id, Sum: "hash", SentAt: sentAt}
	s.repo.On("SetSentRefType", mock.Anything, req, mock.Anything).Return(&resp, nil)

	type args struct {
		ctx context.Context
		req domain.RefTypeSentState
		tx  db.Transaction
	}
	type testCase struct {
		name string
		args args
		want *domain.RefTypeSentState
	}
	tests := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: &resp,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual, err := s.man.SetSentState(test.args.ctx, test.args.req, test.args.tx)
			s.Require().NoError(err)
			s.EqualValues(test.want, actual)
		})
	}
}

func (s *RefTypeManagerTestSuite) TestSend() {
	req := domain.SendRefTypeRequest{
		RefType: domain.RefType{
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
	s.broker.On("SendRefType", mock.Anything, req).Return(nil)

	type args struct {
		ctx context.Context
		req domain.SendRefTypeRequest
	}
	type testCase struct {
		name string
		args args
	}
	tests := []testCase{
		{
			name: "send",
			args: args{ctx: context.Background(), req: req},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			err := s.man.Send(test.args.ctx, test.args.req)
			s.Require().NoError(err)
		})
	}
}

func (s *RefTypeManagerTestSuite) TestGetSender() {
	req := domain.SendRefTypeRequest{
		RefType: domain.RefType{
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
		req domain.SendRefTypeRequest
	}
	type testCase struct {
		name string
		args args
		want *api.Sender
	}
	tests := []testCase{
		{
			name: "get sender",
			args: args{req: req},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			actual := s.man.GetSender(test.args.req)
			s.Implements(test.want, actual)
		})
	}
}
