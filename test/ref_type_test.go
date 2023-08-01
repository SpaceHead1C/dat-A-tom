package test

import (
	"context"
	"errors"
	"fmt"
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
	reqE := domain.AddRefTypeRequest{Name: "error"}
	s.repo.
		On("AddRefType", mock.Anything, req).Return(id, nil).
		On("AddRefType", mock.Anything, reqE).Return(uuid.Nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.AddRefTypeRequest
	}
	type testCase struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}
	cases := []testCase{
		{
			name: "add",
			args: args{ctx: context.Background(), req: req},
			want: id,
		},
		{
			name:    "add error",
			args:    args{ctx: context.Background(), req: reqE},
			want:    uuid.Nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Add(c.args.ctx, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			s.EqualValues(c.want, actual)
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
	reqE := domain.UpdRefTypeRequest{ID: uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")}
	rt := &domain.RefType{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("UpdateRefType", mock.Anything, req).Return(rt, nil).
		On("UpdateRefType", mock.Anything, reqE).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		req domain.UpdRefTypeRequest
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefType
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: rt,
		},
		{
			name:    "update error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Update(c.args.ctx, c.args.req)
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

func (s *RefTypeManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	rt := &domain.RefType{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetRefType", mock.Anything, id).Return(rt, nil).
		On("GetRefType", mock.Anything, idE).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefType
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: rt,
		},
		{
			name:    "get error",
			args:    args{ctx: context.Background(), id: idE},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Get(c.args.ctx, c.args.id)
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

func (s *RefTypeManagerTestSuite) TestGetByKey() {
	id := "12345678-1234-1234-1234-123456789012"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	req := uuid.MustParse(id)
	reqE := uuid.MustParse(idE)
	rt := &domain.RefType{
		ID:          req,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetRefType", mock.Anything, req).Return(rt, nil).
		On("GetRefType", mock.Anything, reqE).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefType
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, id))},
			want: rt,
		},
		{
			name:    "get by key error",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, idE))},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
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
				if c.err != nil {
					s.Require().EqualError(err, c.err.Error())
				}
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *RefTypeManagerTestSuite) TestGetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	rtss := &domain.RefTypeSentState{
		ID:     id,
		Sum:    "hash",
		SentAt: time.Now().UTC(),
	}
	s.repo.
		On("GetRefTypeSentStateForUpdate", mock.Anything, id, mock.Anything).Return(rtss, nil).
		On("GetRefTypeSentStateForUpdate", mock.Anything, idE, mock.Anything).Return(nil, domain.ErrSentDataNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefTypeSentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), id: id},
			want: rtss,
		},
		{
			name:    "get state error",
			args:    args{ctx: context.Background(), id: idE},
			wantErr: true,
			err:     domain.ErrSentDataNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.GetSentState(c.args.ctx, c.args.id, c.args.tx)
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

func (s *RefTypeManagerTestSuite) TestSetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	sentAt := time.Now()
	req := domain.RefTypeSentState{ID: id, Sum: "hash", SentAt: sentAt}
	reqE := domain.RefTypeSentState{ID: idE, Sum: "hash", SentAt: sentAt}
	rtss := &domain.RefTypeSentState{ID: id, Sum: "hash", SentAt: sentAt}
	err := errors.New("error")
	s.repo.
		On("SetSentRefType", mock.Anything, req, mock.Anything).Return(rtss, nil).
		On("SetSentRefType", mock.Anything, reqE, mock.Anything).Return(nil, err)

	type args struct {
		ctx context.Context
		req domain.RefTypeSentState
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RefTypeSentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: rtss,
		},
		{
			name:    "set state error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
			err:     err,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.SetSentState(c.args.ctx, c.args.req, c.args.tx)
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
	reqE := domain.SendRefTypeRequest{}
	s.broker.
		On("SendRefType", mock.Anything, req).Return(nil).
		On("SendRefType", mock.Anything, reqE).Return(errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.SendRefTypeRequest
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
	}
	cases := []testCase{
		{
			name: "send",
			args: args{ctx: context.Background(), req: req},
		},
		{
			name:    "send error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := s.man.Send(c.args.ctx, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
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
