#!/usr/bin/env bash
# Create Yandex Cloud VM + security group for Тихвинский путь.
# Requires: yc CLI authenticated (`yc init`), SSH key at ~/.ssh/id_ed25519.pub
set -euo pipefail

export PATH="${HOME}/yandex-cloud/bin:${PATH}"

ZONE="${YC_ZONE:-ru-central1-a}"
NETWORK_NAME="${YC_NETWORK_NAME:-palomnik-net}"
SUBNET_NAME="${YC_SUBNET_NAME:-palomnik-subnet-a}"
SG_NAME="${YC_SG_NAME:-palomnik-sg}"
VM_NAME="${YC_VM_NAME:-palomnik-vm}"
SSH_KEY="${SSH_KEY:-${HOME}/.ssh/id_ed25519.pub}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

if ! yc config list >/dev/null 2>&1; then
  echo "yc is not authenticated. Run:  export PATH=\"\$HOME/yandex-cloud/bin:\$PATH\" && yc init"
  exit 1
fi

if [[ ! -f "$SSH_KEY" ]]; then
  echo "SSH public key not found: $SSH_KEY"
  exit 1
fi

if ! yc vpc network get --name "$NETWORK_NAME" >/dev/null 2>&1; then
  echo "Creating network $NETWORK_NAME"
  yc vpc network create --name "$NETWORK_NAME" --description "tikhvin-palomnik"
fi

if ! yc vpc subnet get --name "$SUBNET_NAME" >/dev/null 2>&1; then
  echo "Creating subnet $SUBNET_NAME"
  yc vpc subnet create \
    --name "$SUBNET_NAME" \
    --zone "$ZONE" \
    --network-name "$NETWORK_NAME" \
    --range 10.10.0.0/24
fi

if ! yc vpc security-group get --name "$SG_NAME" >/dev/null 2>&1; then
  echo "Creating security group $SG_NAME"
  yc vpc security-group create \
    --name "$SG_NAME" \
    --network-name "$NETWORK_NAME" \
    --rule "direction=ingress,port=22,protocol=tcp,v4-cidrs=[0.0.0.0/0]" \
    --rule "direction=ingress,port=80,protocol=tcp,v4-cidrs=[0.0.0.0/0]" \
    --rule "direction=ingress,port=443,protocol=tcp,v4-cidrs=[0.0.0.0/0]" \
    --rule "direction=egress,protocol=any,v4-cidrs=[0.0.0.0/0]"
fi

if ! yc compute instance get --name "$VM_NAME" >/dev/null 2>&1; then
  echo "Creating VM $VM_NAME"
  yc compute instance create \
    --name "$VM_NAME" \
    --hostname "$VM_NAME" \
    --zone "$ZONE" \
    --platform standard-v3 \
    --cores 2 \
    --memory 4GB \
    --core-fraction 100 \
    --create-boot-disk size=40,type=network-ssd,image-folder-id=standard-images,image-family=ubuntu-2404-lts \
    --network-interface "subnet-name=${SUBNET_NAME},nat-ip-version=ipv4,security-group-name=${SG_NAME}" \
    --ssh-key "$SSH_KEY" \
    --metadata-from-file "user-data=${ROOT_DIR}/deploy/yandex/user-data.yaml"
else
  echo "VM $VM_NAME already exists"
fi

IP="$(yc compute instance get --name "$VM_NAME" --format json | python3 -c "import json,sys; n=json.load(sys.stdin)['network_interfaces'][0]; print(n['primary_v4_address']['one_to_one_nat']['address'])")"
echo "PUBLIC_IP=$IP"
echo
echo "REG.RU DNS A-records (tikhvin-palomnik.ru):"
echo "  @                      A    $IP"
echo "  www                    A    $IP"
echo "  api                    A    $IP"
echo
echo "SSH: ssh yc-user@$IP"
