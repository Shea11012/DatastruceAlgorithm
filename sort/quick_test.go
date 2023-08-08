package sort

import (
	"reflect"
	"testing"
)

func TestQuick_Sort(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "mix",
			args: args{arr: []int{10, 8, 23, 11, 13, 5, 4, 6, 15, 22}},
			want: []int{4, 5, 6, 8, 10, 11, 13, 15, 22, 23},
		},
		{
			name: "sorted",
			args: args{arr: []int{4, 5, 6, 8, 10, 11, 13, 15, 22, 23}},
			want: []int{4, 5, 6, 8, 10, 11, 13, 15, 22, 23},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quick{}
			if got := q.Sort(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quick.Sort() = %v, want %v", got, tt.want)
			}
		})
	}

}
