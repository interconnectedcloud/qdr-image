package router

// QdpidDispatch representation
type QpidDispatch struct {
	Id                    string         `json:"id,omitempty"`
	Role                  RouterRoleType `json:"role,omitempty"`
	Users                 string         `json:"users,omitempty"`
	Listeners             []Listener     `json:"listeners,omitempty"`
	InterRouterListeners  []Listener     `json:"interRouterListeners,omitempty"`
	EdgeListeners         []Listener     `json:"edgeListeners,omitempty"`
	SslProfiles           []SslProfile   `json:"sslProfiles,omitempty"`
	Addresses             []Address      `json:"addresses,omitempty"`
	AutoLinks             []AutoLink     `json:"autoLinks,omitempty"`
	LinkRoutes            []LinkRoute    `json:"linkRoutes,omitempty"`
	Connectors            []Connector    `json:"connectors,omitempty"`
	InterRouterConnectors []Connector    `json:"interRouterConnectors,omitempty"`
	EdgeConnectors        []Connector    `json:"edgeConnectors,omitempty"`
	TcpConnectors         []TcpConnector `json:"tcpConnectors,omitempty"`
	TcpListeners          []TcpListener  `json:"tcpListeners,omitempty"`
	Logs                  []Log          `json:"logs,omitempty"`
}

type RouterRoleType string

const (
	RouterRoleInterior RouterRoleType = "interior"
	RouterRoleEdge                    = "edge"
)

type Address struct {
	Prefix         string `json:"prefix,omitempty"`
	Pattern        string `json:"pattern,omitempty"`
	Distribution   string `json:"distribution,omitempty"`
	Waypoint       bool   `json:"waypoint,omitempty"`
	IngressPhase   *int32 `json:"ingressPhase,omitempty"`
	EgressPhase    *int32 `json:"egressPhase,omitempty"`
	Priority       *int32 `json:"priority,omitempty"`
	EnableFallback bool   `json:"enableFallback,omitempty"`
}

type Listener struct {
	Name             string `json:"name,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int32  `json:"port"`
	RouteContainer   bool   `json:"routeContainer,omitempty"`
	Http             bool   `json:"http,omitempty"`
	Cost             int32  `json:"cost,omitempty"`
	SslProfile       string `json:"sslProfile,omitempty"`
	SaslMechanisms   string `json:"saslMechanisms,omitempty"`
	AuthenticatePeer bool   `json:"authenticatePeer,omitempty"`
	Expose           bool   `json:"expose,omitempty"`
	LinkCapacity     int32  `json:"linkCapacity,omitempty"`
}

type SslProfile struct {
	Name                string `json:"name,omitempty"`
	Credentials         string `json:"credentials,omitempty"`
	CaCert              string `json:"caCert,omitempty"`
	GenerateCredentials bool   `json:"generateCredentials,omitempty"`
	GenerateCaCert      bool   `json:"generateCaCert,omitempty"`
	MutualAuth          bool   `json:"mutualAuth,omitempty"`
	Ciphers             string `json:"ciphers,omitempty"`
	Protocols           string `json:"protocols,omitempty"`
}

type LinkRoute struct {
	Prefix            string `json:"prefix,omitempty"`
	Pattern           string `json:"pattern,omitempty"`
	Direction         string `json:"direction,omitempty"`
	ContainerId       string `json:"containerId,omitempty"`
	Connection        string `json:"connection,omitempty"`
	AddExternalPrefix string `json:"addExternalPrefix,omitempty"`
	DelExternalPrefix string `json:"delExternalPrefix,omitempty"`
}

type Connector struct {
	Name           string `json:"name,omitempty"`
	Host           string `json:"host"`
	Port           int32  `json:"port"`
	RouteContainer bool   `json:"routeContainer,omitempty"`
	Cost           int32  `json:"cost,omitempty"`
	VerifyHostname bool   `json:"verifyHostname,omitempty"`
	SslProfile     string `json:"sslProfile,omitempty"`
	LinkCapacity   int32  `json:"linkCapacity,omitempty"`
}

type AutoLink struct {
	Address         string `json:"address"`
	Direction       string `json:"direction"`
	ContainerId     string `json:"containerId,omitempty"`
	Connection      string `json:"connection,omitempty"`
	ExternalAddress string `json:"externalAddress,omitempty"`
	Phase           *int32 `json:"phase,omitempty"`
	Fallback        bool   `json:"fallback,omitempty"`
}

type TcpConnector struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	Address string `json:"address"`
	SiteId  string `json:"siteId,omitempty"`
}

type TcpListener struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	Address string `json:"address"`
	SiteId  string `json:"siteId,omitempty"`
}

type Log struct {
	Module           string `json:"module"`
	Enable           string `json:"enable"`
	IncludeTimestamp bool   `json:"includeTimestamp,omitempty"`
	IncludeSource    bool   `json:"includeSource,omitempty"`
	OutputFile       string `json:"outputFile,omitempty"`
}
