# Live-coding. Go developer position interview.
## Task description
There is a certain service that accepts requests (TCP, HTTP - it doesn't matter).
The client passes to the input some object with a description of the task.

We must, having received this task, put it in a queue for processing.
We start the task for processing (we simulate useful work through time.Sleep(5*time.Second)),
if we have free handlers.

As soon as the next task is completed, we take the next task from the queue.
If the queue is empty, we wait for new tasks from clients.

The service can simultaneously process no more than N tasks.
The remaining tasks must be queued.