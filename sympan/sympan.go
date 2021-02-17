package sympan

import (
	"errors"
	"strings"

	"github.com/theovassiliou/go-eureka-client/eureka"
)

func ShortenFQName(fqName string) string {
	fqName = strings.TrimSuffix(fqName, ".")
	elements := strings.Split(fqName, ".")
	fqName = strings.Join(elements[:len(elements)-1], ".")
	fqName = strings.TrimSuffix(fqName, ".")
	return fqName
}

// BuildFQWormhole contructs from a given fully-qualified star name
// the fully qualliffied wormhole name
func BuildFQWormhole(fqName string) string {
	s := ShortenFQName(fqName)
	if s == "" {
		return ""
	}
	return s + ".WH"
}

// ResolveApplication returns a callable instance info if the application
// can be resolved.
// resolver will be used for seeking the applications or wormholes
// fqServiceName contains the fully qualified application or service name to look for
// proto if not "", should contain either "grpc" or "http"
// includeWormholes indicates whether also matching wormhole should be considered as qualified application
func ResolveApplication(resolver *eureka.Client, fqServiceName, proto string, includeWormholes bool) (eureka.InstanceInfo, error) {
	theSelectedInstance := eureka.InstanceInfo{}
	cont := true

	for cont {
		app, err := resolver.GetApplication(fqServiceName)

		if err != nil && err.Error() != "EOF" {
			return theSelectedInstance, err
		}

		if app != nil && len(app.Instances) > 0 {
			theSelectedInstance = selectOneOf(app.Instances, proto)
			return theSelectedInstance, nil
		}

		fqServiceName = ShortenFQName(strings.TrimSuffix(fqServiceName, ".WH"))
		if len(fqServiceName) > 0 {
			fqServiceName = fqServiceName + ".WH"
		} else {
			cont = false
		}
	}
	return theSelectedInstance, errors.New("not found")
}

// WormholeResolveApplication returns a callable instanceInfo if the application
// can be resolved. The wormhole first looks whether it can directly call the service
func WormholeResolveApplication(resolver *eureka.Client, scope, fqServiceName, proto string, includeWormholes bool) (eureka.InstanceInfo, error) {
	theSelectedInstance := eureka.InstanceInfo{}

	serviceName := strings.TrimPrefix(fqServiceName, scope+".")

	app, err := resolver.GetApplication(serviceName)
	if err != nil {
		return ResolveApplication(resolver, fqServiceName, proto, includeWormholes)
	}

	if app != nil && len(app.Instances) > 0 {
		theSelectedInstance = selectOneOf(app.Instances, proto)
		return theSelectedInstance, nil
	}

	return theSelectedInstance, errors.New("not found")
}

func selectOneOf(instances []eureka.InstanceInfo, proto string) eureka.InstanceInfo {

	// TODO: Find a better selection mechanism
	if len(instances) > 0 {
		for _, i := range instances {
			if strings.HasPrefix(i.HostName, proto) {
				return i
			}
		}
	}
	return eureka.InstanceInfo{}
}
