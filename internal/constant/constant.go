package constant

type UserType string

const (
	UserTypeInvestor UserType = "Investor"
	UserTypeAdmin    UserType = "Admin"
)

type InvestorType string

const (
	InvestorIndividual  InvestorType = "INDIVIDUAL"
	InvestorInstitution InvestorType = "INSTITUTION"
)

type UserRole string

const (
	RoleSuperadmin UserRole = "SUPERADMIN"
	RoleOperations UserRole = "OPERATIONS"
	RoleFinance    UserRole = "FINANCE"
	RoleBusiness   UserRole = "BUSINESS"
	RoleMarketing  UserRole = "MARKETING"
	RoleInvestor   UserRole = "INVESTOR"
	RoleAuditor    UserRole = "AUDITOR"
	RoleSupport    UserRole = "SUPPORT"
)

type SSOPlatform string

const (
	SSOPlatformGoogle   SSOPlatform = "GOOGLE"
	SSOPlatformApple    SSOPlatform = "APPLE"
	SSOPlatformFacebook SSOPlatform = "FACEBOOK"
)
