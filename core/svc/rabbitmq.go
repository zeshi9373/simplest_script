package svc

import (
	"encoding/json"
	"fmt"
	"simplest_script/core/conf"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/streadway/amqp"
)

var rabbitConn *amqp.Connection
var rabbitChannel = make(map[string]chan *amqp.Channel)
var syncMutex sync.Mutex

// rabbitMQ结构体
type RabbitMQ struct {
	conn *amqp.Connection
	Host string
	//队列名称
	QueueName string
	//交换机名称
	ExchangeName string
	// 交换机类型
	ExchangeType string
	//routing Key 路由
	Key string
	//连接信息
	Mqurl string

	internal bool
	noWait   bool
	args     amqp.Table

	needAck bool
}

// 断开channel 和 connection
func RabbitMQDestory() {
	if rabbitConn != nil {
		rabbitConn.Close()
	}
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		hlog.Info(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQ(queueName string, channelSize int) *RabbitMQ {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	r := &RabbitMQ{
		QueueName: queueName,
	}

	if rabbitConn == nil || rabbitConn.IsClosed() {
		r.Mqurl = conf.Conf.RabbitMQ.Addr
		var err error
		//获取connection
		r.conn, err = amqp.Dial(r.Mqurl)
		r.failOnErr(err, "[rabbitmq] failed to connect!")

		rabbitConn = r.conn

		// 监听连接关闭事件
		go func() {
			<-rabbitConn.NotifyClose(make(chan *amqp.Error))
			rabbitConn = nil // 标记需要重新连接
		}()
	}

	if len(rabbitChannel[queueName]) == 0 {
		//获取channel
		rabbitChannel[queueName] = make(chan *amqp.Channel, channelSize)

		for i := 0; i < channelSize; i++ {
			ch, err := rabbitConn.Channel()

			if err != nil {
				r.failOnErr(err, "[rabbitmq] failed to open a channel")
			}

			rabbitChannel[queueName] <- ch
		}
	}

	return r
}

func (r *RabbitMQ) ChannelGet(queueName string) (*amqp.Channel, error) {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	select {
	case ch := <-rabbitChannel[queueName]:
		return ch, nil
	default:
		return rabbitConn.Channel()
	}
}

func (r *RabbitMQ) ChannelPut(queueName string, ch *amqp.Channel) {
	select {
	case rabbitChannel[queueName] <- ch:
	default:
		ch.Close() // 池已满，关闭通道
	}
}

func (r *RabbitMQ) SetHost(host string) *RabbitMQ {
	r.Host = host
	return r
}

func (r *RabbitMQ) SetExchangeName(exchangeName string) *RabbitMQ {
	r.ExchangeName = exchangeName
	return r
}

func (r *RabbitMQ) SetExchangeType(exchangeType string) *RabbitMQ {
	r.ExchangeType = exchangeType
	return r
}

func (r *RabbitMQ) SetRouteKey(routeKey string) *RabbitMQ {
	r.Key = routeKey
	return r
}

func (r *RabbitMQ) SetInternal(internal bool) *RabbitMQ {
	r.internal = internal
	return r
}

func (r *RabbitMQ) SetNoWait(noWait bool) *RabbitMQ {
	r.noWait = noWait
	return r
}

func (r *RabbitMQ) SetArgs(args amqp.Table) *RabbitMQ {
	r.args = args
	return r
}

func (r *RabbitMQ) SetNeedAck(needAck bool) *RabbitMQ {
	r.needAck = needAck
	return r
}

func (r *RabbitMQ) Publish(message string) {
	ch, err := r.ChannelGet(r.QueueName)

	if err != nil {
		hlog.Error("RabbitMQ channel get error:", err)
		return
	}

	if r.ExchangeName == "" {
		_, err := ch.QueueDeclare(
			r.QueueName,
			//是否持久化
			conf.Conf.RabbitMQ.Durable,
			//是否自动删除
			conf.Conf.RabbitMQ.AutoDelete,
			//是否具有排他性
			r.internal,
			//是否阻塞处理
			r.noWait,
			//额外的属性
			r.args,
		)
		if err != nil {
			r.failOnErr(err, "rabbitmq 队列声明失败")
		}
	} else {
		err := ch.ExchangeDeclare(
			r.ExchangeName,
			r.ExchangeType,
			//是否持久化
			conf.Conf.RabbitMQ.Durable,
			//是否自动删除
			conf.Conf.RabbitMQ.AutoDelete,
			//是否具有排他性
			r.internal,
			//是否阻塞处理
			r.noWait,
			//额外的属性
			r.args,
		)

		if err != nil {
			r.failOnErr(err, "rabbitmq 交换机声明失败")
		}
	}

	//调用channel 发送消息到队列中
	err = ch.Publish(
		r.ExchangeName,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		hlog.Info("发送消息失败 msg => " + message)
	}

	r.ChannelPut(r.QueueName, ch)
}

func (r *RabbitMQ) Consume(f func(string) bool) {
	ch, err := r.ChannelGet(r.QueueName)

	if err != nil {
		hlog.Error("RabbitMQ channel get error:", err)
		return
	}

	if r.ExchangeName == "" {
		_, err := ch.QueueDeclare(
			r.QueueName,
			//是否持久化
			conf.Conf.RabbitMQ.Durable,
			//是否自动删除
			conf.Conf.RabbitMQ.AutoDelete,
			//是否具有排他性
			r.internal,
			//是否阻塞处理
			r.noWait,
			//额外的属性
			r.args,
		)
		if err != nil {
			r.failOnErr(err, "rabbitmq 队列声明失败")
		}
	} else {
		err := ch.ExchangeDeclare(
			r.ExchangeName,
			r.ExchangeType,
			//是否持久化
			conf.Conf.RabbitMQ.Durable,
			//是否自动删除
			conf.Conf.RabbitMQ.AutoDelete,
			//是否具有排他性
			r.internal,
			//是否阻塞处理
			r.noWait,
			//额外的属性
			r.args,
		)

		if err != nil {
			r.failOnErr(err, "rabbitmq 交换机声明失败")
		}
	}

	//接收消息
	msgs, err := ch.Consume(
		r.QueueName, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		!r.needAck, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		s, _ := json.Marshal(msgs)
		hlog.Info("Received a message Error: " + string(s) + " ; err: " + err.Error())
	}

	for d := range msgs {
		//消息逻辑处理，可以自行设计逻辑
		hlog.Info("Received a message: " + string(d.Body))
		res := f(string(d.Body))

		if r.needAck {
			if res {
				//确认消息
				d.Ack(false)
			} else {
				//拒绝消息
				d.Reject(false)
			}
		}
	}

	r.ChannelPut(r.QueueName, ch)
}
