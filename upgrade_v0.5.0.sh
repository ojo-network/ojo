#!/bin/bash

# Set variables
CHAIN_ID="ojo"
ALICE_KEY="alice"
UPGRADE_NAME="v0.5.0"
UPGRADE_HEIGHT=30  # Adjust this to an appropriate block height
VOTING_PERIOD="20s"
BINARY="ojod"

# Ensure the binary is in the PATH
if ! command -v $BINARY &> /dev/null; then
    echo "$BINARY could not be found. Please ensure it's installed and in your PATH."
    exit 1
fi

# Create a temporary JSON file for the proposal
PROPOSAL_FILE=$(mktemp)

# Get the gov module's address
GOV_ADDRESS=$($BINARY query auth module-account gov --output json | jq -r .account.value.address)

echo "GOV_ADDRESS: $GOV_ADDRESS"
# Write the proposal JSON
cat << EOF > "$PROPOSAL_FILE"
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "$GOV_ADDRESS",
      "plan": {
        "name": "$UPGRADE_NAME",
        "height": "$UPGRADE_HEIGHT",
        "info": "Performing upgrade to $UPGRADE_NAME"
      }
    }
  ],
  "title": "Upgrade to $UPGRADE_NAME",
  "summary": "Performing upgrade to $UPGRADE_NAME",
  "metadata": "{\"title\":\"Upgrade to $UPGRADE_NAME\",\"summary\":\"Performing upgrade to $UPGRADE_NAME\"}",
  "deposit": "10000000uojo"
}
EOF

# Submit upgrade proposal
echo "Submitting upgrade proposal..."
$BINARY tx gov submit-proposal "$PROPOSAL_FILE" \
    --from $ALICE_KEY \
    --chain-id $CHAIN_ID \
    --keyring-backend test \
    --gas auto \
    -b sync \
    -y

# Remove the temporary proposal file
rm "$PROPOSAL_FILE"

# Wait for a moment to ensure the proposal is processed
sleep 5

# Get the proposal ID
PROPOSAL_ID=$($BINARY query gov proposals --output json | jq '.proposals[-1].id' -r)

# Vote on the proposal
echo "Voting on proposal $PROPOSAL_ID..."
$BINARY tx gov vote $PROPOSAL_ID yes \
    --from $ALICE_KEY \
    --chain-id $CHAIN_ID \
    --keyring-backend test \
    --gas auto \
    -y

echo "Upgrade proposal submitted and voted. Proposal ID: $PROPOSAL_ID"
echo "Upgrade scheduled for block height: $UPGRADE_HEIGHT"
echo "Please ensure all validators upgrade their software before this block height."

# Wait for the voting period
echo "Waiting for voting period to end..."
sleep $VOTING_PERIOD

# Check proposal status
PROPOSAL_STATUS=$($BINARY query gov proposal $PROPOSAL_ID --output json | jq '.status' -r)
echo "Proposal status: $PROPOSAL_STATUS"

if [ "$PROPOSAL_STATUS" == "PROPOSAL_STATUS_PASSED" ]; then
    echo "Proposal passed. Prepare for upgrade at block height $UPGRADE_HEIGHT."
else
    echo "Proposal did not pass. Please check the voting results."
fi
