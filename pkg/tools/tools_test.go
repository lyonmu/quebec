package tools

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"log"

	"golang.org/x/sync/errgroup"
)

func TestUserAgent(t *testing.T) {
	raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0"
	ua := ParseUserAgent(raw)
	if ua == nil {
		log.Fatal("Expected non-nil UserAgent, got nil")
	}
	log.Printf("ua : %+v", ua)
}


func TestMain(t *testing.T)  {

	testSlice := make([]int, 0, 5)
	fmt.Println(len(testSlice), cap(testSlice))
	for i := 0; i < 10; i++ {
		testSlice = append(testSlice, i)
	}
	fmt.Println(len(testSlice), cap(testSlice))
}

func TestSyncMap(t *testing.T)  {
	var sm sync.Map
	sm.Store("key1", "value1")
	sm.Store("key2", "value2")
	value, ok := sm.Load("key1")
	if ok {
		fmt.Println("Found key1:", value)
	} else {
		fmt.Println("key1 not found")
	}

	sm.Delete("key2")
	_, ok = sm.Load("key2")
	if !ok {
		fmt.Println("key2 successfully deleted")
	}

	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true
	})
}

func TestContext(t *testing.T)  {
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	<-contextWithTimeout.Done()
	fmt.Println("Context error:", contextWithTimeout.Err())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-contextWithTimeout.Done()
	}()
	wg.Wait()

	var eg errgroup.Group
	eg.Go(func() error {
		<-contextWithTimeout.Done()
		return contextWithTimeout.Err()
	})
	if err := eg.Wait(); err != nil {
		fmt.Println("Errgroup error:", err)
	}

}

func TestErrGroup(t *testing.T)  {
		// 创建一个带有取消功能的 errgroup 和 context
	eg, ctx := errgroup.WithContext(context.Background())
	
	// 启动第一个工作协程 - 会定期检查 context 是否被取消
	eg.Go(func() error {
		for i := 0; ; i++ {
			// 检查 context 是否被取消
			select {
			case <-ctx.Done():
				fmt.Println("Worker 1: 收到取消信号，正在退出...")
				return ctx.Err() // 返回 context 的错误
			default:
				fmt.Printf("Worker 1: 正在处理任务 #%d\n", i)
				time.Sleep(500 * time.Millisecond)
			}
			
			// 模拟可能发生的错误
			if i == 5 {
				fmt.Println("Worker 1: 模拟错误发生")
				return fmt.Errorf("worker 1 遇到错误")
			}
		}
	})
	
	// 启动第二个工作协程 - 也会检查 context
	eg.Go(func() error {
		for i := 1; ; i++ {
			// 检查 context 是否被取消
			select {
			case <-ctx.Done():
				fmt.Println("Worker 2: 收到取消信号，正在退出...")
				return ctx.Err()
			default:
				fmt.Printf("Worker 2: 处理数据 #%d\n", i)
				time.Sleep(800 * time.Millisecond)
			}
		}
	})
	
	// 启动第三个工作协程 - 使用 time.After 模拟超时检查
	eg.Go(func() error {
		fmt.Println("Worker 3: 开始执行")
		
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker 3: 收到取消信号，正在退出...")
				return ctx.Err()
			case <-time.After(1 * time.Second):
				fmt.Println("Worker 3: 仍在处理中...")
			}
		}
	})
	
	// 等待所有 goroutine 完成或出错
	fmt.Println("主程序: 等待所有工作协程完成...")
	if err := eg.Wait(); err != nil {
		fmt.Printf("主程序: 检测到错误: %v\n", err)
	} else {
		fmt.Println("主程序: 所有工作协程成功完成")
	}
	
	fmt.Println("主程序: 所有协程已结束，程序退出")
}


type TaskQueue struct {
	tasks []int
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewTaskQueue() *TaskQueue {
	q := &TaskQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *TaskQueue) AddTask(task int) {
	q.mu.Lock()
	q.tasks = append(q.tasks, task)
	fmt.Printf("添加任务 %d，当前队列长度: %d\n", task, len(q.tasks))
	q.mu.Unlock()
	
	// 通知一个等待的消费者
	q.cond.Signal()
}

func (q *TaskQueue) GetTask() int {
	q.mu.Lock()
	// 等待直到队列中有任务
	for len(q.tasks) == 0 {
		fmt.Println("消费者等待任务...")
		q.cond.Wait() // 释放锁并等待，被唤醒时会重新获取锁
	}
	
	// 获取任务
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	fmt.Printf("处理任务 %d，剩余队列长度: %d\n", task, len(q.tasks))
	q.mu.Unlock()
	return task
}

func TestSyncCond(t *testing.T) {
	queue := NewTaskQueue()
	
	// 启动消费者
	go func() {
		for i := 0; i < 5; i++ {
			task := queue.GetTask()
			fmt.Printf("消费者处理完成: 任务 %d\n", task)
			time.Sleep(1 * time.Second) // 模拟处理时间
		}
	}()
	
	// 主 goroutine 作为生产者
	for i := 1; i <= 5; i++ {
		queue.AddTask(i)
		time.Sleep(800 * time.Millisecond) // 控制生产速度
	}
	
	// 等待所有任务处理完成
	time.Sleep(3 * time.Second)
}

func TestAtomic(t *testing.T)  {
	
	var counter int64 = 0
	var wg sync.WaitGroup
	numGoroutines := 5
	incrementsPerGoroutine := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	expected := int64(numGoroutines * incrementsPerGoroutine)
	fmt.Printf("Final counter value: %d, Expected: %d\n", counter, expected)
}

// 基础结构体，用于测试
type Config struct {
	Host     string
	Port     int
	MaxConns int
	Timeout  time.Duration
}

func TestAtomicValue_ConfigUpdatePattern(t *testing.T) {
	// 模拟配置更新场景
	type Service struct {
		config atomic.Value
	}
	
	service := &Service{}
	// 初始化配置
	service.config.Store(Config{
		Host:     "localhost",
		Port:     8080,
		MaxConns: 10,
		Timeout:  5 * time.Second,
	})
	
	// 模拟一个处理请求的函数，它使用当前配置
	processRequest := func() (string, int) {
		cfg := service.config.Load().(Config)
		return cfg.Host, cfg.Port
	}
	
	// 并发测试：一边处理请求，一边更新配置
	var wg sync.WaitGroup
	wg.Add(2)
	
	// 请求处理器
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			host, port := processRequest()
			t.Logf("Processing request with config: %s:%d", host, port)
			time.Sleep(5 * time.Millisecond)
		}
	}()
	
	// 配置更新器
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			newConfig := Config{
				Host:     fmt.Sprintf("server-%d.example.com", i),
				Port:     9000 + i,
				MaxConns: 20 + i*5,
				Timeout:  time.Duration(10+i) * time.Second,
			}
			t.Logf("Updating config to: %+v", newConfig)
			service.config.Store(newConfig)
			time.Sleep(30 * time.Millisecond)
		}
	}()
	
	wg.Wait()
	
	// 验证最终配置
	finalConfig := service.config.Load().(Config)
	if finalConfig.Port != 9005 {
		t.Errorf("Expected final port 9005, got %d", finalConfig.Port)
	}
}