package tools

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// ParseCertificate 解析PEM格式的证书
func ParseCertificate(pemCert string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemCert))
	if block == nil {
		return nil, fmt.Errorf("无法解析PEM块")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("证书解析失败: %w", err)
	}

	return cert, nil
}

// ParsePrivateKey 解析PEM格式的私钥
func ParsePrivateKey(pemKey string) (any, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, fmt.Errorf("无法解析PEM块")
	}

	// 尝试解析RSA私钥
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	// 尝试解析EC私钥
	ecKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err == nil {
		return ecKey, nil
	}

	return nil, fmt.Errorf("不支持的私钥格式")
}

// ExtractCNFromDN 从DN字符串中提取CN
func ExtractCNFromDN(dn string) string {
	parts := strings.Split(dn, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 && strings.TrimSpace(kv[0]) == "CN" {
			return strings.TrimSpace(kv[1])
		}
	}
	return ""
}

// GetFirstOrganization 从Organization数组中获取第一个组织
func GetFirstOrganization(orgs []string) string {
	if len(orgs) > 0 {
		return orgs[0]
	}
	return ""
}
