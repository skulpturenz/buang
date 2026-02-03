package main

import (
	"fmt"

	"github.com/dogmatiq/ferrite"
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	COMPUTE_INSTANCE_NAME = ferrite.
				String("COMPUTE_INSTANCE_NAME", "Unique name for compute instance").
				Required()
	GOOGLE_PROJECT = ferrite.
			String("GOOGLE_PROJECT", "GCP Project").
			Required()
	GOOGLE_SERVICE_ACCOUNT = ferrite.
				String("GOOGLE_SERVICE_ACCOUNT", "Service account linked to vm").
				Required()
	CLOUDFLARE_ZONE_ID = ferrite.
				String("CLOUDFLARE_ZONE_ID", "Cloudflare zone id").
				Required()
	CLOUDFLARE_API_TOKEN = ferrite.
				String("CLOUDFLARE_API_TOKEN", "Cloudflare API token").
				Required()
)

func main() {
	setupNetworking := func(ctx *pulumi.Context) error {
		buangNetwork, err := compute.NewNetwork(ctx, "buang-network", &compute.NetworkArgs{
			Name:        pulumi.String("buang-network"),
			Description: pulumi.String("Allows Cloudflare sources only"),
		})
		if err != nil {
			return err
		}

		ctx.Export("sharedResourceNetwork", buangNetwork.Name)

		_, err = compute.NewFirewall(ctx, "allow-cloudflare", &compute.FirewallArgs{
			Name:        pulumi.String("allow-cloudflare"),
			Network:     buangNetwork.Name,
			Description: pulumi.StringPtr("Allow all traffic from Cloudflare"),
			Allows: compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("0-65535"),
					},
				},
			},
			SourceRanges: pulumi.ToStringArray([]string{
				"173.245.48.0/20",
				"103.21.244.0/22",
				"103.22.200.0/22",
				"103.31.4.0/22",
				"141.101.64.0/18",
				"108.162.192.0/18",
				"190.93.240.0/20",
				"188.114.96.0/20",
				"197.234.240.0/22",
				"198.41.128.0/17",
				"162.158.0.0/15",
				"104.16.0.0/13",
				"104.24.0.0/14",
				"172.64.0.0/13",
				"131.0.72.0/22",
			},
			),
			TargetTags: pulumi.StringArray{
				pulumi.String("allow-cloudflare"),
			},
		})
		if err != nil {
			return err
		}

		_, err = compute.NewFirewall(ctx, "allow-icmp", &compute.FirewallArgs{
			Name:        pulumi.String("allow-icmp"),
			Network:     buangNetwork.Name,
			Description: pulumi.StringPtr("Allow ICMP"),
			Allows: compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("icmp"),
				},
			},
			SourceRanges: pulumi.ToStringArray([]string{
				"0.0.0.0",
			},
			),
			TargetTags: pulumi.StringArray{
				pulumi.String("allow-icmp"),
			},
		})
		if err != nil {
			return err
		}

		_, err = compute.NewFirewall(ctx, "allow-ssh", &compute.FirewallArgs{
			Name:        pulumi.String("allow-ssh"),
			Network:     buangNetwork.Name,
			Description: pulumi.StringPtr("Allow SSH"),
			Allows: compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("22"),
					},
				},
			},
			SourceRanges: pulumi.ToStringArray([]string{
				"0.0.0.0/0",
			},
			),
			TargetTags: pulumi.StringArray{
				pulumi.String("allow-ssh"),
			},
		})
		if err != nil {
			return err
		}

		_, err = compute.NewFirewall(ctx, "allow-health-check", &compute.FirewallArgs{
			Name:        pulumi.String("allow-health-check"),
			Network:     buangNetwork.Name,
			Description: pulumi.StringPtr("Allow GCP health check"),
			Allows: compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("80"),
						pulumi.String("443"),
					},
				},
			},
			SourceRanges: pulumi.ToStringArray([]string{
				"130.211.0.0/22",
				"35.191.0.0/16",
			},
			),
			TargetTags: pulumi.StringArray{
				pulumi.String("allow-health-check"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

	type devResources struct {
		region               string
		zone                 string
		instanceTemplate     *compute.InstanceTemplate
		instanceGroupManager *compute.InstanceGroupManager
	}
	setupDev := func(ctx *pulumi.Context) (*devResources, error) {
		REGION := "us-central1"
		ZONE := "us-central1-a"

		static, err := compute.NewAddress(ctx, fmt.Sprintf("%s-dev", COMPUTE_INSTANCE_NAME.Value()), &compute.AddressArgs{
			Name:   pulumi.Sprintf("%s-dev", COMPUTE_INSTANCE_NAME.Value()),
			Region: pulumi.String(REGION),
		})
		if err != nil {
			return nil, err
		}

		instanceTemplate, err := compute.NewInstanceTemplate(ctx, fmt.Sprintf("%s-dev-template", COMPUTE_INSTANCE_NAME.Value()), &compute.InstanceTemplateArgs{
			// see: https://github.com/pulumi/pulumi-gcp/issues/680#issuecomment-1405680098
			NamePrefix:   pulumi.Sprintf("%s-dev-template", COMPUTE_INSTANCE_NAME.Value()),
			MachineType:  pulumi.String("e2-custom-micro-2048"),
			CanIpForward: pulumi.Bool(false),
			Tags: pulumi.ToStringArray([]string{
				"allow-cloudflare",
				"allow-ssh",
				"allow-health-check",
			}),
			NetworkInterfaces: &compute.InstanceTemplateNetworkInterfaceArray{
				&compute.InstanceTemplateNetworkInterfaceArgs{
					AccessConfigs: &compute.InstanceTemplateNetworkInterfaceAccessConfigArray{
						&compute.InstanceTemplateNetworkInterfaceAccessConfigArgs{
							NatIp: static.Address,
						},
					},
					Network: pulumi.String("buang-network"),
				},
			},
			Disks: compute.InstanceTemplateDiskArray{
				&compute.InstanceTemplateDiskArgs{
					// we want to reuse the same disk if an instance is replaced
					DiskName:    pulumi.Sprintf("%s-dev-disk", COMPUTE_INSTANCE_NAME.Value()),
					SourceImage: pulumi.String("debian-12-bookworm-v20240515"),
					AutoDelete:  pulumi.Bool(false),
					Boot:        pulumi.Bool(true),
					DiskSizeGb:  pulumi.Int(50),
				},
			},
			Scheduling: &compute.InstanceTemplateSchedulingArgs{
				Preemptible:       pulumi.Bool(true),
				AutomaticRestart:  pulumi.Bool(false),
				ProvisioningModel: pulumi.String("SPOT"),
				OnHostMaintenance: pulumi.String("TERMINATE"),
			},
			// Docker setup on Debian 12: https://www.thomas-krenn.com/en/wiki/Docker_installation_on_Debian_12
			// Permanently increase vm.max_map_count value: https://thetechdarts.com/how-to-change-default-vm-max_map_count-on-linux/
			// Enable root login: https://cloud.google.com/compute/docs/connect/root-ssh
			MetadataStartupScript: pulumi.Sprintf(`#! /bin/bash
				curl -sSO https://dl.google.com/cloudagents/add-google-cloud-ops-agent-repo.sh
				sudo bash add-google-cloud-ops-agent-repo.sh --also-install

				sudo apt update &&
				sudo apt install certbot python3-certbot-dns-cloudflare make git ca-certificates curl gnupg apt-transport-https gpg -y &&
				curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker.gpg &&
				echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker.gpg] https://download.docker.com/linux/debian bookworm stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null &&
				sudo apt update &&
				sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-compose -y &&
				sudo grep -qxF 'vm.max_map_count=262144' /etc/sysctl.conf || echo vm.max_map_count=262144 | sudo tee -a /etc/sysctl.conf &&
				sudo sysctl -p &&
				sudo mkdir -p /etc/letsencrypt/renewal-hooks/deploy &&
				echo "dns_cloudflare_api_token = %s" | sudo tee /etc/letsencrypt/dnscloudflare.ini &&
				echo "#! /bin/bash sudo docker service ls -q | xargs -n1 sudo docker service update --force" | sudo tee /etc/letsencrypt/renewal-hooks/deploy/reload-services.sh &&
				sudo chmod 0600 /etc/letsencrypt/dnscloudflare.ini &&
				sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/reload-services.sh &&
				sudo chmod 0600 /etc/letsencrypt/renewal-hooks/deploy/reload-services.sh &&
				sudo certbot certonly -d api-dev.buang.skulpture.xyz,api-prod.buang.skulpture.xyz,buang.skulpture.xyz \
					--dns-cloudflare --dns-cloudflare-credentials /etc/letsencrypt/dnscloudflare.ini \
					--non-interactive --agree-tos \
					--register-unsafely-without-email \
					--dns-cloudflare-propagation-seconds 60

				# Enable root login and restart sshd
				sudo sed -i 's/PermitRootLogin no/PermitRootLogin prohibit-password/g' /etc/ssh/sshd_config
				sudo systemctl restart sshd`, CLOUDFLARE_API_TOKEN.Value()),
			ServiceAccount: &compute.InstanceTemplateServiceAccountArgs{
				Email: pulumi.StringPtr(GOOGLE_SERVICE_ACCOUNT.Value()),
				Scopes: pulumi.ToStringArray([]string{
					"cloud-platform",
				}),
			},
		})
		if err != nil {
			return nil, err
		}

		autohealing, err := compute.NewHealthCheck(ctx, fmt.Sprintf("%s-dev-autohealing", COMPUTE_INSTANCE_NAME.Value()), &compute.HealthCheckArgs{
			Name:               pulumi.Sprintf("%s-dev-autohealing", COMPUTE_INSTANCE_NAME.Value()),
			CheckIntervalSec:   pulumi.Int(5),
			TimeoutSec:         pulumi.Int(5),
			HealthyThreshold:   pulumi.Int(2),
			UnhealthyThreshold: pulumi.Int(10),
			HttpsHealthCheck: &compute.HealthCheckHttpsHealthCheckArgs{
				RequestPath: pulumi.String("/ping"),
				Host:        pulumi.String("api-dev.buang.skulpture.xyz"),
			},
		})
		if err != nil {
			return nil, err
		}

		instanceGroupManager, err := compute.NewInstanceGroupManager(ctx, fmt.Sprintf("%s-dev-igm", COMPUTE_INSTANCE_NAME.Value()), &compute.InstanceGroupManagerArgs{
			Name:             pulumi.String(fmt.Sprintf("%s-dev-igm", COMPUTE_INSTANCE_NAME.Value())),
			BaseInstanceName: pulumi.String(fmt.Sprintf("%s-dev-instance", COMPUTE_INSTANCE_NAME.Value())),
			Zone:             pulumi.String(ZONE),
			TargetSize:       pulumi.Int(1),
			Versions: compute.InstanceGroupManagerVersionArray{
				&compute.InstanceGroupManagerVersionArgs{
					InstanceTemplate: instanceTemplate.SelfLinkUnique,
					Name:             pulumi.String("primary"),
				},
			},
			StandbyPolicy: &compute.InstanceGroupManagerStandbyPolicyArgs{
				Mode: pulumi.String("MANUAL"),
			},
			UpdatePolicy: &compute.InstanceGroupManagerUpdatePolicyArgs{
				MinimalAction:       pulumi.String("REPLACE"),
				Type:                pulumi.String("PROACTIVE"),
				MaxSurgeFixed:       pulumi.Int(0),
				MaxUnavailableFixed: pulumi.Int(1),
				ReplacementMethod:   pulumi.String("RECREATE"),
			},
			NamedPorts: compute.InstanceGroupManagerNamedPortArray{
				&compute.InstanceGroupManagerNamedPortArgs{
					Name: pulumi.String("https"),
					Port: pulumi.Int(443),
				},
			},
			AutoHealingPolicies: &compute.InstanceGroupManagerAutoHealingPoliciesArgs{
				HealthCheck:     autohealing.ID(),
				InitialDelaySec: pulumi.Int(300),
			},
		})
		if err != nil {
			return nil, err
		}

		// EXPENSIVE
		// https://cloud.google.com/vpc/network-pricing?hl=en#section-2

		// devManagedSslCertificate, err := compute.NewManagedSslCertificate(ctx, fmt.Sprintf("%s-dev-managed-ssl-cert", COMPUTE_INSTANCE_NAME.Value()), &compute.ManagedSslCertificateArgs{
		// 	Name: pulumi.Sprintf("%s-dev-managed-ssl-cert", COMPUTE_INSTANCE_NAME.Value()),
		// 	Managed: &compute.ManagedSslCertificateManagedArgs{
		// 		Domains: pulumi.StringArray{
		// 			pulumi.String("buang-api-dev.skulpture.xyz"),
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// devHttpsHealthCheck, err := compute.NewHttpsHealthCheck(ctx, fmt.Sprintf("%s-dev-backend-healthcheck", COMPUTE_INSTANCE_NAME.Value()), &compute.HttpsHealthCheckArgs{
		// 	Name:             pulumi.Sprintf("%s-dev-backend-healthcheck", COMPUTE_INSTANCE_NAME.Value()),
		// 	RequestPath:      pulumi.String("/ping"),
		// 	CheckIntervalSec: pulumi.Int(5),
		// 	TimeoutSec:       pulumi.Int(5),
		// 	Host:             pulumi.String("buang-api-dev.skulpture.xyz"),
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// devBackendService, err := compute.NewBackendService(ctx, fmt.Sprintf("%s-dev-backend", COMPUTE_INSTANCE_NAME.Value()), &compute.BackendServiceArgs{
		// 	Name:         pulumi.Sprintf("%s-dev-backend", COMPUTE_INSTANCE_NAME.Value()),
		// 	Protocol:     pulumi.String("HTTPS"),
		// 	PortName:     pulumi.String("https"),
		// 	HealthChecks: devHttpsHealthCheck.ID(),
		// 	Backends: compute.BackendServiceBackendArray{
		// 		compute.BackendServiceBackendArgs{
		// 			Group: instanceGroupManager.InstanceGroup,
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// devUrlMap, err := compute.NewURLMap(ctx, fmt.Sprintf("%s-dev-url-map", COMPUTE_INSTANCE_NAME.Value()), &compute.URLMapArgs{
		// 	Name:           pulumi.Sprintf("%s-dev-url-map", COMPUTE_INSTANCE_NAME.Value()),
		// 	DefaultService: devBackendService.ID(),
		// 	HostRules: compute.URLMapHostRuleArray{
		// 		&compute.URLMapHostRuleArgs{
		// 			Hosts: pulumi.StringArray{
		// 				pulumi.String("*"),
		// 			},
		// 			PathMatcher: pulumi.String("allpaths"),
		// 		},
		// 	},
		// 	PathMatchers: compute.URLMapPathMatcherArray{
		// 		&compute.URLMapPathMatcherArgs{
		// 			Name:           pulumi.String("allpaths"),
		// 			DefaultService: devBackendService.ID(),
		// 			PathRules: compute.URLMapPathMatcherPathRuleArray{
		// 				&compute.URLMapPathMatcherPathRuleArgs{
		// 					Paths: pulumi.StringArray{
		// 						pulumi.String("/*"),
		// 					},
		// 					Service: devBackendService.ID(),
		// 				},
		// 			},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// devHttpsProxy, err := compute.NewTargetHttpsProxy(ctx, fmt.Sprintf("%s-dev-proxy", COMPUTE_INSTANCE_NAME.Value()), &compute.TargetHttpsProxyArgs{
		// 	Name:   pulumi.Sprintf("%s-dev-proxy", COMPUTE_INSTANCE_NAME.Value()),
		// 	UrlMap: devUrlMap.ID(),
		// 	SslCertificates: pulumi.StringArray{
		// 		devManagedSslCertificate.ID(),
		// 	},
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// devLoadBalancer, err := compute.NewGlobalForwardingRule(ctx, fmt.Sprintf("%s-dev-lb", COMPUTE_INSTANCE_NAME.Value()), &compute.GlobalForwardingRuleArgs{
		// 	Name:      pulumi.Sprintf("%s-dev-lb", COMPUTE_INSTANCE_NAME.Value()),
		// 	Target:    devHttpsProxy.ID(),
		// 	PortRange: pulumi.String("443"),
		// })
		// if err != nil {
		// 	return nil, err
		// }

		ctx.Export("devStaticAddress", static.Address)
		// ctx.Export("devGlobalForwardingRuleAddress", devLoadBalancer.IpAddress)

		_, err = cloudflare.NewRecord(ctx, fmt.Sprintf("%s-api-dev", COMPUTE_INSTANCE_NAME.Value()), &cloudflare.RecordArgs{
			ZoneId:  pulumi.String(CLOUDFLARE_ZONE_ID.Value()),
			Name:    pulumi.String("buang"),
			Content: static.Address, // devLoadBalancer.IpAddress,
			Type:    pulumi.String("A"),
			Proxied: pulumi.Bool(true),
		})
		if err != nil {
			return nil, err
		}

		result := devResources{
			region:               REGION,
			zone:                 ZONE,
			instanceTemplate:     instanceTemplate,
			instanceGroupManager: instanceGroupManager,
		}

		return &result, nil
	}

	setupIdentityPool := func(ctx *pulumi.Context, devResources *devResources) error {
		version, err := random.NewRandomInteger(ctx, "wif-pool-count", &random.RandomIntegerArgs{
			Min: pulumi.Int(1),
			Max: pulumi.Int(1000),
			Keepers: pulumi.StringMap{
				"version": pulumi.String("1"),
			},
		})
		if err != nil {
			return err
		}

		pool, err := iam.NewWorkloadIdentityPool(ctx, "buang-wif-pool", &iam.WorkloadIdentityPoolArgs{
			// workload identity pool ids must be unique and they are soft deleted for 30 days
			// so if the stack is torn down and brought up again then it results in an error about the id being used
			// whenever we want to modify the instance group template or what not the stack can be torn down entirely
			WorkloadIdentityPoolId: pulumi.Sprintf("buang-wif-pool-%d", version.Result),
		})
		if err != nil {
			return err
		}

		ctx.Export("devWifPool", pool.Name)

		githubProvider, err := iam.NewWorkloadIdentityPoolProvider(ctx, "buang-wif-provider-gh", &iam.WorkloadIdentityPoolProviderArgs{
			WorkloadIdentityPoolId:         pool.WorkloadIdentityPoolId,
			WorkloadIdentityPoolProviderId: pulumi.String("github"),
			DisplayName:                    pulumi.String("Github"),
			AttributeMapping: pulumi.StringMap{
				"google.subject":       pulumi.String("assertion.sub"),
				"attribute.actor":      pulumi.String("assertion.actor"),
				"attribute.repository": pulumi.String("assertion.repository"),
				"attribute.ref":        pulumi.String("assertion.ref"),
			},
			Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
				IssuerUri: pulumi.String("https://token.actions.githubusercontent.com"),
			},
			AttributeCondition: pulumi.String("attribute.repository==assertion.repository"),
		})
		if err != nil {
			return err
		}

		ctx.Export("devGithubWIFProvider", githubProvider.Name)

		const REPOSITORY = "skulpturenz/buang"

		devResources.instanceTemplate.Name.ApplyT(func(instanceTemplateName string) error {
			pool.Name.ApplyT(func(workloadIdentityPoolId string) error {
				principalSet := fmt.Sprintf("principalSet://iam.googleapis.com/%s/attribute.repository/%s", workloadIdentityPoolId, REPOSITORY)

				_, err = compute.NewInstanceTemplateIamBinding(ctx, fmt.Sprintf("%s-compute-admin-dev", instanceTemplateName), &compute.InstanceTemplateIamBindingArgs{
					Name:    pulumi.String(instanceTemplateName),
					Role:    pulumi.String("roles/compute.admin"),
					Members: pulumi.ToStringArray([]string{principalSet}),
				})
				if err != nil {
					return err
				}

				return nil
			})

			return nil
		})

		return nil
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		err := setupNetworking(ctx)
		if err != nil {
			return err
		}

		devResources, err := setupDev(ctx)
		if err != nil {
			return err
		}

		err = setupIdentityPool(ctx, devResources)
		if err != nil {
			return err
		}

		return nil
	})
}
