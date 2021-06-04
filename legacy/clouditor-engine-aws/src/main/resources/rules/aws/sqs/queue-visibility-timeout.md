# Configure a Visibility Timeout for Queues

Checks for all queues if they have a visibility timeout of at most 30 seconds set.

```ccl
Queue has queueAttributes.VisibilityTimeout <= 30
```
