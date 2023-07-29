package test

import (
	"datatom/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
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
	tests := []testCase{
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
	for _, test := range tests {
		s.Run(test.name, func() {
			out, err := domain.ValueAsJSON(test.args.v, test.args.t)
			s.Require().NoError(err, "domain.ValueAsJSON(%v, %v)", test.args.v, test.args.t)
			s.EqualValues(test.want, out, "domain.ValueAsJSON(%v, %v)", test.args.v, test.args.t)
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
	tests := []testCase{
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
	for _, test := range tests {
		s.Run(test.name, func() {
			out, err := domain.ValueAsJSON(test.args.v, test.args.t)
			s.Require().Error(err, "domain.ValueAsJSON(%v, %v)", test.args.v, test.args.t)
			s.Require().ErrorIs(err, test.err, "domain.ValueAsJSON(%v, %v)", test.args.v, test.args.t)
			s.Nil(out, "domain.ValueAsJSON(%v, %v)", test.args.v, test.args.t)
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
	tests := []testCase{
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
	for _, test := range tests {
		s.Run(test.name, func() {
			out, err := domain.ValidatedValue(test.args.v, test.args.t)
			s.Require().NoError(err, "domain.ValidatedValue(%v, %v)", test.args.v, test.args.t)
			s.EqualValues(test.want, out, "domain.ValidatedValue(%v, %v)", test.args.v, test.args.t)
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
	tests := []testCase{
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
	for _, test := range tests {
		s.Run(test.name, func() {
			out, err := domain.ValidatedValue(test.args.v, test.args.t)
			s.Require().Error(err, "domain.ValidatedValue(%v, %v)", test.args.v, test.args.t)
			s.Require().ErrorIs(err, test.err, "domain.ValidatedValue(%v, %v)", test.args.v, test.args.t)
			s.Nil(out, "domain.ValidatedValue(%v, %v)", test.args.v, test.args.t)
		})
	}
}
