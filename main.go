package main

const (
	error429 = 429
	noError  = 0
)

// ID -
type ID int

// Queue - самописная тупая очередь
type Queue struct {
	buffer []int
}

func (q *Queue) Push(val int) {
	q.buffer = append(q.buffer, val)
}

func (q *Queue) Top() (int, bool) {
	if len(q.buffer) > 0 {
		return q.buffer[0], false
	}
	return 0, true
}

func (q *Queue) Pop() (int, bool) {
	if len(q.buffer) > 0 {
		top := q.buffer[0]
		q.buffer = q.buffer[1:]
		return top, false
	}

	return 0, true
}

func NewQueue() Queue {
	return Queue{
		buffer: []int{},
	}
}

type Request struct {
	id   ID  // уникальный номер запроса пользователя
	time int // время запроса
}

type Response struct {
	responses map[ID][]int
}

type RatelimitData struct {
	ids       map[ID]struct{} // какие пользователи ограничены
	timeLimit int             // единица времени
	maxCount  int             // максимальное количество запросов за единицу времени
}

type Ratelimit struct {
	rateLimitData       *RatelimitData
	queues              map[ID]Queue // очередь времени, когда освобождает еще одно место для нового запроса, для каждого id
	currentRequestCount map[ID]int   // количество текущих запросов, для каждого id
}

func (r *Ratelimit) Check(req Request) bool {
	if _, inIdsList := r.rateLimitData.ids[req.id]; !inIdsList {
		return true
	}

	currentRequestCount, countExist := r.currentRequestCount[req.id]
	if !countExist {
		currentRequestCount = 1
	} else {
		currentRequestCount++
	}

	queue, queueExist := r.queues[req.id]
	if !queueExist {
		newQueue := NewQueue()
		newQueue.Push(req.time + r.rateLimitData.timeLimit)
		r.queues[req.id] = newQueue
	} else {
		queue.Push(req.time + r.rateLimitData.timeLimit)
		top, isEmpty := queue.Top()

		for top < req.time && !isEmpty {
			currentRequestCount--
			top, isEmpty = queue.Pop()
		}
	}

	r.currentRequestCount[req.id] = currentRequestCount
	if currentRequestCount <= r.rateLimitData.maxCount {
		return true
	}

	return false

}

func pseudoServer(requests []Request, responses Response, rateLimit Ratelimit) {
	for _, req := range requests {
		if rateLimit.Check(req) {
			responses.responses[req.id] = append(responses.responses[req.id], noError)
		} else {
			responses.responses[req.id] = append(responses.responses[req.id], error429)
		}
	}
}

func main() {

}
