package main

import (
	"reflect"
	"testing"
)

type TestCase struct {
	requests       []Request
	response       Response
	responseResult Response
	ratelimit      Ratelimit
}

func TestSimpleCase(t *testing.T) {
	tests := []*TestCase{
		&TestCase{
			ratelimit: Ratelimit{
				rateLimitData: &RatelimitData{
					maxCount:  2,
					timeLimit: 4,
					ids: map[ID]struct{}{
						0: struct{}{},
					},
				},
				queues:              make(map[ID]*Queue),
				currentRequestCount: make(map[ID]int),
			},
			requests: []Request{
				Request{
					id:   0,
					time: 0,
				},
				Request{
					id:   0,
					time: 3,
				},
				Request{
					id:   0,
					time: 5,
				},
				Request{
					id:   1,
					time: 0,
				},
				Request{
					id:   1,
					time: 4,
				},
			},
			response: Response{
				responses: map[ID][]int{
					0: []int{},
					1: []int{},
				},
			},
			responseResult: Response{
				responses: map[ID][]int{
					0: []int{
						noError,
						noError,
						noError,
					},
					1: []int{
						noError,
						noError,
					},
				},
			},
		},
		&TestCase{
			ratelimit: Ratelimit{
				rateLimitData: &RatelimitData{
					maxCount:  4,
					timeLimit: 4,
					ids: map[ID]struct{}{
						0: struct{}{},
						1: struct{}{},
					},
				},
				queues:              make(map[ID]*Queue),
				currentRequestCount: make(map[ID]int),
			},
			requests: []Request{
				Request{
					id:   0,
					time: 0,
				},
				Request{
					id:   0,
					time: 1,
				},
				Request{
					id:   0,
					time: 2,
				},
				Request{
					id:   0,
					time: 3,
				},
				Request{
					id:   0,
					time: 4,
				},
				Request{
					id:   1,
					time: 3,
				},
				Request{
					id:   1,
					time: 4,
				},
			},
			response: Response{
				responses: map[ID][]int{
					0: []int{},
					1: []int{},
				},
			},
			responseResult: Response{
				responses: map[ID][]int{
					0: []int{
						noError,
						noError,
						noError,
						noError,
						error429,
					},
					1: []int{
						noError,
						noError,
					},
				},
			},
		},
	}

	for idx, test := range tests {
		pseudoServer(test.requests, test.response, test.ratelimit)

		for id, resp := range test.response.responses {

			if !reflect.DeepEqual(test.responseResult.responses[id], resp) {
				t.Errorf("[%v] Unexpected response: %#v\n %#v was expected",
					idx, resp, test.responseResult.responses[id])
			}
		}
	}
}
