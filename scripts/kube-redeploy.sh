at=$(date +%s)
echo "redeployed at $at"
kubectl config use-context docker-desktop
kubectl patch deployment user-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${at}'" }}}}}'
kubectl patch deployment auth-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${at}'" }}}}}'
kubectl patch deployment post-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${at}'" }}}}}'
