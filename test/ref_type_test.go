package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/internal/handlers"
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

type RefTypeHandlersTestSuite struct {
	suite.Suite
	man  *api.RefTypeManager
	repo *mocks.RefTypeRepository
}

func TestRefTypeHandlers(t *testing.T) {
	suite.Run(t, new(RefTypeHandlersTestSuite))
}

func (s *RefTypeHandlersTestSuite) SetupTest() {
	s.man, s.repo, _ = newTestRefTypeMockedManager(s.T())
}

func (s *RefTypeHandlersTestSuite) TestAdd() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	mockReq := domain.AddRefTypeRequest{Name: "rt"}
	mockReqE := domain.AddRefTypeRequest{Name: "error"}
	req := handlers.AddRefTypeRequestSchema{Name: mockReq.Name}
	reqE := handlers.AddRefTypeRequestSchema{Name: mockReqE.Name}
	s.repo.
		On("AddRefType", mock.Anything, mockReq).Return(id, nil).
		On("AddRefType", mock.Anything, mockReqE).Return(uuid.Nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req handlers.AddRefTypeRequestSchema
	}
	type testCase struct {
		name    string
		args    args
		want    handlers.TextResult
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "add",
			args: args{ctx: context.Background(), req: req},
			want: handlers.TextResult{
				Payload: id.String(),
				Status:  http.StatusCreated,
			},
		},
		{
			name:    "add error",
			args:    args{ctx: context.Background(), req: reqE},
			want:    handlers.TextResult{Status: http.StatusInternalServerError},
			wantErr: true,
			err:     errors.New("error"),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.AddRefType(c.args.ctx, s.man, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
				s.Require().EqualError(err, c.err.Error())
			} else {
				s.Require().NoError(err)
			}
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *RefTypeHandlersTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	mockReq := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	mockReqE := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:        &name,
		Description: &descr,
	}
	mockReqENF := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:        &name,
		Description: &descr,
	}
	req := handlers.UpdRefTypeRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqE := handlers.UpdRefTypeRequestSchema{
		ID:          mockReqE.ID.String(),
		Name:        mockReqE.Name,
		Description: mockReqE.Description,
	}
	reqENF := handlers.UpdRefTypeRequestSchema{
		ID:          mockReqENF.ID.String(),
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	reqEName := handlers.UpdRefTypeRequestSchema{
		ID:          mockReq.ID.String(),
		Description: mockReq.Description,
	}
	reqEDescr := handlers.UpdRefTypeRequestSchema{
		ID:   mockReq.ID.String(),
		Name: mockReqENF.Name,
	}
	reqEParse := handlers.UpdRefTypeRequestSchema{
		ID:          "hello",
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	rt := &domain.RefType{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("UpdateRefType", mock.Anything, mockReq).Return(rt, nil).
		On("UpdateRefType", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateRefType", mock.Anything, mockReqENF).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdRefTypeRequestSchema
	}
	type testCase struct {
		name    string
		args    args
		want    handlers.Result
		wantErr bool
		likeErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: handlers.Result{Status: http.StatusNoContent},
		},
		{
			name:    "update error",
			args:    args{ctx: context.Background(), req: reqE},
			want:    handlers.Result{Status: http.StatusInternalServerError},
			wantErr: true,
			err:     errors.New("error"),
		},
		{
			name:    "update error not found",
			args:    args{ctx: context.Background(), req: reqENF},
			want:    handlers.Result{Status: http.StatusNotFound},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
		},
		{
			name:    "update error name expected",
			args:    args{ctx: context.Background(), req: reqEName},
			want:    handlers.Result{Status: http.StatusBadRequest},
			wantErr: true,
			likeErr: true,
			err:     domain.ErrExpected,
		},
		{
			name:    "update error description expected",
			args:    args{ctx: context.Background(), req: reqEDescr},
			want:    handlers.Result{Status: http.StatusBadRequest},
			wantErr: true,
			likeErr: true,
			err:     domain.ErrExpected,
		},
		{
			name:    "update parse JSON error",
			args:    args{ctx: context.Background(), req: reqEParse},
			want:    handlers.Result{Status: http.StatusBadRequest},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.UpdateRefType(c.args.ctx, s.man, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
				if c.err != nil {
					if c.likeErr {
						s.Require().ErrorIs(err, c.err)
					} else {
						s.Require().EqualError(err, c.err.Error())
					}
				}
			} else {
				s.Require().NoError(err)
			}
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *RefTypeHandlersTestSuite) TestPatch() {
	id := "12345678-1234-1234-1234-123456789012"
	name := "name"
	descr := "description"
	mockReq := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse(id),
		Name:        &name,
		Description: &descr,
	}
	mockReqE := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:        &name,
		Description: &descr,
	}
	mockReqENF := domain.UpdRefTypeRequest{
		ID:          uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:        &name,
		Description: &descr,
	}
	req := handlers.UpdRefTypeRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqE := handlers.UpdRefTypeRequestSchema{
		ID:          mockReqE.ID.String(),
		Name:        mockReqE.Name,
		Description: mockReqE.Description,
	}
	reqENF := handlers.UpdRefTypeRequestSchema{
		ID:          mockReqENF.ID.String(),
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	reqEParse := handlers.UpdRefTypeRequestSchema{
		ID:          "hello",
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	rt := &domain.RefType{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s"}`, id, name, descr))
	s.repo.
		On("UpdateRefType", mock.Anything, mockReq).Return(rt, nil).
		On("UpdateRefType", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateRefType", mock.Anything, mockReqENF).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdRefTypeRequestSchema
	}
	type testCase struct {
		name    string
		args    args
		want    handlers.Result
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "patch",
			args: args{ctx: context.Background(), req: req},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payload,
			},
		},
		{
			name:    "patch error",
			args:    args{ctx: context.Background(), req: reqE},
			want:    handlers.Result{Status: http.StatusInternalServerError},
			wantErr: true,
			err:     errors.New("error"),
		},
		{
			name:    "patch error not found",
			args:    args{ctx: context.Background(), req: reqENF},
			want:    handlers.Result{Status: http.StatusNotFound},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
		},
		{
			name:    "patch parse JSON error",
			args:    args{ctx: context.Background(), req: reqEParse},
			want:    handlers.Result{Status: http.StatusBadRequest},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.PatchRefType(c.args.ctx, s.man, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
				if c.err != nil {
					s.Require().EqualError(err, c.err.Error())
				}
			} else {
				s.Require().NoError(err)
			}
			s.EqualValues(c.want, actual)
		})
	}
}

func (s *RefTypeHandlersTestSuite) TestGet() {
	id := "12345678-1234-1234-1234-123456789012"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	idENF := "ffffffff-ffff-ffff-ffff-ffffffffffff"
	name := "name"
	descr := "description"
	rt := &domain.RefType{
		ID:          uuid.MustParse(id),
		Name:        name,
		Description: descr,
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s"}`, id, name, descr))
	payloadE := []byte("parse reference type id error: ")
	s.repo.
		On("GetRefType", mock.Anything, uuid.MustParse(id)).Return(rt, nil).
		On("GetRefType", mock.Anything, uuid.MustParse(idE)).Return(nil, errors.New("error")).
		On("GetRefType", mock.Anything, uuid.MustParse(idENF)).Return(nil, domain.ErrRefTypeNotFound)

	type args struct {
		ctx context.Context
		id  string
	}
	type testCase struct {
		name         string
		args         args
		want         handlers.Result
		wantFromLeft bool
		wantErr      bool
		err          error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payload,
			},
		},
		{
			name:    "get error",
			args:    args{ctx: context.Background(), id: idE},
			want:    handlers.Result{Status: http.StatusInternalServerError},
			wantErr: true,
			err:     errors.New("error"),
		},
		{
			name:    "get error not found",
			args:    args{ctx: context.Background(), id: idENF},
			want:    handlers.Result{Status: http.StatusNotFound},
			wantErr: true,
			err:     domain.ErrRefTypeNotFound,
		},
		{
			name: "get error parse ID",
			args: args{ctx: context.Background(), id: "hello"},
			want: handlers.Result{
				Status:  http.StatusBadRequest,
				Payload: payloadE,
			},
			wantFromLeft: true,
			wantErr:      true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.GetRefType(c.args.ctx, s.man, c.args.id)
			if c.wantErr {
				s.Require().Error(err)
				if c.err != nil {
					s.Require().EqualError(err, c.err.Error())
				}
			} else {
				s.Require().NoError(err)
			}
			if c.wantFromLeft {
				right := len(c.want.Payload)
				s.Require().LessOrEqual(right, len(actual.Payload))
				s.Require().Greater(right, 0)
				actual.Payload = actual.Payload[:right]
			}
			s.EqualValues(c.want, actual)
		})
	}
}
