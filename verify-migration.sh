#!/bin/bash

# Migration Verification Script
# This script checks if the backend API migration is complete

echo "=== CanIFly API Migration Verification ==="
echo

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track if any issues are found
ISSUES_FOUND=0

# Check for legacy endpoints in backend
echo "1. Checking for legacy endpoints in router.go..."
LEGACY_ENDPOINTS=$(grep -E "(app-data|toggle-account|update-account-name|remove-account|update-character|remove-character|choose-settings|save-user-selections|get-skill-plan|save-skill-plan|delete-skill-plan|associate-character|unassociate-character)" internal/server/router.go 2>/dev/null)
if [ -n "$LEGACY_ENDPOINTS" ]; then
    echo -e "${RED}❌ Found legacy endpoints in router.go:${NC}"
    echo "$LEGACY_ENDPOINTS"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No legacy endpoints found in router.go${NC}"
fi
echo

# Check for legacy API calls in frontend
echo "2. Checking for legacy API calls in frontend..."
LEGACY_API_CALLS=$(grep -r -E "(app-data|toggle-account|update-account-name|remove-account|update-character|remove-character|choose-settings|save-user-selections)" renderer/src/api/ --include="*.js" --include="*.jsx" 2>/dev/null | grep -v "test\.")
if [ -n "$LEGACY_API_CALLS" ]; then
    echo -e "${RED}❌ Found legacy API calls in frontend:${NC}"
    echo "$LEGACY_API_CALLS"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No legacy API calls found in frontend${NC}"
fi
echo

# Check for legacy handler methods
echo "3. Checking for legacy handler methods..."
LEGACY_HANDLERS=$(grep -r "func.*\(ToggleAccountStatus\|ToggleAccountVisibility\|UpdateAccountName\|RemoveAccount\|UpdateCharacter[^R]\|RemoveCharacter\|GetDashboardData\|GetDashboardDataNoCache\|SaveUserSelections\|ChooseSettingsDir\|ResetToDefaultDir\|AssociateCharacter\|UnassociateCharacter\|SaveSkillPlan\|GetSkillPlanFile\|DeleteSkillPlan[^R]\).*http\.HandlerFunc" internal/handlers/ --include="*.go" 2>/dev/null | grep -v "test\.")
if [ -n "$LEGACY_HANDLERS" ]; then
    echo -e "${RED}❌ Found legacy handler methods:${NC}"
    echo "$LEGACY_HANDLERS"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No legacy handler methods found${NC}"
fi
echo

# Check for getDashboards in frontend (should not exist)
echo "4. Checking for getDashboards calls..."
GET_DASHBOARDS=$(grep -r "getDashboards" renderer/src/ --include="*.js" --include="*.jsx" 2>/dev/null | grep -v "test\.")
if [ -n "$GET_DASHBOARDS" ]; then
    echo -e "${RED}❌ Found getDashboards calls (dashboards don't exist as entities):${NC}"
    echo "$GET_DASHBOARDS"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No getDashboards calls found${NC}"
fi
echo

# Check for required RESTful endpoints
echo "5. Checking for required RESTful endpoints..."
REQUIRED_PATTERNS=(
    "HandleFunc.*\/api\/accounts.*GET"
    "HandleFunc.*\/api\/accounts.*PATCH"
    "HandleFunc.*\/api\/accounts.*DELETE"
    "HandleFunc.*\/api\/characters.*GET"
    "HandleFunc.*\/api\/characters.*PATCH"
    "HandleFunc.*\/api\/characters.*DELETE"
    "HandleFunc.*\/api\/config.*GET"
    "HandleFunc.*\/api\/config.*PATCH"
)

MISSING_ENDPOINTS=0
for pattern in "${REQUIRED_PATTERNS[@]}"; do
    if ! grep -q "$pattern" internal/server/router.go 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Missing RESTful endpoint pattern: $pattern${NC}"
        MISSING_ENDPOINTS=1
    fi
done

if [ $MISSING_ENDPOINTS -eq 0 ]; then
    echo -e "${GREEN}✅ All required RESTful endpoints found${NC}"
else
    ISSUES_FOUND=1
fi
echo

# Summary
echo "=== Summary ==="
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ Migration appears to be complete!${NC}"
    echo "All legacy code has been removed and RESTful endpoints are in place."
else
    echo -e "${RED}❌ Migration is NOT complete!${NC}"
    echo "Please address the issues found above before considering the migration done."
fi
echo

# Additional checks that require manual verification
echo "=== Manual Verification Required ==="
echo "Please also verify:"
echo "1. All frontend components use the new API patterns"
echo "2. All tests have been updated for new endpoints"
echo "3. The application works end-to-end without errors"
echo "4. No 404 errors appear in the browser console"
echo "5. All user workflows function correctly"

exit $ISSUES_FOUND