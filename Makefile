CURRENT_TIME=$(shell date +%s)
VERSION_MAJOR=$(shell cat production-version-major.txt)
VERSION_MINOR=$(shell cat production-version-minor.txt)
VERSION_PATCH=$(shell cat production-version-patch.txt)
BUILD_NUMBER=$(shell cat production-build-number.txt)
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + 1 )))
PRODUCTION_BINARY_NAME=./production/c3rl-cli_${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}

.PHONY: production
production:

	GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w -X 'main.BuildType=production' -X 'main.command_version_time_unix=${CURRENT_TIME}' -X 'main.command_version_major=${VERSION_MAJOR}' -X 'main.command_version_minor=${VERSION_MINOR}' -X 'main.command_version_patch=${VERSION_PATCH}' -X 'main.command_version_build_number=${BUILD_NUMBER}'" -o ${PRODUCTION_BINARY_NAME}_darwin_amd64

	GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w -X 'main.BuildType=production' -X 'main.command_version_time_unix=${CURRENT_TIME}' -X 'main.command_version_major=${VERSION_MAJOR}' -X 'main.command_version_minor=${VERSION_MINOR}' -X 'main.command_version_patch=${VERSION_PATCH}' -X 'main.command_version_build_number=${BUILD_NUMBER}'" -o ${PRODUCTION_BINARY_NAME}_darwin_arm64

	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w -X 'main.BuildType=production' -X 'main.command_version_time_unix=${CURRENT_TIME}' -X 'main.command_version_major=${VERSION_MAJOR}' -X 'main.command_version_minor=${VERSION_MINOR}' -X 'main.command_version_patch=${VERSION_PATCH}' -X 'main.command_version_build_number=${BUILD_NUMBER}'" -o ${PRODUCTION_BINARY_NAME}_linux_amd64

	GOARCH=arm64 GOOS=linux go build -ldflags="-s -w -X 'main.BuildType=production' -X 'main.command_version_time_unix=${CURRENT_TIME}' -X 'main.command_version_major=${VERSION_MAJOR}' -X 'main.command_version_minor=${VERSION_MINOR}' -X 'main.command_version_patch=${VERSION_PATCH}' -X 'main.command_version_build_number=${BUILD_NUMBER}'" -o ${PRODUCTION_BINARY_NAME}_linux_arm64 

	GOARCH=arm GOOS=linux go build -ldflags="-s -w -X 'main.BuildType=production' -X 'main.command_version_time_unix=${CURRENT_TIME}' -X 'main.command_version_major=${VERSION_MAJOR}' -X 'main.command_version_minor=${VERSION_MINOR}' -X 'main.command_version_patch=${VERSION_PATCH}' -X 'main.command_version_build_number=${BUILD_NUMBER}'" -o ${PRODUCTION_BINARY_NAME}_linux_arm 

	# update build number
	# sed -i 	"s/${BUILD_NUMBER}\$$/${NEW_BUILD_NUMBER}/g" production-build-number.txt 
	echo ${NEW_BUILD_NUMBER} > production-build-number.txt 
