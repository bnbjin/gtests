package test

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

func defer_call() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	// panic("触发异常")
}

func TestDefer(t *testing.T) {
	defer_call()

	/*
		打印后
		打印中
		打印前
		panic: 触发异常
	*/
}

type student struct {
	Name string
	Age  int
}

func TestMapStrPStruct(t *testing.T) {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _, stu := range stus {
		m[stu.Name] = &stu
	}

	fmt.Print(stus)
	fmt.Print(m)
}

func TestGoRutineas1(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func TestDuckType(t *testing.T) {
	tc := Teacher{}
	tc.ShowA()
}

func TestSelectOrder(t *testing.T) {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}

	/*
		output:
		1

		would not panic
	*/
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func TestDefer2(t *testing.T) {
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1

	/*
	   20
	   0
	   1
	   1
	   2
	   0
	   1
	   1
	   10
	   0
	   1
	   1
	   1
	   0
	   1
	   1
	*/

	/*
		正确答案：
		10 1 2 3
		20 0 2 2
		2 0 2 2
		1 1 3 4

		defer调用会先计算参数，在调用时刻。
	*/
}

func TestAppend(t *testing.T) {
	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	fmt.Println(s)

	/* [0 0 0 0 0 1 2 3] */
}

type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

// 读的时候没加锁

/*
func (set *threadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		set.RLock()

		for elem := range set.s {
			ch <- elem
		}

		close(ch)
		set.RUnlock()

	}()
	return ch
}

在返回ch之前，可能就已经close掉了
*/

type People1 interface {
	Speak(string) string
}

type Stduent1 struct{}

// 如果传入的类对象是作为指针的形式，就不能编译，可能是因为golang严格的类型校验
// 如果是对象本身的类型，就可以
func (stu *Stduent1) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

/*
func TestDuckType1(t *testing.T) {
	var peo People1 = Stduent1{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
}

不能编译
*/

type People2 interface {
	Show()
}

type Student2 struct{}

func (stu *Student2) Show() {}

func live() People2 {
	var stu *Student2
	return stu
}

func TestDuckType2(t *testing.T) {
	obj := live()
	if obj == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}

// 在返回值作用了多态以后，不再与nil相等，但内部的data成员仍然是nil

func TestNilChannel(t *testing.T) {
	var cc chan int
	cc <- 1
}

func TestNilChannel2(t *testing.T) {
	var cc chan int
	val := <-cc

	fmt.Print(val)
}

func TestChan2(t *testing.T) {
	cc := make(chan int)

	go func() {
		cc <- 1
		fmt.Println(1)
	}()

	go func() {
		cc <- 2
		fmt.Println(2)
	}()

	fmt.Println(10 + <-cc)
	fmt.Println(20 + <-cc)
}

func TestMoveSlice(t *testing.T) {
	nums1 := []int{1, 2, 3, 4, 5}

	// can not assign like this
	// nums1[0:5] = nums1[1:cap(nums1)]
	nums1 = nums1[1:cap(nums1)]

	fmt.Println(nums1)
}

func moveBackward(nums []int) {
	for i := len(nums) - 1; i >= 0; i-- {
		nums[i+1] = nums[i]
	}
}

func TestMoveSliceOutOfIndex(t *testing.T) {
	nums1 := []int{1, 2, 3, 4, 5}
	// will be panic
	// moveBackward(nums1)

	fmt.Println(nums1)
}

func TestSwitch(t *testing.T) {
	switch {
	default:
		fmt.Println("defalut")
	case false:
		fmt.Println("hello 1")
	case true:
		fmt.Println("hello 2")
	case true:
		fmt.Println("hello 3")
	}
}

func TestSelect(t *testing.T) {
	c := make(chan int)
	done := make(chan int)

	// blocked
	// c <- 3

	select {
	case <-c:
		fmt.Println("cccc")
	case <-done:
		fmt.Println("done")
	default:
		fmt.Println("default")
	}
}

func TestCloseChan(t *testing.T) {
	done := make(chan interface{})

	close(done)

	select {
	case <-done:
		fmt.Println("done")
	default:
		fmt.Println("default")
	}

	select {
	default:
		fmt.Println("default")
	case <-done:
		fmt.Println("done")
	}
}

func TestBufferedChan(t *testing.T) {
	bc := make(chan int, 5)

	for i := 0; i < 10; i++ {
		fmt.Println("insertion")
		bc <- i
	}

	<-bc
}

func TestOrChannel(t *testing.T) {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func TestErrorHandlingInConcurrency(t *testing.T) {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()

		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

func TestEnclosure(t *testing.T) {
	var err error
	setError := func() <-chan interface{} {
		done := make(chan interface{})

		go func() {
			defer close(done)
			err := fmt.Errorf("it's an error")
			fmt.Println(err)
		}()

		return done
	}

	done := setError()

	<-done

	fmt.Println(err)
}

func TestString(t *testing.T) {
	var abc string

	fmt.Println(abc == "")
	fmt.Println(abc)
}
