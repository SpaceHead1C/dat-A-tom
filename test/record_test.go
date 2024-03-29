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

type RecordManagerTestSuite struct {
	suite.Suite
	man    *api.RecordManager
	repo   *mocks.RecordRepository
	broker *mocks.RecordBroker
}

func TestRecordManager(t *testing.T) {
	suite.Run(t, new(RecordManagerTestSuite))
}

func (s *RecordManagerTestSuite) SetupTest() {
	s.man, s.repo, s.broker = newTestRecordMockedManager(s.T())
}

func (s *RecordManagerTestSuite) TestAdd() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	req := domain.AddRecordRequest{Name: "prop"}
	reqE := domain.AddRecordRequest{Name: "error"}
	s.repo.
		On("AddRecord", mock.Anything, req).Return(id, nil).
		On("AddRecord", mock.Anything, reqE).Return(uuid.Nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.AddRecordRequest
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

func (s *RecordManagerTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	req := domain.UpdRecordRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	reqE := domain.UpdRecordRequest{ID: uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")}
	rec := &domain.Record{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("UpdateRecord", mock.Anything, req).Return(rec, nil).
		On("UpdateRecord", mock.Anything, reqE).Return(nil, domain.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		req domain.UpdRecordRequest
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Record
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: rec,
		},
		{
			name:    "update error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
			err:     domain.ErrRecordNotFound,
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

func (s *RecordManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	rec := &domain.Record{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetRecord", mock.Anything, id).Return(rec, nil).
		On("GetRecord", mock.Anything, idE).Return(nil, domain.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Record
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: rec,
		},
		{
			name:    "get error",
			args:    args{ctx: context.Background(), id: idE},
			wantErr: true,
			err:     domain.ErrRecordNotFound,
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

func (s *RecordManagerTestSuite) TestGetByKey() {
	id := "12345678-1234-1234-1234-123456789012"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	req := uuid.MustParse(id)
	reqE := uuid.MustParse(idE)
	rec := &domain.Record{
		ID:          req,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetRecord", mock.Anything, req).Return(rec, nil).
		On("GetRecord", mock.Anything, reqE).Return(nil, domain.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Record
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, id))},
			want: rec,
		},
		{
			name:    "get by key error",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, idE))},
			wantErr: true,
			err:     domain.ErrRecordNotFound,
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

func (s *RecordManagerTestSuite) TestGetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	rss := &domain.RecordSentState{
		ID:     id,
		Sum:    "hash",
		SentAt: time.Now().UTC(),
	}
	s.repo.
		On("GetRecordSentStateForUpdate", mock.Anything, id, mock.Anything).Return(rss, nil).
		On("GetRecordSentStateForUpdate", mock.Anything, idE, mock.Anything).Return(nil, domain.ErrSentDataNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RecordSentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), id: id},
			want: rss,
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

func (s *RecordManagerTestSuite) TestSetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	sentAt := time.Now()
	req := domain.RecordSentState{ID: id, Sum: "hash", SentAt: sentAt}
	reqE := domain.RecordSentState{ID: idE, Sum: "hash", SentAt: sentAt}
	rss := &domain.RecordSentState{ID: id, Sum: "hash", SentAt: sentAt}
	err := errors.New("error")
	s.repo.
		On("SetSentRecord", mock.Anything, req, mock.Anything).Return(rss, nil).
		On("SetSentRecord", mock.Anything, reqE, mock.Anything).Return(nil, err)

	type args struct {
		ctx context.Context
		req domain.RecordSentState
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.RecordSentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: rss,
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

func (s *RecordManagerTestSuite) TestSend() {
	req := domain.SendRecordRequest{
		Record: domain.Record{
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
	reqE := domain.SendRecordRequest{}
	s.broker.
		On("SendRecord", mock.Anything, req).Return(nil).
		On("SendRecord", mock.Anything, reqE).Return(errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.SendRecordRequest
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

func (s *RecordManagerTestSuite) TestGetSender() {
	req := domain.SendRecordRequest{
		Record: domain.Record{
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
		req domain.SendRecordRequest
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

type RecordHandlersTestSuite struct {
	suite.Suite
	man  *api.RecordManager
	repo *mocks.RecordRepository
}

func TestRecordHandlers(t *testing.T) {
	suite.Run(t, new(RecordHandlersTestSuite))
}

func (s *RecordHandlersTestSuite) SetupTest() {
	s.man, s.repo, _ = newTestRecordMockedManager(s.T())
}

func (s *RecordHandlersTestSuite) TestAdd() {
	validUUID := "88888888-4444-4444-4444-cccccccccccc"
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	mockReq := domain.AddRecordRequest{Name: "rec"}
	mockReqRT := domain.AddRecordRequest{
		Name:            "rec",
		ReferenceTypeID: uuid.MustParse(validUUID),
	}
	mockReqE := domain.AddRecordRequest{Name: "error"}
	req := handlers.AddRecordRequestSchema{Name: mockReq.Name}
	reqRT := handlers.AddRecordRequestSchema{
		Name:            mockReq.Name,
		ReferenceTypeID: validUUID,
	}
	reqE := handlers.AddRecordRequestSchema{Name: mockReqE.Name}
	reqERTParse := handlers.AddRecordRequestSchema{
		Name:            mockReq.Name,
		ReferenceTypeID: "hello",
	}
	s.repo.
		On("AddRecord", mock.Anything, mockReq).Return(id, nil).
		On("AddRecord", mock.Anything, mockReqRT).Return(id, nil).
		On("AddRecord", mock.Anything, mockReqE).Return(uuid.Nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req handlers.AddRecordRequestSchema
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
			name: "add with reference type",
			args: args{ctx: context.Background(), req: reqRT},
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
		{
			name:    "add parse reference id error",
			args:    args{ctx: context.Background(), req: reqERTParse},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.AddRecord(c.args.ctx, s.man, c.args.req)
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

func (s *RecordHandlersTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	delMark := true
	mockReq := domain.UpdRecordRequest{
		ID:           uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	mockReqE := domain.UpdRecordRequest{
		ID:           uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	mockReqENF := domain.UpdRecordRequest{
		ID:           uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	req := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Name:         mockReq.Name,
		Description:  mockReq.Description,
		DeletionMark: mockReq.DeletionMark,
	}
	reqE := handlers.UpdRecordRequestSchema{
		ID:           mockReqE.ID.String(),
		Name:         mockReqE.Name,
		Description:  mockReqE.Description,
		DeletionMark: mockReqE.DeletionMark,
	}
	reqENF := handlers.UpdRecordRequestSchema{
		ID:           mockReqENF.ID.String(),
		Name:         mockReqENF.Name,
		Description:  mockReqENF.Description,
		DeletionMark: mockReqENF.DeletionMark,
	}
	reqEName := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Description:  mockReq.Description,
		DeletionMark: mockReq.DeletionMark,
	}
	reqEDescr := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Name:         mockReq.Name,
		DeletionMark: mockReq.DeletionMark,
	}
	reqEDelMark := handlers.UpdRecordRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqEParse := handlers.UpdRecordRequestSchema{
		ID:           "hello",
		Name:         mockReq.Name,
		Description:  mockReq.Description,
		DeletionMark: mockReq.DeletionMark,
	}
	rec := &domain.Record{
		ID:           uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:         name,
		Description:  descr,
		DeletionMark: delMark,
	}
	s.repo.
		On("UpdateRecord", mock.Anything, mockReq).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateRecord", mock.Anything, mockReqENF).Return(nil, domain.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdRecordRequestSchema
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
			err:     domain.ErrRecordNotFound,
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
			name:    "update error deletion mark expected",
			args:    args{ctx: context.Background(), req: reqEDelMark},
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
			actual, err := handlers.UpdateRecord(c.args.ctx, s.man, c.args.req)
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

func (s *RecordHandlersTestSuite) TestPatch() {
	id := "12345678-1234-1234-1234-123456789012"
	name := "name"
	descr := "description"
	delMark := true
	mockReq := domain.UpdRecordRequest{
		ID:           uuid.MustParse(id),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	mockReqWoN := domain.UpdRecordRequest{
		ID:           uuid.MustParse(id),
		Description:  &descr,
		DeletionMark: &delMark,
	}
	mockReqWoD := domain.UpdRecordRequest{
		ID:           uuid.MustParse(id),
		Name:         &name,
		DeletionMark: &delMark,
	}
	mockReqWoDM := domain.UpdRecordRequest{
		ID:          uuid.MustParse(id),
		Name:        &name,
		Description: &descr,
	}
	mockReqWoAll := domain.UpdRecordRequest{ID: uuid.MustParse(id)}
	mockReqE := domain.UpdRecordRequest{
		ID:           uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	mockReqENF := domain.UpdRecordRequest{
		ID:           uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:         &name,
		Description:  &descr,
		DeletionMark: &delMark,
	}
	req := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Name:         mockReq.Name,
		Description:  mockReq.Description,
		DeletionMark: mockReq.DeletionMark,
	}
	reqWoN := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Description:  mockReq.Description,
		DeletionMark: mockReq.DeletionMark,
	}
	reqWoD := handlers.UpdRecordRequestSchema{
		ID:           mockReq.ID.String(),
		Name:         mockReq.Name,
		DeletionMark: mockReq.DeletionMark,
	}
	reqWoDM := handlers.UpdRecordRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqWoAll := handlers.UpdRecordRequestSchema{ID: mockReq.ID.String()}
	reqE := handlers.UpdRecordRequestSchema{
		ID:           mockReqE.ID.String(),
		Name:         mockReqE.Name,
		Description:  mockReqE.Description,
		DeletionMark: mockReqE.DeletionMark,
	}
	reqENF := handlers.UpdRecordRequestSchema{
		ID:           mockReqENF.ID.String(),
		Name:         mockReqENF.Name,
		Description:  mockReqENF.Description,
		DeletionMark: mockReqENF.DeletionMark,
	}
	reqEParse := handlers.UpdRecordRequestSchema{
		ID:           "hello",
		Name:         mockReqE.Name,
		Description:  mockReqE.Description,
		DeletionMark: mockReqE.DeletionMark,
	}
	rec := &domain.Record{
		ID:           uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:         name,
		Description:  descr,
		DeletionMark: delMark,
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","deletion_mark":%v,"reference_type_id":null}`, id, name, descr, delMark))
	s.repo.
		On("UpdateRecord", mock.Anything, mockReq).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqWoN).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqWoD).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqWoDM).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqWoAll).Return(rec, nil).
		On("UpdateRecord", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateRecord", mock.Anything, mockReqENF).Return(nil, domain.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdRecordRequestSchema
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
			name: "patch without name",
			args: args{ctx: context.Background(), req: reqWoN},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payload,
			},
		},
		{
			name: "patch without description",
			args: args{ctx: context.Background(), req: reqWoD},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payload,
			},
		},
		{
			name: "patch without deletion mark",
			args: args{ctx: context.Background(), req: reqWoDM},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payload,
			},
		},
		{
			name: "patch only ID",
			args: args{ctx: context.Background(), req: reqWoAll},
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
			err:     domain.ErrRecordNotFound,
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
			actual, err := handlers.PatchRecord(c.args.ctx, s.man, c.args.req)
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

func (s *RecordHandlersTestSuite) TestGet() {
	id := "12345678-1234-1234-1234-123456789012"
	idR := "11111111-1111-1111-1111-111111111111"
	idRT := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	idENF := "ffffffff-ffff-ffff-ffff-ffffffffffff"
	name := "name"
	descr := "description"
	delMark := true
	rec := &domain.Record{
		ID:           uuid.MustParse(id),
		Name:         name,
		Description:  descr,
		DeletionMark: delMark,
	}
	recRT := &domain.Record{
		ID:              uuid.MustParse(idR),
		Name:            name,
		Description:     descr,
		DeletionMark:    delMark,
		ReferenceTypeID: uuid.MustParse(idRT),
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","deletion_mark":%v,"reference_type_id":null}`, id, name, descr, delMark))
	payloadRT := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","deletion_mark":%v,"reference_type_id":"%s"}`, idR, name, descr, delMark, idRT))
	payloadE := []byte("parse record id error: ")
	s.repo.
		On("GetRecord", mock.Anything, uuid.MustParse(id)).Return(rec, nil).
		On("GetRecord", mock.Anything, uuid.MustParse(idR)).Return(recRT, nil).
		On("GetRecord", mock.Anything, uuid.MustParse(idE)).Return(nil, errors.New("error")).
		On("GetRecord", mock.Anything, uuid.MustParse(idENF)).Return(nil, domain.ErrRecordNotFound)

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
			name: "get with reference ID",
			args: args{ctx: context.Background(), id: idR},
			want: handlers.Result{
				Status:  http.StatusOK,
				Payload: payloadRT,
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
			err:     domain.ErrRecordNotFound,
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
			actual, err := handlers.GetRecord(c.args.ctx, s.man, c.args.id)
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
