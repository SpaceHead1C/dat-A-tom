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
	reqE := domain.AddPropertyRequest{Name: "error"}
	s.repo.
		On("AddProperty", mock.Anything, req).Return(id, nil).
		On("AddProperty", mock.Anything, reqE).Return(uuid.Nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.AddPropertyRequest
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

func (s *PropertyManagerTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	req := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	reqE := domain.UpdPropertyRequest{ID: uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")}
	prop := &domain.Property{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("UpdateProperty", mock.Anything, req).Return(prop, nil).
		On("UpdateProperty", mock.Anything, reqE).Return(nil, domain.ErrPropertyNotFound)

	type args struct {
		ctx context.Context
		req domain.UpdPropertyRequest
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Property
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "update",
			args: args{ctx: context.Background(), req: req},
			want: prop,
		},
		{
			name:    "update error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
			err:     domain.ErrPropertyNotFound,
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

func (s *PropertyManagerTestSuite) TestGet() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	prop := &domain.Property{
		ID:          id,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetProperty", mock.Anything, id).Return(prop, nil).
		On("GetProperty", mock.Anything, idE).Return(nil, domain.ErrPropertyNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Property
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), id: id},
			want: prop,
		},
		{
			name:    "get error",
			args:    args{ctx: context.Background(), id: idE},
			wantErr: true,
			err:     domain.ErrPropertyNotFound,
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

func (s *PropertyManagerTestSuite) TestGetByKey() {
	id := "12345678-1234-1234-1234-123456789012"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	req := uuid.MustParse(id)
	reqE := uuid.MustParse(idE)
	prop := &domain.Property{
		ID:          req,
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("GetProperty", mock.Anything, req).Return(prop, nil).
		On("GetProperty", mock.Anything, reqE).Return(nil, domain.ErrPropertyNotFound)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Property
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, id))},
			want: prop,
		},
		{
			name:    "get by key error",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"id":"%s"}`, idE))},
			wantErr: true,
			err:     domain.ErrPropertyNotFound,
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

func (s *PropertyManagerTestSuite) TestGetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	pss := &domain.PropertySentState{
		ID:     id,
		Sum:    "hash",
		SentAt: time.Now().UTC(),
	}
	s.repo.
		On("GetPropertySentStateForUpdate", mock.Anything, id, mock.Anything).Return(pss, nil).
		On("GetPropertySentStateForUpdate", mock.Anything, idE, mock.Anything).Return(nil, domain.ErrSentDataNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.PropertySentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), id: id},
			want: pss,
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

func (s *PropertyManagerTestSuite) TestSetSentState() {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	idE := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	sentAt := time.Now()
	req := domain.PropertySentState{ID: id, Sum: "hash", SentAt: sentAt}
	reqE := domain.PropertySentState{ID: idE, Sum: "hash", SentAt: sentAt}
	pss := &domain.PropertySentState{ID: id, Sum: "hash", SentAt: sentAt}
	err := errors.New("error")
	s.repo.
		On("SetSentProperty", mock.Anything, req, mock.Anything).Return(pss, nil).
		On("SetSentProperty", mock.Anything, reqE, mock.Anything).Return(nil, err)

	type args struct {
		ctx context.Context
		req domain.PropertySentState
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.PropertySentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: pss,
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
	reqE := domain.SendPropertyRequest{}
	s.broker.
		On("SendProperty", mock.Anything, req).Return(nil).
		On("SendProperty", mock.Anything, reqE).Return(errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.SendPropertyRequest
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

type PropertyHandlersTestSuite struct {
	suite.Suite
	man  *api.PropertyManager
	repo *mocks.PropertyRepository
}

func TestPropertyHandlers(t *testing.T) {
	suite.Run(t, new(PropertyHandlersTestSuite))
}

func (s *PropertyHandlersTestSuite) SetupTest() {
	s.man, s.repo, _ = newTestPropertyMockedManager(s.T())
}

func (s *PropertyHandlersTestSuite) TestAdd() {
	validUUID1 := "88888888-4444-4444-4444-cccccccccccc"
	validUUID2 := "11111111-1111-1111-1111-111111111111"

	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	mockReq := domain.AddPropertyRequest{
		Name:           "prop",
		Types:          []domain.Type{domain.TypeText},
		RefTypeIDs:     []uuid.UUID{},
		OwnerRefTypeID: uuid.MustParse(validUUID1),
	}
	mockReqM := domain.AddPropertyRequest{
		Name:       "prop multitype",
		Types:      []domain.Type{domain.TypeText, domain.TypeReference},
		RefTypeIDs: []uuid.UUID{uuid.MustParse(validUUID1), uuid.MustParse(validUUID2)},
	}
	mockReqE := domain.AddPropertyRequest{
		Name:       "error",
		Types:      []domain.Type{domain.TypeText},
		RefTypeIDs: []uuid.UUID{},
	}
	mockReqEPG := domain.AddPropertyRequest{
		Name:       "PG error",
		Types:      []domain.Type{domain.TypeText},
		RefTypeIDs: []uuid.UUID{uuid.MustParse(validUUID1)},
	}
	req := handlers.AddPropertyRequestSchema{
		Name:           mockReq.Name,
		Types:          []string{domain.TypeText.Code()},
		OwnerRefTypeID: validUUID1,
	}
	reqM := handlers.AddPropertyRequestSchema{
		Name:       mockReqM.Name,
		Types:      []string{domain.TypeText.Code(), domain.TypeReference.Code()},
		RefTypeIDs: []string{validUUID1, validUUID2},
	}
	reqE := handlers.AddPropertyRequestSchema{
		Name:  mockReqE.Name,
		Types: []string{"text"},
	}
	reqEPG := handlers.AddPropertyRequestSchema{
		Name:       mockReqEPG.Name,
		Types:      []string{"text"},
		RefTypeIDs: []string{validUUID1},
	}
	reqEWoT := handlers.AddPropertyRequestSchema{
		Name: mockReqE.Name,
	}
	reqEUT := handlers.AddPropertyRequestSchema{
		Name:  mockReqE.Name,
		Types: []string{"text", "hello"},
	}
	reqERTParse := handlers.AddPropertyRequestSchema{
		Name:       mockReqE.Name,
		Types:      []string{"ref"},
		RefTypeIDs: []string{validUUID1, "hello", validUUID2},
	}
	reqEOParse := handlers.AddPropertyRequestSchema{
		Name:           mockReqE.Name,
		Types:          []string{"ref"},
		OwnerRefTypeID: "hello",
	}
	s.repo.
		On("AddProperty", mock.Anything, mockReq).Return(id, nil).
		On("AddProperty", mock.Anything, mockReqM).Return(id, nil).
		On("AddProperty", mock.Anything, mockReqE).Return(uuid.Nil, errors.New("error")).
		On("AddProperty", mock.Anything, mockReqEPG).Return(uuid.Nil, domain.ErrTypesConditionNotMatchedPG)

	type args struct {
		ctx context.Context
		req handlers.AddPropertyRequestSchema
	}
	type testCase struct {
		name    string
		args    args
		want    handlers.TextResult
		wantErr bool
		likeErr bool
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
			name: "add multiple types",
			args: args{ctx: context.Background(), req: reqM},
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
			name:    "add DB error",
			args:    args{ctx: context.Background(), req: reqEPG},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
			err:     domain.ErrTypesConditionNotMatchedPG,
		},
		{
			name:    "add without types error",
			args:    args{ctx: context.Background(), req: reqEWoT},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
			likeErr: true,
			err:     domain.ErrExpected,
		},
		{
			name:    "add unknown types error",
			args:    args{ctx: context.Background(), req: reqEUT},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
		},
		{
			name:    "add parse reference type ID error",
			args:    args{ctx: context.Background(), req: reqERTParse},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
		},
		{
			name:    "add parse owner ID error",
			args:    args{ctx: context.Background(), req: reqEOParse},
			want:    handlers.TextResult{Status: http.StatusBadRequest},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := handlers.AddProperty(c.args.ctx, s.man, c.args.req)
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

func (s *PropertyHandlersTestSuite) TestUpdate() {
	name := "name"
	descr := "description"
	mockReq := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        &name,
		Description: &descr,
	}
	mockReqE := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:        &name,
		Description: &descr,
	}
	mockReqENF := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:        &name,
		Description: &descr,
	}
	req := handlers.UpdPropertyRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqE := handlers.UpdPropertyRequestSchema{
		ID:          mockReqE.ID.String(),
		Name:        mockReqE.Name,
		Description: mockReqE.Description,
	}
	reqENF := handlers.UpdPropertyRequestSchema{
		ID:          mockReqENF.ID.String(),
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	reqEName := handlers.UpdPropertyRequestSchema{
		ID:          mockReq.ID.String(),
		Description: mockReq.Description,
	}
	reqEDescr := handlers.UpdPropertyRequestSchema{
		ID:   mockReq.ID.String(),
		Name: mockReqENF.Name,
	}
	reqEParse := handlers.UpdPropertyRequestSchema{
		ID:          "hello",
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	rt := &domain.Property{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
	}
	s.repo.
		On("UpdateProperty", mock.Anything, mockReq).Return(rt, nil).
		On("UpdateProperty", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateProperty", mock.Anything, mockReqENF).Return(nil, domain.ErrPropertyNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdPropertyRequestSchema
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
			err:     domain.ErrPropertyNotFound,
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
			actual, err := handlers.UpdateProperty(c.args.ctx, s.man, c.args.req)
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

func (s *PropertyHandlersTestSuite) TestPatch() {
	id := "12345678-1234-1234-1234-123456789012"
	name := "name"
	descr := "description"
	mockReq := domain.UpdPropertyRequest{
		ID:          uuid.MustParse(id),
		Name:        &name,
		Description: &descr,
	}
	mockReqWoN := domain.UpdPropertyRequest{
		ID:          uuid.MustParse(id),
		Description: &descr,
	}
	mockReqWoD := domain.UpdPropertyRequest{
		ID:   uuid.MustParse(id),
		Name: &name,
	}
	mockReqE := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		Name:        &name,
		Description: &descr,
	}
	mockReqENF := domain.UpdPropertyRequest{
		ID:          uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		Name:        &name,
		Description: &descr,
	}
	req := handlers.UpdPropertyRequestSchema{
		ID:          mockReq.ID.String(),
		Name:        mockReq.Name,
		Description: mockReq.Description,
	}
	reqWoN := handlers.UpdPropertyRequestSchema{
		ID:          mockReq.ID.String(),
		Description: mockReq.Description,
	}
	reqWoD := handlers.UpdPropertyRequestSchema{
		ID:   mockReq.ID.String(),
		Name: mockReqENF.Name,
	}
	reqE := handlers.UpdPropertyRequestSchema{
		ID:          mockReqE.ID.String(),
		Name:        mockReqE.Name,
		Description: mockReqE.Description,
	}
	reqENF := handlers.UpdPropertyRequestSchema{
		ID:          mockReqENF.ID.String(),
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	reqEParse := handlers.UpdPropertyRequestSchema{
		ID:          "hello",
		Name:        mockReqENF.Name,
		Description: mockReqENF.Description,
	}
	prop := &domain.Property{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Name:        "name",
		Description: "description",
		Types:       []domain.Type{domain.TypeText},
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","types":["%s"],"reference_type_ids":null,"owner_reference_type_id":null}`, id, name, descr, prop.Types[0].Code()))
	s.repo.
		On("UpdateProperty", mock.Anything, mockReq).Return(prop, nil).
		On("UpdateProperty", mock.Anything, mockReqWoN).Return(prop, nil).
		On("UpdateProperty", mock.Anything, mockReqWoD).Return(prop, nil).
		On("UpdateProperty", mock.Anything, mockReqE).Return(nil, errors.New("error")).
		On("UpdateProperty", mock.Anything, mockReqENF).Return(nil, domain.ErrPropertyNotFound)

	type args struct {
		ctx context.Context
		req handlers.UpdPropertyRequestSchema
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
			err:     domain.ErrPropertyNotFound,
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
			actual, err := handlers.PatchProperty(c.args.ctx, s.man, c.args.req)
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

func (s *PropertyHandlersTestSuite) TestGet() {
	id := "12345678-1234-1234-1234-123456789012"
	idE := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	idENF := "ffffffff-ffff-ffff-ffff-ffffffffffff"
	name := "name"
	descr := "description"
	prop := &domain.Property{
		ID:          uuid.MustParse(id),
		Name:        name,
		Description: descr,
		Types:       []domain.Type{domain.TypeText},
	}
	payload := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","types":["%s"],"reference_type_ids":null,"owner_reference_type_id":null}`, id, name, descr, prop.Types[0].Code()))
	payloadE := []byte("parse property id error: ")
	s.repo.
		On("GetProperty", mock.Anything, uuid.MustParse(id)).Return(prop, nil).
		On("GetProperty", mock.Anything, uuid.MustParse(idE)).Return(nil, errors.New("error")).
		On("GetProperty", mock.Anything, uuid.MustParse(idENF)).Return(nil, domain.ErrPropertyNotFound)

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
			err:     domain.ErrPropertyNotFound,
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
			actual, err := handlers.GetProperty(c.args.ctx, s.man, c.args.id)
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
