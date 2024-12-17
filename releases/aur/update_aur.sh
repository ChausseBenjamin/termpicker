#!/bin/sh

aur_arch="Linux_x86_64"

# Initial setup/Sanity check {{{

# Confirms necessary env variables are present before running the rest
# of the script.
# $1: Env Variable to check
# $@: Message to send to stderr before quitting
assert_env() {
  [ -n "$1" ] || { echo "ERROR: $@" 1>&2; exit 1; }
}

assert_env "$AUR_PRIVATE_KEY" "Couldn't retrieve a private key to publish to the AUR..."
assert_env "$AUR_PUBLIC_KEY"  "Couldn't retrieve a public key to publish to the AUR..."
assert_env "$PKG_REPO_URI"    "Cound't retrieve a URI to pull the package from"
assert_env "$PKG_NAME"        "Couldn't retrieve the package name"
assert_env "$GIT_USER"        "Couldn't retrieve the git username to pull the release from"

latest_tag="$(git ls-remote --tags "$PKG_REPO_URI" | awk '
  # Process lines without ^{} and matching vX.X.X format
  !/\^\{\}$/ && $2 ~ /refs\/tags\/v[0-9]+\.[0-9]+\.[0-9]+$/ {
      tag = $2    # Store the tag reference
  }

  # Print the latest tag without the prefix
  END {
      gsub("refs/tags/v", "", tag)
      print tag
  }')"

# }}}
# Retrieving the checksums for the latest tag {{{
checksum_url="https://github.com/${GIT_USER}/${PKG_NAME}/releases/download/v${latest_tag}/${PKG_NAME}_${latest_tag}_checksums.txt"

checksums="$( wget -q "$checksum_url" -O - )"

checksum="$(echo "$checksums" | awk -v arch="$aur_arch" -v pkg="$PKG_NAME" '{
    for (i = 1; i <= NF; i++) {
        if ($i == pkg "_" arch ".tar.gz") {
            print $(i-1)
        }
    }
  }')"
# }}}
# Cloning and updating the PKGBUILD {{{

git clone "ssh://aur@aur.archlinux.org/${PKG_NAME}"
cd ${PKG_NAME} || { echo "ERROR: could not clone PKGBUILD repo from the aur" 1>&2; exit 1; }

awk -v new_hash="\'$checksum\'" -v new_version="$latest_tag" '
/sha256sums/ {
    # Surround the checksum with single quotes
    $0 = "sha256sums=("'new_hash'")"
}
/^pkgver/ {
    # Only change the pkgver at the beginning of the line
    $0 = "pkgver=" new_version
}
/pkgrel/ {
    # Increment the value of pkgrel by 1
    sub(/^pkgrel=[0-9]+/, "pkgrel=" int($NF) + 1)
}
{ print }' PKGBUILD > PKGBUILD.new && mv PKGBUILD.new PKGBUILD

# }}}
# Commit and push the changes {{{

echo "$AUR_PRIVATE_KEY" > ~/.ssh/id_ed25519
echo "$AUR_PUBLIC_KEY" > ~/.ssh/id_ed25519.pub

git commit -am "Updated package to v${latest_tag}"

# Uncomment only once script is verified and complete:
git push

# }}}
