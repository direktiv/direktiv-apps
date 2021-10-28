REPOSITORY=$(printenv REPOSITORY)
CONTAINER=$(printenv CONTAINER)
echo $CONTAINER/$REPOSITORY
VERSION=$(curl -L -s "https://registry.hub.docker.com/v2/repositories/$REPOSITORY/$CONTAINER/tags?page_size=1024"|jq ".\"results\"[1].name")
VERSION_WITHOUT_QUOTES=`echo $VERSION | sed 's/.\(.*\)/\1/' | sed 's/\(.*\)./\1/'`
VERSION_NUMBER=`echo $VERSION_WITHOUT_QUOTES | cut -c 2-`
VERSION_NUMBER=$(($VERSION_NUMBER+1))

echo v$VERSION_NUMBER > $CONTAINER/VERSION
docker build $CONTAINER -t $REPOSITORY/$CONTAINER:latest -t $REPOSITORY/$CONTAINER:v$VERSION_NUMBER
docker push $REPOSITORY/$CONTAINER:v$VERSION_NUMBER
docker push $REPOSITORY/$CONTAINER:latest

