# srq >>= Stupid Redis message Queue

It is a so stupid message queue based redis that I think it is ***Simple***

## Install

```
go get -u -v github.com/Lispre/srq
```

## Feature

1. It supports multi producer push message at the same time, and the same message enqueue *Only* once.
2. It supports multi consumer fetch message from queue at the same time, and it make sure that every consumer get the different message every time from each other.
3. It supports the priority configuration of message with **waitWeight** parameter

## Example

### producer

```
package producer
    
import (
    "fmt"
    "github.com/Lispre/srq"
)

func main() {
    conn, err := srq.NewConnection("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Printf("connect redis error\n")
        return
    }
    defer conn.Close()
    queue := srq.NewQueue("test_queue", conn)
    
    status, err := queue.push("message1", 11)
    if err != nil {
        fmt.Printf("push message error\n")
    }
    if status {
        fmt.Printf("push message success\n")
    }
    
    status, err := queue.push("message2", 22)
    if err != nil {
        fmt.Printf("push message error\n")
    }
    if status {
        fmt.Printf("push message success\n")
    }
    
    status, err := queue.push("message3", 33)
    if err != nil {
        fmt.Printf("push message error\n")
    }
    if status {
        fmt.Printf("push message success\n")
    }
}
```

### consumer

```
import (
    "fmt"
    "github.com/Lispre/srq"
)

func main() {
    conn, err := srq.NewConnection("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Printf("connect redis error\n")
        return
    }
    defer conn.Close()
    
    queue := srq.NewQueue("test_queue", conn)
    
    msg1, err := queue.Pop()
    if err != nil {
        fmt.Printf("received message error\n")
    }
    if msg1 == "" {
        fmt.Printf("No message in queue")
    }
    fmt.Printf("message is: %v", msg1)
    
    msg2, err := queue.Pop()
    if err != nil {
        fmt.Printf("received message error\n")
    }
    if msg1 == "" {
        fmt.Printf("No message in queue")
    }
    fmt.Printf("message is: %v", msg2)
    
    msg3, err := queue.Pop()
    if err != nil {
        fmt.Printf("received message error\n")
    }
    if msg1 == "" {
        fmt.Printf("No message in queue")
    }
    fmt.Printf("message is: %v", msg3)
}
```
