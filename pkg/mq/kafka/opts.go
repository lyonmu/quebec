package kafka

import "github.com/IBM/sarama"

type Options struct {
	cfg        *sarama.Config
	partitions uint32
}

type Option func(*Options)

func WithSASLGSSAPI(service, krbConf, keytab, relam, username string) Option {
	return func(o *Options) {
		o.cfg.Net.SASL.Enable = true
		o.cfg.Net.SASL.Mechanism = sarama.SASLTypeGSSAPI
		o.cfg.Net.SASL.GSSAPI.ServiceName = service
		o.cfg.Net.SASL.GSSAPI.AuthType = sarama.KRB5_KEYTAB_AUTH
		o.cfg.Net.SASL.GSSAPI.KeyTabPath = keytab
		o.cfg.Net.SASL.GSSAPI.KerberosConfigPath = krbConf
		o.cfg.Net.SASL.GSSAPI.Realm = relam
	}
}

func WithSASLPlain(username, password string) Option {
	return func(o *Options) {
		o.cfg.Net.SASL.Enable = true
		o.cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		o.cfg.Net.SASL.User = username
		o.cfg.Net.SASL.Password = password
	}
}

func WithCompressionZSTD() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionZSTD
	}
}

func WithCompressionGZIP() Option {
	return func(o *Options) {
		o.cfg.Producer.Compression = sarama.CompressionGZIP
	}
}

func WithPartitions(partitons uint32) Option {
	return func(o *Options) {
		o.partitions = partitons
	}
}
