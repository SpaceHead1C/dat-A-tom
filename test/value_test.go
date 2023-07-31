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

type ValueTypeTestSuite struct {
	suite.Suite
}

func TestValueType(t *testing.T) {
	suite.Run(t, new(ValueTypeTestSuite))
}

func (s *ValueTypeTestSuite) TestValueAsJSON() {
	type args struct {
		v any
		t domain.Type
	}
	type testCase struct {
		name string
		args args
		want []byte
	}
	cases := []testCase{
		{
			name: "value as text",
			args: args{"text", domain.TypeText},
			want: []byte(`{"v":"text"}`),
		},
		{
			name: "value as empty text",
			args: args{"", domain.TypeText},
			want: []byte(`{"v":""}`),
		},
		{
			name: "value as boolean",
			args: args{true, domain.TypeBool},
			want: []byte(`{"v":true}`),
		},
		{
			name: "value as date",
			args: args{time.Date(2023, 2, 13, 21, 21, 21, 0, time.UTC), domain.TypeDate},
			want: []byte(`{"v":"2023-02-13T21:21:21Z"}`),
		},
		{
			name: "value as string date",
			args: args{"2023-02-13T21:21:21-07:00", domain.TypeDate},
			want: []byte(`{"v":"2023-02-13T21:21:21-07:00"}`),
		},
		{
			name: "value as int number",
			args: args{7, domain.TypeNumber},
			want: []byte(`{"v":7}`),
		},
		{
			name: "value as real number",
			args: args{7.7, domain.TypeNumber},
			want: []byte(`{"v":7.7}`),
		},
		{
			name: "value as reference ID",
			args: args{uuid.MustParse("12345678-4321-0123-4567-123456789abc"), domain.TypeReference},
			want: []byte(`{"v":"12345678-4321-0123-4567-123456789abc"}`),
		},
		{
			name: "value as string reference ID",
			args: args{"12345678-4321-0123-4567-123456789abc", domain.TypeReference},
			want: []byte(`{"v":"12345678-4321-0123-4567-123456789abc"}`),
		},
		{
			name: "value as nil reference ID",
			args: args{uuid.Nil, domain.TypeReference},
			want: []byte(`{"v":"00000000-0000-0000-0000-000000000000"}`),
		},
		{
			name: "value as UUID",
			args: args{uuid.MustParse("12345678-4321-0123-4567-123456789abc"), domain.TypeUUID},
			want: []byte(`{"v":"12345678-4321-0123-4567-123456789abc"}`),
		},
		{
			name: "value as string UUID",
			args: args{"12345678-4321-0123-4567-123456789abc", domain.TypeUUID},
			want: []byte(`{"v":"12345678-4321-0123-4567-123456789abc"}`),
		},
		{
			name: "value as nil UUID",
			args: args{uuid.Nil, domain.TypeUUID},
			want: []byte(`{"v":"00000000-0000-0000-0000-000000000000"}`),
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			out, err := domain.ValueAsJSON(c.args.v, c.args.t)
			s.Require().NoError(err, "domain.ValueAsJSON(%v, %v)", c.args.v, c.args.t)
			s.EqualValues(c.want, out, "domain.ValueAsJSON(%v, %v)", c.args.v, c.args.t)
		})
	}
}

func (s *ValueTypeTestSuite) TestValueAsJSONError() {
	type args struct {
		v any
		t domain.Type
	}
	type testCase struct {
		name string
		args args
		err  error
	}
	cases := []testCase{
		{
			name: "value not as text",
			args: args{1337, domain.TypeText},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value not as boolean",
			args: args{"true", domain.TypeBool},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value as string date but invalid",
			args: args{"2023-02-13T21:21:21Z-07:00", domain.TypeDate},
			err:  domain.ErrParseError,
		},
		{
			name: "value not as date",
			args: args{1676348481, domain.TypeDate},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value not as number",
			args: args{"7", domain.TypeNumber},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value not as reference ID",
			args: args{nil, domain.TypeReference},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value as string reference ID but invalid",
			args: args{"12345678-4321-0123-4567-123456789xyz", domain.TypeReference},
			err:  domain.ErrParseError,
		},
		{
			name: "value not as UUID",
			args: args{nil, domain.TypeUUID},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value as string UUID but invalid",
			args: args{"12345678-4321-0123-4567-123456789xyz", domain.TypeUUID},
			err:  domain.ErrParseError,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			out, err := domain.ValueAsJSON(c.args.v, c.args.t)
			s.Require().Error(err, "domain.ValueAsJSON(%v, %v)", c.args.v, c.args.t)
			s.Require().ErrorIs(err, c.err, "domain.ValueAsJSON(%v, %v)", c.args.v, c.args.t)
			s.Nil(out, "domain.ValueAsJSON(%v, %v)", c.args.v, c.args.t)
		})
	}
}

func (s *ValueTypeTestSuite) TestValidatedValue() {
	type args struct {
		v any
		t domain.Type
	}
	type testCase struct {
		name string
		args args
		want any
	}
	cases := []testCase{
		{
			name: "value is text",
			args: args{"text", domain.TypeText},
			want: "text",
		},
		{
			name: "value is boolean",
			args: args{true, domain.TypeBool},
			want: true,
		},
		{
			name: "value is date",
			args: args{time.Date(2023, 2, 13, 21, 21, 21, 0, time.UTC), domain.TypeDate},
			want: time.Date(2023, 2, 13, 21, 21, 21, 0, time.UTC),
		},
		{
			name: "value is date as string",
			args: args{"2023-02-13T21:21:21Z", domain.TypeDate},
			want: time.Date(2023, 2, 13, 21, 21, 21, 0, time.UTC),
		},
		{
			name: "value is int number",
			args: args{7, domain.TypeNumber},
			want: 7,
		},
		{
			name: "value is real number",
			args: args{7.7, domain.TypeNumber},
			want: 7.7,
		},
		{
			name: "value is reference ID",
			args: args{uuid.MustParse("12345678-4321-0123-4567-123456789abc"), domain.TypeReference},
			want: uuid.MustParse("12345678-4321-0123-4567-123456789abc"),
		},
		{
			name: "value is reference ID as string",
			args: args{"12345678-4321-0123-4567-123456789abc", domain.TypeReference},
			want: uuid.MustParse("12345678-4321-0123-4567-123456789abc"),
		},
		{
			name: "value is nil reference ID",
			args: args{uuid.Nil, domain.TypeReference},
			want: uuid.Nil,
		},
		{
			name: "value is UUID",
			args: args{uuid.MustParse("12345678-4321-0123-4567-123456789abc"), domain.TypeUUID},
			want: uuid.MustParse("12345678-4321-0123-4567-123456789abc"),
		},
		{
			name: "value is UUID as string",
			args: args{"12345678-4321-0123-4567-123456789abc", domain.TypeUUID},
			want: uuid.MustParse("12345678-4321-0123-4567-123456789abc"),
		},
		{
			name: "value as nil UUID",
			args: args{uuid.Nil, domain.TypeUUID},
			want: uuid.Nil,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			out, err := domain.ValidatedValue(c.args.v, c.args.t)
			s.Require().NoError(err, "domain.ValidatedValue(%v, %v)", c.args.v, c.args.t)
			s.EqualValues(c.want, out, "domain.ValidatedValue(%v, %v)", c.args.v, c.args.t)
		})
	}
}

func (s *ValueTypeTestSuite) TestValidatedValueError() {
	type args struct {
		v any
		t domain.Type
	}
	type testCase struct {
		name string
		args args
		err  error
	}
	cases := []testCase{
		{
			name: "value is not text",
			args: args{1337, domain.TypeText},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is not boolean",
			args: args{"true", domain.TypeBool},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is string date but invalid",
			args: args{"2023-02-13T21:21:21Z-07:00", domain.TypeDate},
			err:  domain.ErrParseError,
		},
		{
			name: "value is not date",
			args: args{1676348481, domain.TypeDate},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is not number",
			args: args{"7", domain.TypeNumber},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is not reference ID",
			args: args{nil, domain.TypeReference},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is string reference ID but invalid",
			args: args{"12345678-4321-0123-4567-123456789xyz", domain.TypeReference},
			err:  domain.ErrParseError,
		},
		{
			name: "value is not UUID",
			args: args{nil, domain.TypeUUID},
			err:  domain.ErrUnexpectedTypePG,
		},
		{
			name: "value is string UUID but invalid",
			args: args{"12345678-4321-0123-4567-123456789xyz", domain.TypeUUID},
			err:  domain.ErrParseError,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			out, err := domain.ValidatedValue(c.args.v, c.args.t)
			s.Require().Error(err, "domain.ValidatedValue(%v, %v)", c.args.v, c.args.t)
			s.Require().ErrorIs(err, c.err, "domain.ValidatedValue(%v, %v)", c.args.v, c.args.t)
			s.Nil(out, "domain.ValidatedValue(%v, %v)", c.args.v, c.args.t)
		})
	}
}

type ValueManagerTestSuite struct {
	suite.Suite
	man    *api.ValueManager
	repo   *mocks.ValueRepository
	broker *mocks.ValueBroker
}

func TestValueManager(t *testing.T) {
	suite.Run(t, new(ValueManagerTestSuite))
}

func (s *ValueManagerTestSuite) SetupTest() {
	s.man, s.repo, s.broker = newTestValueMockedManager(s.T())
}

func (s *ValueManagerTestSuite) TestSet() {
	rID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	pID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	req := domain.SetValueRequest{RecordID: rID, PropertyID: pID, Type: domain.TypeNumber, Value: 7}
	reqErr := domain.SetValueRequest{RecordID: rID, PropertyID: pID, Type: domain.TypeBool, Value: false}
	resp := &domain.Value{RecordID: rID, PropertyID: pID}
	s.repo.
		On("SetValue", mock.Anything, req).Return(resp, nil).
		On("SetValue", mock.Anything, reqErr).Return(nil, domain.ErrUnexpectedTypePG)

	type args struct {
		ctx context.Context
		req domain.SetValueRequest
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Value
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "set",
			args: args{ctx: context.Background(), req: req},
			want: resp,
		},
		{
			name:    "set error",
			args:    args{ctx: context.Background(), req: reqErr},
			wantErr: true,
			err:     domain.ErrUnexpectedType,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Set(c.args.ctx, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
				s.ErrorIs(err, c.err)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *ValueManagerTestSuite) TestGet() {
	rID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	pID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	pIDe := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	req := domain.GetValueRequest{RecordID: rID, PropertyID: pID}
	reqENF := domain.GetValueRequest{RecordID: rID, PropertyID: pIDe}
	val := &domain.Value{RecordID: rID, PropertyID: pID}
	s.repo.
		On("GetValue", mock.Anything, req).Return(val, nil).
		On("GetValue", mock.Anything, reqENF).Return(nil, domain.ErrValueNotFound)

	type args struct {
		ctx context.Context
		req domain.GetValueRequest
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Value
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get",
			args: args{ctx: context.Background(), req: req},
			want: val,
		},
		{
			name:    "get error not found",
			args:    args{ctx: context.Background(), req: reqENF},
			wantErr: true,
			err:     domain.ErrNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.Get(c.args.ctx, c.args.req)
			if c.wantErr {
				s.Require().Error(err)
				s.ErrorIs(err, c.err)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *ValueManagerTestSuite) TestGetByKey() {
	rID := "11111111-1111-1111-1111-111111111111"
	pID := "22222222-2222-2222-2222-222222222222"
	pIDe := "33333333-3333-3333-3333-333333333333"
	rUUID := uuid.MustParse(rID)
	pUUID := uuid.MustParse(pID)
	pUUIDe := uuid.MustParse(pIDe)
	req := domain.GetValueRequest{RecordID: rUUID, PropertyID: pUUID}
	reqENF := domain.GetValueRequest{RecordID: rUUID, PropertyID: pUUIDe}
	val := &domain.Value{RecordID: rUUID, PropertyID: pUUID}
	s.repo.
		On("GetValue", mock.Anything, req).Return(val, nil).
		On("GetValue", mock.Anything, reqENF).Return(nil, domain.ErrValueNotFound)

	type args struct {
		ctx context.Context
		key []byte
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.Value
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get by key",
			args: args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"owner_id":"%s","property_id":"%s"}`, rID, pID))},
			want: val,
		},
		{
			name:    "get by key but record ID as invalid UUID",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"owner_id":"hello","property_id":"%s"}`, pID))},
			wantErr: true,
		},
		{
			name:    "get by key but property ID as invalid UUID",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"owner_id":"%s","property_id":"hello"}`, rID))},
			wantErr: true,
		},
		{
			name:    "get by key without record ID",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"property_id":"%s"}`, pID))},
			wantErr: true,
		},
		{
			name:    "get by key without property ID",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"owner_id":"%s"}`, rID))},
			wantErr: true,
		},
		{
			name:    "get by key without IDs",
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
		{
			name:    "get error not found",
			args:    args{ctx: context.Background(), key: []byte(fmt.Sprintf(`{"owner_id":"%s","property_id":"%s"}`, rID, pIDe))},
			wantErr: true,
			err:     domain.ErrNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.GetByKey(c.args.ctx, c.args.key)
			if c.wantErr {
				s.Require().Error(err)
				if c.err != nil {
					s.Require().ErrorIs(err, c.err)
				}
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *ValueManagerTestSuite) TestGetSentState() {
	rID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	pID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	pIDe := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	req := domain.GetValueRequest{RecordID: rID, PropertyID: pID}
	reqENF := domain.GetValueRequest{RecordID: rID, PropertyID: pIDe}
	vs := &domain.ValueSentState{RecordID: rID, PropertyID: pID}
	s.repo.
		On("GetValueSentStateForUpdate", mock.Anything, req, mock.Anything).Return(vs, nil).
		On("GetValueSentStateForUpdate", mock.Anything, reqENF, mock.Anything).Return(nil, domain.ErrSentDataNotFound)

	type args struct {
		ctx context.Context
		req domain.GetValueRequest
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.ValueSentState
		wantErr bool
		err     error
	}
	cases := []testCase{
		{
			name: "get state",
			args: args{ctx: context.Background(), req: req},
			want: vs,
		},
		{
			name:    "get state error not found",
			args:    args{ctx: context.Background(), req: req},
			wantErr: true,
			err:     domain.ErrNotFound,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.GetSentState(c.args.ctx, c.args.req, c.args.tx)
			if c.wantErr {
				s.Require().Error(err)
				s.Require().ErrorIs(err, c.err)
				s.Nil(actual)
			} else {
				s.Require().NoError(err)
				s.EqualValues(c.want, actual)
			}
		})
	}
}

func (s *ValueManagerTestSuite) TestSetSentState() {
	rID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	pID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	sentAt := time.Now()
	req := domain.ValueSentState{RecordID: rID, PropertyID: pID, Sum: "hash", SentAt: sentAt}
	reqE := domain.ValueSentState{RecordID: rID, PropertyID: pID, Sum: "", SentAt: sentAt}
	resp := &domain.ValueSentState{RecordID: rID, PropertyID: pID, Sum: "hash", SentAt: sentAt}
	s.repo.
		On("SetSentValue", mock.Anything, req, mock.Anything).Return(resp, nil).
		On("SetSentValue", mock.Anything, reqE, mock.Anything).Return(nil, errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.ValueSentState
		tx  db.Transaction
	}
	type testCase struct {
		name    string
		args    args
		want    *domain.ValueSentState
		wantErr bool
	}
	cases := []testCase{
		{
			name: "set state",
			args: args{ctx: context.Background(), req: req},
			want: resp,
		},
		{
			name:    "set state error",
			args:    args{ctx: context.Background(), req: reqE},
			wantErr: true,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			actual, err := s.man.SetSentState(c.args.ctx, c.args.req, c.args.tx)
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

func (s *ValueManagerTestSuite) TestSend() {
	req := domain.SendValueRequest{
		Value: domain.Value{
			RecordID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			PropertyID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Sum:        "hash",
			ChangeAt:   time.Now().UTC(),
		},
		TomID:       uuid.MustParse("88888888-4444-4444-4444-cccccccccccc"),
		Exchange:    "exhange",
		RoutingKeys: []string{"routing.key"},
	}
	reqE := domain.SendValueRequest{}
	s.broker.
		On("SendValue", mock.Anything, req).Return(nil).
		On("SendValue", mock.Anything, reqE).Return(errors.New("error"))

	type args struct {
		ctx context.Context
		req domain.SendValueRequest
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
			name: "send error",
			args: args{ctx: context.Background(), req: reqE},
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

func (s *ValueManagerTestSuite) TestGetSender() {
	req := domain.SendValueRequest{
		Value: domain.Value{
			RecordID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			PropertyID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Sum:        "hash",
			ChangeAt:   time.Now().UTC(),
		},
		TomID:       uuid.MustParse("88888888-4444-4444-4444-cccccccccccc"),
		Exchange:    "exhange",
		RoutingKeys: []string{"routing.key"},
	}

	type args struct {
		req domain.SendValueRequest
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
