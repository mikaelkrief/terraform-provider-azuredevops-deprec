#!/usr/bin/env bash
VERSION=$1
OUTPUTPATH=$2
package_name="terraform-provider-azuredevops_v"
platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/arm" "linux/386")



# clean up
echo "-> running clean up...."
rm -rf ${OUTPUTPATH}/*

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="${OUTPUTPATH}/${GOOS}_${GOARCH}/${package_name}${VERSION}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    echo "--> building $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done


for PLATFORM in $(find ${OUTPUTPATH} -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../terraform-provider-azuredevops_${OSARCH}_v${VERSION}.zip ./*
    popd >/dev/null 2>&1
done