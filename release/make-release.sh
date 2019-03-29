#!/bin/bash

VER=$1
PROJECT=$2

DOTCNT=$(echo $VER | grep -o "\." | wc -l)
if [ "$DOTCNT" -gt 1 ]; then
  PRE="true"
else
  PRE="false"
  VER="${VER}.0"
fi

if [ -z "$PROJECT" ]; then
  echo "Project not specified. Exiting..."
  exit 0
fi

echo Creating a new release: $VER

IFS="." read -a PARTS <<< "$VER"
MAJOR=${PARTS[0]}
MINOR=${PARTS[1]}
MINOR=$((MINOR+1))
VERN=${MAJOR}.${MINOR}

pushd $GOPATH/src/v2ray.com/core
echo "Adding a new tag: " "v$VER"
git tag -s -a "v$VER" -m "Version ${VER}"
sed -i '' "s/\(version *= *\"\).*\(\"\)/\1$VERN\2/g" core.go
echo "Committing core.go (may not necessary)"
git commit core.go -S -m "Update version"
echo "Pushing changes"
git push --follow-tags
popd

echo "Launching build machine."
DIR="$(dirname "$0")"
RAND="$(openssl rand -hex 5)"
gcloud compute instances create "v2raycore-${RAND}" \
    --machine-type=n1-standard-2 \
    --metadata=release_tag=v${VER},prerelease=${PRE} \
    --metadata-from-file=startup-script=${DIR}/release-ci.sh \
    --zone=europe-west4-c \
    --project ${PROJECT} \
    --scopes "https://www.googleapis.com/auth/compute,https://www.googleapis.com/auth/devstorage.read_write" \

