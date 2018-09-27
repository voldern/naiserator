package fixtures

import (
	nais "github.com/nais/naiserator/pkg/apis/naiserator/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Name                      = "app"
	Namespace                 = "default"
	Port                      = 8080
	ImageName                 = "user/image:version"
	TeamName                  = "pandas"
	MinReplicas               = 1
	MaxReplicas               = 2
	CpuThresholdPercentage    = 69
	ReadinessPath             = "isready"
	ReadinessInitialDelay     = 1
	ReadinessTimeout          = 2
	ReadinessFailureThreshold = 3
	ReadinessPeriodSeconds    = 4
	LivenessPath              = "isalive"
	LivenessInitialDelay      = 5
	LivenessTimeout           = 6
	LivenessFailureThreshold  = 7
	LivenessPeriodSeconds     = 8
	RequestCpu                = "200m"
	RequestMemory             = "256Mi"
	LimitCpu                  = "500m"
	LimitMemory               = "512Mi"
	PrometheusPath            = "metrics"
	PrometheusPort            = "http"
	PrometheusEnabled         = true
	IstioEnabled              = true
	WebProxyEnabled           = true
	IngressDisabled           = true
	LeaderElectionEnabled     = true
	SecretsEnabled            = true
	PreStopHookPath           = "die"
	LogFormat                 = "accesslog"
	LogTransform              = "dns_loglevel"
)

func Application() *nais.Application {
	app := &nais.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: Namespace,
		},
		Spec: nais.ApplicationSpec{
			Port:  Port,
			Image: ImageName,
			Team:  TeamName,
			Replicas: nais.Replicas{
				Min:                    MinReplicas,
				Max:                    MaxReplicas,
				CpuThresholdPercentage: CpuThresholdPercentage,
			},
			Healthcheck: nais.Healthcheck{
				Readiness: nais.Probe{
					Path:             ReadinessPath,
					InitialDelay:     ReadinessInitialDelay,
					FailureThreshold: ReadinessFailureThreshold,
					Timeout:          ReadinessTimeout,
					PeriodSeconds:    ReadinessPeriodSeconds,
				},
				Liveness: nais.Probe{
					Path:             LivenessPath,
					InitialDelay:     LivenessInitialDelay,
					FailureThreshold: LivenessFailureThreshold,
					Timeout:          LivenessTimeout,
					PeriodSeconds:    LivenessPeriodSeconds,
				},
			},
			Resources: nais.ResourceRequirements{
				Requests: nais.ResourceSpec{
					Memory: RequestMemory,
					Cpu:    RequestCpu,
				},
				Limits: nais.ResourceSpec{
					Memory: LimitMemory,
					Cpu:    LimitCpu,
				},
			},
			Prometheus: nais.PrometheusConfig{
				Path:    PrometheusPath,
				Port:    PrometheusPort,
				Enabled: PrometheusEnabled,
			},
			Istio: nais.IstioConfig{
				Enabled: IstioEnabled,
			},
			Logtransform:    LogTransform,
			Logformat:       LogFormat,
			WebProxy:        WebProxyEnabled,
			PreStopHookPath: PreStopHookPath,
			Ingress: nais.Ingress{
				Disabled: IngressDisabled,
			},
			LeaderElection: LeaderElectionEnabled,
			Secrets:        SecretsEnabled,
		}}

	return app
}
