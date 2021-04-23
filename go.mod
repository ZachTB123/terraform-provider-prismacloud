module github.com/terraform-providers/terraform-provider-prismacloud

go 1.13

require (
	github.com/hashicorp/terraform-plugin-sdk v1.9.0
	github.com/paloaltonetworks/prisma-cloud-go v0.4.0
)

//replace github.com/paloaltonetworks/prisma-cloud-go => ../prisma-cloud-go

replace github.com/paloaltonetworks/prisma-cloud-go v0.4.0 => github.com/ZachTB123/prisma-cloud-go v0.4.2
