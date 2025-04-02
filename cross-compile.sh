#!/bin/bash

VERSION="$1"
CROSS_COMPILE_OS_AND_ARCHS=$(cat <<-END
darwin	amd64
darwin	arm64
freebsd	386
freebsd	amd64
freebsd	arm
linux	386
linux	amd64
linux	arm
linux	arm64
linux	ppc64
linux	ppc64le
linux	mips
linux	mipsle
linux	mips64
linux	mips64le
netbsd	386
netbsd	amd64
netbsd	arm
openbsd	386
openbsd	amd64
openbsd	arm
windows	386
windows	amd64
END
)

mkdir -p cross-compilation

while IFS= read -r line
do
    os=$(echo "$line" | awk '{print $1}')
    arch=$(echo "$line" | awk '{print $2}')
    echo "Compile Subify $VERSION for OS $os and architecture $arch..."
    if [[ "$os" == "windows" ]]; then
        $(GOOS=$os GOARCH=$arch go build -o ./cross-compilation/subify_${VERSION}_${os}_${arch}.exe)
    else
        $(GOOS=$os GOARCH=$arch go build -o ./cross-compilation/subify_${VERSION}_${os}_${arch})
    fi
done <<< "$CROSS_COMPILE_OS_AND_ARCHS"



