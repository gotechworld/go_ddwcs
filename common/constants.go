package common

const GLOBAL_CONFIG_DEFAULT_SCOPE = 0
const ALTEX_WEBSITE_ID = 1
const MEDIAGALAXY_WEBSITE_ID = 2

const ALTEX_WEBSITE_CODE = "altex"
const MEDIAGALAXY_WEBSITE_CODE = "mediagalaxy"

// Get website id by website code
func GetWebsiteId(websiteCode string) int {
	switch {
	case ALTEX_WEBSITE_CODE == websiteCode:
		return ALTEX_WEBSITE_ID
	case MEDIAGALAXY_WEBSITE_CODE == websiteCode:
		return MEDIAGALAXY_WEBSITE_ID
	default:
		return ALTEX_WEBSITE_ID
	}
}

func IsValidWebsiteCode(websiteCode string) bool {
	switch {
	case ALTEX_WEBSITE_CODE == websiteCode:
		return true
	case MEDIAGALAXY_WEBSITE_CODE == websiteCode:
		return true
	default:
		return false
	}
}