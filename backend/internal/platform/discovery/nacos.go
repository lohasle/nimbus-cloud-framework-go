package discovery

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Registry struct {
	client naming_client.INamingClient
	ip     string
	group  string
}

func New() (*Registry, error) {
	host := os.Getenv("NIMBUS_NACOS_HOST")
	if host == "" {
		return &Registry{}, nil
	}
	port, _ := strconv.ParseUint(env("NIMBUS_NACOS_PORT", "28858"), 10, 64)
	servers := []constant.ServerConfig{*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos"))}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(env("NIMBUS_NACOS_NAMESPACE", "public")),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir(".run/nacos/cache"),
		constant.WithLogDir(".run/nacos/log"),
		constant.WithLogLevel("warn"),
	)
	client, err := clients.NewNamingClient(vo.NacosClientParam{ClientConfig: &clientConfig, ServerConfigs: servers})
	if err != nil {
		return nil, fmt.Errorf("create nacos client: %w", err)
	}
	return &Registry{client: client, ip: env("NIMBUS_ADVERTISE_IP", "127.0.0.1"), group: "NIMBUS_GROUP"}, nil
}

func (r *Registry) Enabled() bool { return r != nil && r.client != nil }

func (r *Registry) Register(service string, port uint64) error {
	if !r.Enabled() {
		slog.Info("nacos disabled; using static local routing", "service", service)
		return nil
	}
	ok, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip: r.ip, Port: port, Weight: 1, Enable: true, Healthy: true, Ephemeral: true,
		ServiceName: service, GroupName: r.group, ClusterName: "DEFAULT", Metadata: map[string]string{"framework": "nimbus-go"},
	})
	if err != nil || !ok {
		return fmt.Errorf("register %s in nacos: %w", service, err)
	}
	slog.Info("service registered", "registry", "nacos", "service", service, "ip", r.ip, "port", port)
	return nil
}

func (r *Registry) Deregister(service string, port uint64) {
	if !r.Enabled() {
		return
	}
	_, _ = r.client.DeregisterInstance(vo.DeregisterInstanceParam{Ip: r.ip, Port: port, ServiceName: service, GroupName: r.group, Cluster: "DEFAULT", Ephemeral: true})
}

func (r *Registry) Resolve(service, fallback string) string {
	if !r.Enabled() {
		return fallback
	}
	instance, err := r.client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{ServiceName: service, GroupName: r.group})
	if err != nil || instance == nil {
		slog.Warn("nacos resolution failed; using fallback", "service", service, "error", err)
		return fallback
	}
	return fmt.Sprintf("http://%s:%d", instance.Ip, instance.Port)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
