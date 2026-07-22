#!/bin/sh

if [ "$(which dpkg-deb)" = "" ]; then
	echo "ERROR: 'dpkg-deb' is required for building .deb packages!"
fi

cd $(dirname $0)/..  # ensure correct dir

mkdir -p .generated/deb/DEBIAN
mkdir -p .generated/deb/usr/bin

install -m 0755 $OUTDIR/$NAME -t .generated/deb/usr/bin

cat > .generated/deb/DEBIAN/control <<EOF
Package: $NAME
Version: $VERSION
Architecture: $ARCH
Maintainer: Tomte Bender <bender@solidapo.de>
Priority: optional
Section: database
Description: A database that stores relations using tags.
 LisDB is a database that stores relations between nodes using tags.
EOF
# TODO: write longer description

dpkg-deb --build --root-owner-group .generated/deb $OUTDIR/${NAME}_${VERSION}_${ARCH}.deb

exit $?
