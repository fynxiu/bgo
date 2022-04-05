package config

import (
	"reflect"
	"testing"
)

func Test_jsonUnmarshaller(t *testing.T) {
	type args struct {
		contents []byte
		actual   *Config
		expected *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "hp",
			args: args{
				contents: []byte(`{
					"dockerizing": true,
					"buildAll": true,
					"buildPath": "buildPath",
					"services": [
						{
							"name": "foo",
							"path": "fooPath",
							"relatives": [
								"r1",
								"r2"
							]
						}
					],
					"updateHooks": [
						{
							"url":"http://foo.bar/biz",
							"Header":{"foo":["bar"]}
						}
					]
				}`),
				actual: &Config{},
				expected: &Config{
					Dockerizing: true,
					BuildAll:    true,
					BuildPath:   "buildPath",
					Services: []Service{
						{
							Name: "foo",
							Path: "fooPath",
							Relatives: []string{
								"r1", "r2",
							},
						},
					},
					UpdatedHooks: []UpdatedHook{
						{
							URL: "http://foo.bar/biz",
							Header: map[string][]string{
								"foo": {"bar"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonUnmarshaller(tt.args.contents, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonUnmarshaller() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(tt.args.expected, tt.args.actual) {
				t.Errorf("jsonUnmarshaller()  actual = %v, expected %v", tt.args.actual, tt.args.expected)
			}
		})
	}
}

func Test_yamlUnmarshaller(t *testing.T) {
	type args struct {
		contents []byte
		actual   *Config
		expected *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "hp",
			args: args{
				contents: []byte(
					`
dockerizing: true
buildAll: true
buildPath: buildPath
services:
  - name: foo
    path: fooPath
    relatives:
      - r1
      - r2
updateHooks:
  - url: http://foo.bar/biz
    header:
      foo: [bar]
`),
				actual: &Config{},
				expected: &Config{
					Dockerizing: true,
					BuildAll:    true,
					BuildPath:   "buildPath",
					Services: []Service{
						{
							Name: "foo",
							Path: "fooPath",
							Relatives: []string{
								"r1", "r2",
							},
						},
					},
					UpdatedHooks: []UpdatedHook{
						{
							URL: "http://foo.bar/biz",
							Header: map[string][]string{
								"foo": {"bar"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yamlUnmarshaller(tt.args.contents, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonUnmarshaller() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(tt.args.expected, tt.args.actual) {
				t.Errorf("jsonUnmarshaller()  actual = %v, expected %v", tt.args.actual, tt.args.expected)
			}
		})
	}
}
