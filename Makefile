#-----------------------------------------#
# Commands
#-----------------------------------------#

GO ?= go
GORUN ?= $(GO) run
GOTEST ?= gotest
GOTOOL ?= $(GO) tool

#-----------------------------------------#
# Directories
#-----------------------------------------#

DIR_COVER ?= ./.cover
GO_COVERPKG ?= ./...
COVER_IGNORED_FILES = "tests\|examples/"

#-----------------------------------------#
# Files
#-----------------------------------------#

COVER_UNIT_FILE ?= unit.out

#-----------------------------------------#
# Tests
#-----------------------------------------#

.PHONY: _test/general
_test/general: _cover/setup

.PHONY: _test/internal/unit
_test/internal/unit: GO_PKGS_UNIT_TEST ?= `go list $(GO_COVERPKG) | grep -v ${COVER_IGNORED_FILES}`
_test/internal/unit:
	@echo "Running go unit tests..."
	@$(GOTEST) -race -failfast -v \
		-coverpkg=${GO_COVERPKG} -covermode=atomic -coverprofile=${DIR_COVER}/${COVER_UNIT_FILE}.tmp \
		${GO_PKGS_UNIT_TEST}

.PHONY: test/local/unit
test/local/unit: _test/general _test/internal/unit _cover/patch 

.PHONY: test/action/unit
test/action/unit: _test/general _test/internal/unit _cover/patch 

#-----------------------------------------#
# Cover
#-----------------------------------------#

.PHONY: cover
cover: _cover/unit

.PHONY: _cover/unit
_cover/unit:
	@echo "Coverage for: $(COVER_UNIT_FILE)"
	@$(GOTOOL) cover -func=$(DIR_COVER)/$(COVER_UNIT_FILE)

.PHONY: _cover/setup
_cover/setup:
	@rm -rf $(DIR_COVER)
	@mkdir -p $(DIR_COVER)

.PHONY: _cover/patch
_cover/patch:
	@touch ${DIR_COVER}/${COVER_UNIT_FILE}.tmp
	@cat ${DIR_COVER}/${COVER_UNIT_FILE}.tmp | { grep -v ${COVER_IGNORED_FILES} > ${DIR_COVER}/${COVER_UNIT_FILE} || true; }

#-----------------------------------------#
# Report
#-----------------------------------------#

.PHONY: report/unit
report/unit:
	@$(GO) tool cover -o ${DIR_COVER}/unit.html -html=${DIR_COVER}/${COVER_UNIT_FILE}
