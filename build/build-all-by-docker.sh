#!/bin/bash

cd $(dirname "$0")

# init variable `builds`
source ./build-all.inc.sh

prefix=$(realpath ../)
rm -rf "$prefix/output/"

buildByDocker() {
  gover="$1"
  shift
  docker pull golang:"$gover"

  docker run \
    --rm \
    -v "$prefix":/mnt \
    -e EX_UID="$(id -u)" \
    -e EX_GID="$(id -g)" \
    golang:"$gover" \
    /bin/bash -c '
      sed -i -e "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apt/sources.list;
      apt-get update;
      apt-get install -yq --force-yes git zip;
      /bin/bash /mnt/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID /mnt/output;
    ' \
    'argv_0_placeholder' \
    "$@"
}

gover=1.2
builds=('windows 386 -2000')
#builds=("${builds[@]}" 'freebsd 386 -8' 'freebsd amd64 -8')
buildByDocker "$gover" "${builds[@]}"
