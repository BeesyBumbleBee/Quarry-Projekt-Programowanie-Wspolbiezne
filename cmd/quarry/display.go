package main

type (
	WorkerMessage struct {
		WorkerId int
		Position int
		Enqueued bool
	}

	MessageStack struct {
		top    *node
		length int
	}

	node struct {
		value WorkerMessage
		prev  *node
	}
)

func NewMessageStack() *MessageStack {
	return &MessageStack{nil, 0}
}

func (ms *MessageStack) Pop() (WorkerMessage, bool) {
	if ms.length == 0 || ms.top == nil {
		return WorkerMessage{}, false
	}
	msg := ms.top
	ms.top = msg.prev
	ms.length--
	return msg.value, true
}

func (ms *MessageStack) Push(value WorkerMessage) {
	ms.top = &node{value, ms.top}
	ms.length++
}

type WorkerPositionController struct {
	WorkersChannel chan WorkerMessage
	Buffer         *MessageStack
	DisplayChannel chan WorkerMessage
}

func NewWorkerPositionController(workerAmount int) *WorkerPositionController {
	return &WorkerPositionController{
		WorkersChannel: make(chan WorkerMessage, workerAmount),
		Buffer:         NewMessageStack(),
		DisplayChannel: make(chan WorkerMessage),
	}
}

// async function
func (c *WorkerPositionController) ManageWorkersPositions() {
	for {
		msg := <-c.WorkersChannel
		c.Buffer.Push(msg)
	}
}

func (c *WorkerPositionController) GetEnqueuedMessages(done *bool) {
	for msg, ok := c.Buffer.Pop(); ok; msg, ok = c.Buffer.Pop() {
		c.DisplayChannel <- msg
	}
	*done = false
	c.DisplayChannel <- WorkerMessage{WorkerId: -1}
}
