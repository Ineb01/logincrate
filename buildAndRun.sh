export $(xargs < public.env)
docker build -t $LOGINCRATE_REGISTRY/$LOGINCRATE_IMAGENAME:$LOGINCRATE_VERSION .