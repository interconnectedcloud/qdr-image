package router

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QpidDispatchConfigMap core.ConfigMap

func NewConfigMap(q QpidDispatch, namespace string, labels map[string]string) *QpidDispatchConfigMap {
	cm := &QpidDispatchConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      q.Id,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"qdrouterd.conf": ConfigForQpidDispatch(q),
		},
	}

	return cm
}

// ConfigForQpidDispatch returns a string representation of a
// serialized QpidDispatch configuration
func ConfigForQpidDispatch(q QpidDispatch) string {
	config := `
router {
    mode: {{.Role}}
    id: {{.Id}}
}
{{range .Listeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .RouteContainer}}
    role: route-container
    {{- else }}
    role: normal
    {{- end}}
    {{- if .Http}}
    http: true
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
{{range .InterRouterListeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: inter-router
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
{{range .EdgeListeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: edge
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
{{range .SslProfiles}}
sslProfile {
   name: {{.Name}}
   {{- if .Credentials}}
   certFile: /etc/router-certs/{{.Name}}/{{.Credentials}}/tls.crt
   privateKeyFile: /etc/router-certs/{{.Name}}/{{.Credentials}}/tls.key
   {{- end}}
   {{- if .CaCert}}
       {{- if eq .CaCert .Credentials}}
   caCertFile: /etc/router-certs/{{.Name}}/{{.CaCert}}/ca.crt
       {{- else if and .GenerateCredentials .MutualAuth}}
   caCertFile: /etc/router-certs/{{.Name}}/{{.Credentials}}/ca.crt
       {{- else}}
   caCertFile: /etc/router-certs/{{.Name}}/{{.CaCert}}/tls.crt
       {{- end}}
   {{- end}}
}
{{- end}}
{{range .Addresses}}
address {
    {{- if .Prefix}}
    prefix: {{.Prefix}}
    {{- end}}
    {{- if .Pattern}}
    pattern: {{.Pattern}}
    {{- end}}
    {{- if .Distribution}}
    distribution: {{.Distribution}}
    {{- end}}
    {{- if .Waypoint}}
    waypoint: {{.Waypoint}}
    {{- end}}
    {{- if .IngressPhase}}
    ingressPhase: {{.IngressPhase}}
    {{- end}}
    {{- if .EgressPhase}}
    egressPhase: {{.EgressPhase}}
    {{- end}}
    {{- if .Priority}}
    priority: {{.Priority}}
    {{- end}}
    {{- if .EnableFallback}}
    enableFallback: {{.EnableFallback}}
    {{- end}}
}
{{- end}}
{{range .AutoLinks}}
autoLink {
    {{- if .Address}}
    addr: {{.Address}}
    {{- end}}
    {{- if .Direction}}
    direction: {{.Direction}}
    {{- end}}
    {{- if .ContainerId}}
    containerId: {{.ContainerId}}
    {{- end}}
    {{- if .Connection}}
    connection: {{.Connection}}
    {{- end}}
    {{- if .ExternalAddress}}
    externalAddress: {{.ExternalAddress}}
    {{- end}}
    {{- if .Phase}}
    phase: {{.Phase}}
    {{- end}}
    {{- if .Fallback}}
    fallback: {{.Fallback}}
    {{- end}}
}
{{- end}}
{{range .LinkRoutes}}
linkRoute {
    {{- if .Prefix}}
    prefix: {{.Prefix}}
    {{- end}}
    {{- if .Pattern}}
    pattern: {{.Pattern}}
    {{- end}}
    {{- if .Direction}}
    direction: {{.Direction}}
    {{- end}}
    {{- if .Connection}}
    connection: {{.Connection}}
    {{- end}}
    {{- if .ContainerId}}
    containerId: {{.ContainerId}}
    {{- end}}
    {{- if .AddExternalPrefix}}
    addExternalPrefix: {{.AddExternalPrefix}}
    {{- end}}
    {{- if .DelExternalPrefix}}
    delExternalPrefix: {{.DelExternalPrefix}}
    {{- end}}
}
{{- end}}
{{range .Connectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .RouteContainer}}
    role: route-container
    {{- else}}
    role: normal
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}
{{range .InterRouterConnectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: inter-router
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}
{{range .EdgeConnectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: edge
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}
{{range .TcpConnectors}}
tcpConnector {
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Address}}
    address: {{.Address}}
    {{- end}}
    {{- if .SiteId}}
    siteId: {{.SiteId}}
    {{- end}}
}
{{- end}}
{{range .TcpListeners}}
tcpListener {
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Address}}
    address: {{.Address}}
    {{- end}}
    {{- if .SiteId}}
    siteId: {{.SiteId}}
    {{- end}}
}
{{- end}}
{{range .Logs}}
log {
    {{- if .Module}}
    module: {{.Module}}
    {{- end}}
    {{- if .Enable}}
    enable: {{.Enable}}
    {{- end}}
    {{- if .IncludeTimestamp}}
    includeTimestamp: {{.IncludeTimestamp}}
    {{- end}}
    {{- if .IncludeSource}}
    includeSource: {{.IncludeSource}}
    {{- end}}
    {{- if .OutputFile}}
    outputFile: {{.OutputFile}}
    {{- end}}
}
{{- end}}
`
	var buff bytes.Buffer
	qdrconfig := template.Must(template.New("qdrconfig").Parse(config))
	qdrconfig.Execute(&buff, q)
	return buff.String()
}
