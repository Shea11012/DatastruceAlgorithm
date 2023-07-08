package sort

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPQ_Insert(t *testing.T) {
	type fields struct {
		size int
		fn Cmp
	}
	type args struct {
		vlaue []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect func(*PQ)
	}{
		{
			name: "normal insert",
			fields: fields{
				size: 2,
				fn: max,
			},
			expect: func(p *PQ) {
				if len(p.data) != 4 {
					t.Errorf("wrong size: expected %d got %d\n",4,len(p.data))
				}

				if !reflect.DeepEqual(p.data,[]int{0,3,2,0}) {
					t.Errorf("PQ.data = %v, want %v", p.data,[]int{0,3,2,0})
				}
			},
			args: args{
				vlaue: []int{2,3},
			},
		},

		{
			name: "grow insert",
			fields: fields{
				size: 2,
				fn: max,
			},
			expect: func(p *PQ) {
				if len(p.data) != 2*2 {
					t.Errorf("wrong size: expected size: %d, got %d\n",2*2,len(p.data)-1)
				}
				
				if !reflect.DeepEqual(p.data,[]int{0,5,2,3}) {
					t.Errorf("PQ.data = %v, want %v", p.data,[]int{0,5,2,3})
				}
			},
			args: args{
				vlaue: []int{2,3,5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPQ(tt.fields.size,tt.fields.fn)
			for _,v := range tt.args.vlaue {
				p.Insert(v)
			}
			tt.expect(p)
		})
	}
}

func TestPQ_MAX_GetTop(t *testing.T) {
	p := NewPQ(4,max)
	p.Insert(3)
	p.Insert(2)
	p.Insert(6)
	p.Insert(9)
	p.Insert(11)
	p.Insert(8)

	value := p.GetTop()
	if value != 11 {
		t.Errorf("PQ.GetTop() expected %d, got %d\n",11,value)
	}
}

func TestPQ_MIN_GetTop(t *testing.T) {
	p := NewPQ(4,less)
	p.Insert(3)
	p.Insert(2)
	p.Insert(6)
	p.Insert(9)
	p.Insert(11)
	p.Insert(8)

	value := p.GetTop()
	if value != 2 {
		t.Errorf("PQ.GetTop() expected %d, got %d\n",2,value)
	}
}

func TestPQ_Delete(t *testing.T) {
	p := NewPQ(4,less)
	p.Insert(3)
	p.Insert(2)
	p.Insert(6)
	p.Insert(9)
	p.Insert(11)
	p.Insert(8)
	fmt.Printf("%v\n",p.data)
	value := p.Delete()
	if value != 2 {
		t.Errorf("PQ.Delete() expected %d, got %d\n",2,value)
	}

	fmt.Printf("%v\n",p.data)
	value = p.GetTop()
	if value != 3 {
		t.Errorf("PQ.GetTop expected %d, got %d\n",3,value)
	}
}
