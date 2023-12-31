package goerror

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "single error",
			args: args{
				msg: "test error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New("", tt.args.msg, nil, false)
			assert.Error(t, err)
			assert.Equal(t, tt.args.msg, err.Error())
			assert.Equal(t, NoType, GetType(err))
		})
	}
}

func TestNewWithCustomErrorType(t *testing.T) {
	var CustomType Type = "CustomError"

	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "single custom error",
			args: args{
				msg: "test custom error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New("", tt.args.msg, &CustomType, false)
			assert.Error(t, err)
			assert.Equal(t, tt.args.msg, err.Error())
			assert.Equal(t, CustomType, GetType(err))
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one level wrapped error",
			args: args{
				err: New("", "original error", nil, false),
				msg: "child error",
			},
		},
		{
			name: "two level wrapped error",
			args: args{
				err: Wrap(New("", "original error", nil, false), nil, "", "child error", nil, false),
				msg: "child error 2",
			},
		},
		{
			name: "one level wrapped generic go error",
			args: args{
				err: fmt.Errorf("original error"),
				msg: "child error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.err, nil, "", tt.args.msg, nil, false); err != nil {
				switch tt.name {
				case "one level wrapped error":
					assert.Equal(t, "child error: original error", err.Error())
					assert.Equal(t, "original error", Unwrap(err).Error())
					assert.Nil(t, Unwrap(Unwrap(err)))
				case "two level wrapped error":
					assert.Equal(t, "child error 2: child error: original error", err.Error())
					assert.Equal(t, "child error: original error", Unwrap(err).Error())
					assert.Equal(t, "original error", Unwrap(Unwrap(err)).Error())
					assert.Nil(t, Unwrap(Unwrap(Unwrap(err))))
				case "one level wrapped generic go error":
					assert.Equal(t, "child error: original error", err.Error())
					assert.Equal(t, "original error", Unwrap(err).Error())
					assert.Nil(t, Unwrap(Unwrap(err)))
				}
			}
		})
	}
}

func TestWrapWithCustomErrorType(t *testing.T) {
	var CustomType Type = "CustomError"

	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "one level wrapped error",
			args: args{
				err: New("", "original error", &CustomType, false),
				msg: "child error",
			},
			wantErr: false,
		},
		{
			name: "two level wrapped error",
			args: args{
				err: Wrap(New("", "original error", &CustomType, false), nil, "", "child error", &CustomType, false),
				msg: "child error 2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.err, nil, "", tt.args.msg, nil, false); (err != nil) != tt.wantErr {
				switch tt.name {
				case "one level wrapped error":
					assert.Equal(t, "child error: original error", err.Error())
					assert.Equal(t, "original error", Unwrap(err).Error())
					assert.Nil(t, Unwrap(Unwrap(err)))
					assert.Equal(t, NoType, GetType(err))
					assert.Equal(t, CustomType, GetType(Unwrap(err)))
				case "two level wrapped error":
					assert.Equal(t, "child error 2: child error: original error", err.Error())
					assert.Equal(t, "child error: original error", Unwrap(err).Error())
					assert.Equal(t, "original error", Unwrap(Unwrap(err)).Error())
					assert.Nil(t, Unwrap(Unwrap(Unwrap(err))))
					assert.Equal(t, NoType, GetType(err))
					assert.Equal(t, CustomType, GetType(Unwrap(err)))
					assert.Equal(t, CustomType, GetType(Unwrap(Unwrap(err))))
				}
			}
		})
	}
}

func TestSetContext(t *testing.T) {
	type args struct {
		err   error
		key   interface{}
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "add single valued context",
			args: args{
				err:   New("", "testing error context", nil, false),
				key:   "Key1",
				value: "Value1",
			},
			want: map[string]interface{}{"field": "Key1", "message": "Value1"},
		},
		{
			name: "add single valued context for generic go error",
			args: args{
				err:   fmt.Errorf("generic go error"),
				key:   "Key1",
				value: "Value1",
			},
			want: map[string]interface{}{"field": "Key1", "message": "Value1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetContext(tt.args.err, tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, GetContext(err))
		})
	}
}

func TestGetType(t *testing.T) {
	err := DBError

	var CustomErrorType Type = "CustomErrorType"
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "test get type when no type is passed",
			args: args{
				err: New("", "NoType error", nil, false),
			},
			want: NoType,
		},
		{
			name: "test get type when error type is passed",
			args: args{
				err: New("", "DBError type", &err, false),
			},
			want: DBError,
		},
		{
			name: "test get type when generic error is passed",
			args: args{
				err: fmt.Errorf("generic go error"),
			},
			want: NoType,
		},
		{
			name: "test get type when custom error type is passed",
			args: args{
				err: New("", "DBError type", &CustomErrorType, false),
			},
			want: CustomErrorType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.args.err); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetType(t *testing.T) {
	br := BadRequest
	nf := NotFound
	type args struct {
		err error
		t   Type
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Set error type for generic go error",
			args: args{
				err: fmt.Errorf("generic go error"),
				t:   DBError,
			},
			want: DBError,
		},
		{
			name: "Set different error type",
			args: args{
				err: New("", "resource not found", &br, false),
				t:   nf,
			},
			want: nf,
		},
		{
			name: "Set error type for NoType error",
			args: args{
				err: New("", "resource not found", nil, false),
				t:   br,
			},
			want: br,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetType(tt.args.err, tt.args.t)
			assert.Equal(t, tt.want, GetType(err))
		})
	}
}

func TestGetContext(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Get error context when no context is set for generic go error",
			args: args{
				err: fmt.Errorf("generic error"),
			},
			want: nil,
		},
		{
			name: "Get error context when no context is set",
			args: args{
				err: New("", "generic error", nil, false),
			},
			want: nil,
		},
		{
			name: "Get error context when no context is set",
			args: args{
				err: SetContext(New("", "request error", nil, false), "Key1", "Value1"),
			},
			want: map[string]interface{}{"field": "Key1", "message": "Value1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContext(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tErr := fmt.Errorf("test error")
	goTErr := New("", "test error", nil, false)
	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "goerror comparison for same error",
			args: args{
				err:    New("", "test error", nil, false),
				target: New("", "test error", nil, false),
			},
			want: false,
		},
		{
			name: "goerror comparison for same error",
			args: args{
				err:    goTErr,
				target: goTErr,
			},
			want: true,
		},
		{
			name: "goerror comparison for different error",
			args: args{
				err:    New("", "test error", nil, false),
				target: New("", "test error 2", nil, false),
			},
			want: false,
		},
		{
			name: "go generic error comparison for same error",
			args: args{
				err:    fmt.Errorf("test error"),
				target: fmt.Errorf("test error"),
			},
			want: false,
		},
		{
			name: "go generic error comparison for same error",
			args: args{
				err:    tErr,
				target: tErr,
			},
			want: true,
		},
		{
			name: "go generic error comparison for different error",
			args: args{
				err:    fmt.Errorf("test error"),
				target: fmt.Errorf("test error 2"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAs(t *testing.T) {

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GoError comparison for no type specified",
			args: args{
				err:    New("", "test error", nil, false),
				target: New("", "test error", nil, false),
			},
			want: true,
		},
		{
			name: "GoError comparison when error type is specified",
			args: args{
				err:    New("", "test error", &NotFound, false),
				target: New("", "test error", &NotFound, false),
			},
			want: true,
		},
		{
			name: "GoError comparison when error different type is specified",
			args: args{
				err:    New("", "test error", &BadRequest, false),
				target: New("", "test error", &NotFound, false),
			},
			want: false,
		},
		{
			name: "GoError comparison with generic go error",
			args: args{
				err:    New("", "test error", &NotFound, false),
				target: fmt.Errorf("test error"),
			},
			want: false,
		},
		{
			name: "generic go error comparison with generic go error",
			args: args{
				err:    fmt.Errorf("test error"),
				target: fmt.Errorf("test error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := As(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("As() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "getting error message",
			args: args{
				New("", "test-error", nil, false),
			},
			want: "test-error",
		},
		{
			name: "getting error message with generic go error",
			args: args{
				fmt.Errorf("test-error-1"),
			},
			want: "test-error-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Error(tt.args.err); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
