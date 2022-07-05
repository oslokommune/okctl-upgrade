# Functions
function run() {
  CMD=$@
  echo
  echo "~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~"
  echo -e "Running command: \e[96m${CMD}\e[0m"
  echo "~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~"
  $CMD
  ERROR_CODE=$?

  echo "~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~"

  if [[ ! $ERROR_CODE == 0 ]]; then
    echo Command failed with error code $ERROR_CODE: "$CMD"
    echo Aborting.

    exit $ERROR_CODE
  fi
}
# Args check
if [[ $* == "-h" || -z "$1" ]]
then
    ME=$(basename $0)
    echo "USAGE:"
    echo "$ME cluster-manifest.yaml"
    exit 0
fi

CLUSTER_MANIFEST="$1"

if [[ ! -f "$CLUSTER_MANIFEST" ]]; then
  echo "File does not exist: $CLUSTER_MANIFEST"
fi

# Test requirements
if ! command -v yq &> /dev/null
then
    echo 'yq' could not be found. Install before retrying.
    exit
fi

CLUSTER_NAME=$(yq e '.metadata.name' $CLUSTER_MANIFEST)

echo "Upgrading EKS to 1.22"
echo "Cluster manifest: $CLUSTER_CLUSTER_MANIFEST"
echo "Cluster name: $CLUSTER_NAME"

echo "------------------------------------------------------------------------------------------------------------------------"
echo "Verify that cluster exists"
echo "------------------------------------------------------------------------------------------------------------------------"
run "eksctl get cluster $CLUSTER_NAME"

