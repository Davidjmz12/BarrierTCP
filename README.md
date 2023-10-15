# BarrierTCP
BarrierTCD is a Go library that provides a simple and efficient implementation of a distributed barrier synchronization mechanism. 
It allows you to synchronize a group of distributed processes, ensuring that they all reach a designated synchronization point before proceeding further.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Example](#example)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Distributed Barrier**: BarrierTCD facilitates the coordination of distributed processes, allowing them to synchronize their execution.
- **Simple API**: BarrierTCD provides a straightforward API for creating and using barriers in your Go applications.
- **Configurable**: You can customize the behavior of BarrierTCD, such as the number of participants required to trigger the barrier.

## Installation

To use BarrierTCD in your Go project, you can use the `go get` command:

```sh
go get github.com/Davidjmz12/BarrierTCP
```

## Usage
To use BarrierTCD in your Go application, you need to import the library and create a barrier. The basic usage involves 
initializing a barrier with the desired number of participants and then using it for synchronization.

### Example

In the following example, it can be seen all the main features of our library. We have two clients who communicate locally.
The first one starts a non-trigger block where the trigger thread *waitSentThree* notifies that the barrier should be starting
but he will not be blocked.
```go
package bar

import (
	"fmt"
	"github.com/DistributedClocks/GoVector/govec"
	"testing"
	"time"
)

func simulateClient(bm *BarrierManager, t time.Duration) {
	go waitSendThree(bm, t)
	bm.StartBarrier()
}

func waitSendThree(bm *BarrierManager, t time.Duration) {
	for i := 1; i <= 3; i++ {
		time.Sleep(1000 * t * time.Millisecond)
		bm.Ready(false)
	}
	fmt.Println("I ACCESS")
}

func waitBlock(bm *BarrierManager, t time.Duration) {

	go func() {
		time.Sleep(1000 * t * time.Millisecond)
		bm.Ready(true)
		fmt.Println("I ACCESS")
	}()

	bm.StartBarrier()
}

func TestBlockingBarrier(t *testing.T) {
	logger1 := govec.InitGoVector("Process_1", "log/LogFile_1", govec.GetDefaultConfig())
	logger2 := govec.InitGoVector("Process_2", "log/LogFile_2", govec.GetDefaultConfig())
	bm1 := New(3000, 3, "./barrier.txt", 1, "", logger1)
	bm2 := New(3000, 1, "./barrier.txt", 2, "", logger2)

	go simulateClient(&bm1, 4)

	waitBlock(&bm2, 1)

}
```
In the second client, not only have we blocked the main thread which is waiting for the barrier to be done,
but also the one that triggers the start of the barrier.

Frequently, it does not matter at all if the triggering process is blocked or not. For example, imagine
we have a TCP server that first listens, then triggers the barrier, and keeps waiting for messages. 
In this case, we need a barrier to ensure that everyone is listening but we actually do not care if the
receiving thread is waiting or not (because it will be blocked anyway).

## Contributing

Contributions to BarrierTCD are welcome! 
If you have any suggestions, feature requests, or bug reports, please open an issue on the GitHub repository. 
If you'd like to contribute code, please fork the repository, create a new branch, and submit a pull request.

## License
BarrierTCD is licensed under the GNU License. See the [LICENSE](LICENSE) file for details.
