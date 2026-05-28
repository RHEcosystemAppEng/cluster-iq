# Vault Integration Guide for Cluster IQ

Cluster IQ supports HashiCorp Vault for secure storage and retrieval of cloud provider credentials. This guide explains how to configure your existing Vault instance and deploy Cluster IQ with Vault integration enabled.

> **Note:** If you do not have a Vault instance and want to install one for testing or development, please refer to the [Reference Installation Guide](#appendix-reference-installation-guide) in the Appendix.

## Prerequisites

1. **Secrets Store CSI Driver Operator** installed on the cluster.
    * If not installed, follow the [CSI Driver Installation](#1-install-secrets-store-csi-driver-operator) steps in the Appendix.
2. **HashiCorp Vault** installed and accessible.

## Part I: Vault Configuration Requirements

Ensure your Vault instance is configured with the following resources.

### 1. Kubernetes Authentication

Vault must have the Kubernetes authentication method enabled and configured to communicate with your OpenShift cluster.

### 2. Policy

Create a policy named `cluster-iq-credentials` that allows reading the secret.

* **Example Policy File:** [`doc/vault/manifests/20-vault-policy.hcl`](manifests/20-vault-policy.hcl)

### 3. Kubernetes Roles

Create roles to bind the application's ServiceAccounts to the policy.

| Role Name | ServiceAccount Name | Namespace | Policy |
|-----------|---------------------|-----------|--------|
| `scanner` | `scanner` *         | `cluster-iq` | `cluster-iq-credentials` |
| `agent`   | `agent` *           | `cluster-iq` | `cluster-iq-credentials` |

*\* Note: ServiceAccount names must match those defined in your Helm values.*

### 4. Credentials Secret

Store your cloud provider credentials in the KV v2 secrets engine.

* **Secret Path:** `secret/cluster-iq/credentials`
* **Key:** `credentials`
* **Value:** Content of your credentials file.

---

## Part II: Application Deployment

Once Vault is configured, you can deploy (or upgrade) Cluster IQ with Vault integration enabled.

**1. Create a Helm values file (`vault-config.yaml`)**

Enable Vault and provide the connection details.

```yaml
vault:
  enabled: true
  # Address of your Vault instance
  address: "http://vault.hashicorp-vault.svc.cluster.local:8200"
  # Path to the secret
  secretPath: "secret/data/cluster-iq/credentials"
  # Key within the secret
  secretKey: "credentials"
```

**2. Deploy/Upgrade Cluster IQ**

```bash
export APP_NS="cluster-iq"
helm upgrade --install cluster-iq deployments/helm/cluster-iq -n $APP_NS -f vault-config.yaml
```

**3. Verification**

```bash
# Check scanner logs
oc logs deployment/scanner -n $APP_NS

# Verify credentials mount in scanner pod
oc exec -n $APP_NS deployment/scanner -- ls -la /credentials

# Verify credentials mount in agent pod
oc exec -n $APP_NS deployment/agent -- ls -la /credentials
```

---

# Appendix: Reference Installation Guide

This section provides a complete, step-by-step guide to installing and configuring a reference Vault instance on OpenShift. Use this for development, testing, or if you are setting up Vault from scratch.

## 1. Install Secrets Store CSI Driver Operator

### 1.1 Create Subscription

```bash
export CSI_NS="openshift-cluster-csi-drivers"
oc create -f doc/vault/manifests/00-csi-operator-subscription.yaml
```

### 1.2 Approve InstallPlan (if Manual approval mode)

```bash
oc patch installplan $(oc get subscription secrets-store-csi-driver-operator -n $CSI_NS -o jsonpath='{.status.installPlanRef.name}') -n $CSI_NS --type merge -p '{"spec":{"approved":true}}'
```

### 1.3 Verify Installation

```bash
oc get csv -n $CSI_NS | grep secrets-store
```

### 1.4 Create ClusterCSIDriver

```bash
oc create -f doc/vault/manifests/01-csi-driver.yaml
```

### 1.5 Verify CSI Driver pods

```bash
oc get pods -n $CSI_NS | grep secrets-store
```

## 2. Install Reference Vault Instance

Use this to set up a production-ready Vault instance on OpenShift for testing/demo purposes.

**Choose one of the following installation options:**

* **[Option A: Standard Installation (Manual Unseal)](#option-a-standard-installation-manual-unseal)**
  * Simpler setup, no external dependencies.
  * Requires manual unseal with keys after every restart.
* **[Option B: Production Installation (AWS KMS Auto-Unseal)](#option-b-production-installation-aws-kms-auto-unseal)**
  * Production-grade setup.
  * Requires AWS access (KMS, IAM).
  * Vault unseals automatically.

### Common Preparation (Required for BOTH options)

```bash
# Add Helm repo
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Create Namespace
export VAULT_NS="hashicorp-vault"
oc new-project $VAULT_NS

# Configure SCC
oc adm policy add-scc-to-user privileged -z vault-csi-provider -n $VAULT_NS
```

### Option A: Standard Installation (Manual Unseal)

**1. Install Vault**

```bash
helm upgrade --install vault hashicorp/vault -n $VAULT_NS \
    -f doc/vault/manifests/10-vault-values-standard.yaml
```

**2. Initialize & Unseal**

```bash
# Initialize vault-0
oc exec -n $VAULT_NS vault-0 -- vault operator init -key-shares=5 -key-threshold=3 -format=json > /tmp/vault-init.json

# Extract keys
UNSEAL_KEY_1=$(jq -r '.unseal_keys_b64[0]' /tmp/vault-init.json)
UNSEAL_KEY_2=$(jq -r '.unseal_keys_b64[1]' /tmp/vault-init.json)
UNSEAL_KEY_3=$(jq -r '.unseal_keys_b64[2]' /tmp/vault-init.json)

# Unseal vault-0
oc exec -n $VAULT_NS vault-0 -- vault operator unseal $UNSEAL_KEY_1
oc exec -n $VAULT_NS vault-0 -- vault operator unseal $UNSEAL_KEY_2
oc exec -n $VAULT_NS vault-0 -- vault operator unseal $UNSEAL_KEY_3

# Join and unseal followers
for pod in vault-1 vault-2; do
  echo "Joining and unsealing $pod..."
  oc exec -n $VAULT_NS $pod -- vault operator raft join http://vault-0.vault-internal:8200
  oc exec -n $VAULT_NS $pod -- vault operator unseal $UNSEAL_KEY_1
  oc exec -n $VAULT_NS $pod -- vault operator unseal $UNSEAL_KEY_2
  oc exec -n $VAULT_NS $pod -- vault operator unseal $UNSEAL_KEY_3
done
```

### Option B: Production Installation (AWS KMS Auto-Unseal)

**1. Set AWS Variables**

```bash
export AWS_REGION="us-east-1"
export IAM_USER_NAME="vault-unseal-user"
export KMS_KEY_TAG_NAME="vault-autounseal-key"
```

**2. Create AWS KMS key**

```bash
aws kms create-key \
    --description "Vault auto-unseal key" \
    --tags TagKey=Name,TagValue=$KMS_KEY_TAG_NAME \
    --region $AWS_REGION \
    --output json | tee /tmp/kms-key.json

# Extract Key ARN
export KMS_KEY_ARN=$(jq -r '.KeyMetadata.Arn' /tmp/kms-key.json)
echo "KMS Key ARN: $KMS_KEY_ARN"
```

**3. Configure IAM Access**

Choose one of the following sub-options to configure IAM access.

* **Sub-Option B1 (Automatic):** Uses OpenShift `CredentialsRequest` (Mint Mode). **Requirement:** You must have permissions to create `CredentialsRequest` objects in the cluster.
* **Sub-Option B2 (Manual):** Manually creates IAM users/policies and Kubernetes Secrets.

**Sub-Option B1: Automatic Setup (OpenShift Mint Mode)**

```bash
# Prepare Manifest
cp doc/vault/manifests/20-vault-credentials-request.yaml /tmp/vault-credentials-request.yaml
KMS_KEY_ARN_ESCAPED=$(echo $KMS_KEY_ARN | sed 's/\//\\\//g')
sed -i "s/REPLACE_WITH_YOUR_KMS_KEY_ARN/$KMS_KEY_ARN_ESCAPED/g" /tmp/vault-credentials-request.yaml
sed -i "s/namespace: hashicorp-vault/namespace: $VAULT_NS/g" /tmp/vault-credentials-request.yaml

# Apply
oc apply -f /tmp/vault-credentials-request.yaml

# Wait for Secret
while ! oc get secret vault-kms-creds -n $VAULT_NS > /dev/null 2>&1; do sleep 5; echo "Waiting..."; done
```

**Sub-Option B2: Manual Setup (Universal)**

```bash
# Create IAM policy
cat > /tmp/vault-kms-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["kms:Encrypt", "kms:Decrypt", "kms:DescribeKey"],
      "Resource": "$KMS_KEY_ARN"
    }
  ]
}
EOF

aws iam create-policy --policy-name VaultKMSUnseal --policy-document file:///tmp/vault-kms-policy.json --output json | tee /tmp/iam-policy.json
export POLICY_ARN=$(jq -r '.Policy.Arn' /tmp/iam-policy.json)

# Create IAM user
aws iam create-user --user-name $IAM_USER_NAME --tags Key=Purpose,Value=VaultAutoUnseal
aws iam attach-user-policy --user-name $IAM_USER_NAME --policy-arn $POLICY_ARN
aws iam create-access-key --user-name $IAM_USER_NAME --output json | tee /tmp/vault-aws-credentials.json

export VAULT_KMS_ACCESS_KEY_ID=$(jq -r '.AccessKey.AccessKeyId' /tmp/vault-aws-credentials.json)
export VAULT_KMS_SECRET_ACCESS_KEY=$(jq -r '.AccessKey.SecretAccessKey' /tmp/vault-aws-credentials.json)

# Create Kubernetes Secret
oc create secret generic vault-kms-creds -n $VAULT_NS \
    --from-literal=AWS_ACCESS_KEY_ID=$VAULT_KMS_ACCESS_KEY_ID \
    --from-literal=AWS_SECRET_ACCESS_KEY=$VAULT_KMS_SECRET_ACCESS_KEY \
    --from-literal=AWS_REGION=$AWS_REGION
```

**Optional:** Add label to exclude from Velero backups

```bash
oc label secret vault-kms-creds -n $VAULT_NS velero.io/exclude-from-backup=true
```

**4. Install Vault with Auto-Unseal**

```bash
# Prepare values file
cp doc/vault/manifests/10-vault-values-autounseal.yaml /tmp/vault-values-autounseal.yaml
KMS_KEY_ARN_ESCAPED=$(echo $KMS_KEY_ARN | sed 's/\//\\\//g')
sed -i "s/REPLACE_WITH_YOUR_REGION/$AWS_REGION/g" /tmp/vault-values-autounseal.yaml
sed -i "s/REPLACE_WITH_YOUR_KMS_KEY_ARN/$KMS_KEY_ARN_ESCAPED/g" /tmp/vault-values-autounseal.yaml

# Install Vault
helm upgrade --install vault hashicorp/vault -n $VAULT_NS \
    -f /tmp/vault-values-autounseal.yaml
```

**5. Initialize Vault**

```bash
# Initialize
oc exec -n $VAULT_NS vault-0 -- vault operator init -format=json > /tmp/vault-init.json

# Join followers
oc exec -n $VAULT_NS vault-1 -- vault operator raft join http://vault-0.vault-internal:8200
oc exec -n $VAULT_NS vault-2 -- vault operator raft join http://vault-0.vault-internal:8200

# Verify status
oc exec -n $VAULT_NS vault-0 -- vault status
```

## 3. Configure Vault (Reference)

This section details how to configure the reference Vault instance for Cluster IQ.

**3.1 Enable Kubernetes Authentication**

```bash
# Get Root Token
ROOT_TOKEN=$(jq -r '.root_token' /tmp/vault-init.json)
oc exec -n $VAULT_NS vault-0 -- vault login $ROOT_TOKEN

# Enable auth method
oc exec -n $VAULT_NS vault-0 -- vault auth enable kubernetes

# Configure auth
oc exec -n $VAULT_NS vault-0 -- sh -c 'vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"'
```

**3.2 Create Policy and Roles**

```bash
export APP_NS="cluster-iq"

# Create Policy
cat doc/vault/manifests/20-vault-policy.hcl | \
  oc exec -i -n $VAULT_NS vault-0 -- vault policy write cluster-iq-credentials -

# Create Role for Scanner
oc exec -n $VAULT_NS vault-0 -- vault write auth/kubernetes/role/scanner \
    bound_service_account_names=scanner \
    bound_service_account_namespaces=$APP_NS \
    policies=cluster-iq-credentials \
    ttl=24h \
    audience=vault

# Create Role for Agent
oc exec -n $VAULT_NS vault-0 -- vault write auth/kubernetes/role/agent \
    bound_service_account_names=agent \
    bound_service_account_namespaces=$APP_NS \
    policies=cluster-iq-credentials \
    ttl=24h \
    audience=vault
```

**3.3 Load Credentials**

```bash
# Enable KV v2 engine
oc exec -n $VAULT_NS vault-0 -- vault secrets enable -path=secret kv-v2

# Replace with the real path to credentials file
cat /path/to/credentials/file | oc exec -i -n $VAULT_NS vault-0 -- \
  sh -c "VAULT_TOKEN=$ROOT_TOKEN vault kv put secret/cluster-iq/credentials credentials=-"
```

## 4. Troubleshooting & Verification

**Check CSI Driver Logs**
If pods are failing to start, check the CSI driver logs on the node where the pod is scheduled.

```bash
oc logs -n openshift-cluster-csi-drivers -l app=secrets-store-csi-driver-node
```

**Verify Secret Content in Vault**

```bash
oc exec -n $VAULT_NS vault-0 -- vault kv get secret/cluster-iq/credentials
```

**Check Vault Cluster Status**

```bash
# Check overall status (Sealed/Unsealed, HA Mode)
oc exec -n $VAULT_NS vault-0 -- vault status

# List Raft peers (Leader/Follower status)
oc exec -n $VAULT_NS vault-0 -- vault operator raft list-peers

# Check all pods status
oc get pods -n $VAULT_NS
```

**Login to Vault CLI**

```bash
# Login with Root Token (from initialization step)
oc exec -ti -n $VAULT_NS vault-0 -- vault login <ROOT_TOKEN>

# Or extract from init file
ROOT_TOKEN=$(jq -r '.root_token' /tmp/vault-init.json)
oc exec -ti -n $VAULT_NS vault-0 -- vault login $ROOT_TOKEN
```

**Common Vault Operations**

```bash
# List all secrets engines
oc exec -n $VAULT_NS vault-0 -- vault secrets list

# List all auth methods
oc exec -n $VAULT_NS vault-0 -- vault auth list

# List all policies
oc exec -n $VAULT_NS vault-0 -- vault policy list

# Read a specific policy
oc exec -n $VAULT_NS vault-0 -- vault policy read cluster-iq-credentials

# List Kubernetes auth roles
oc exec -n $VAULT_NS vault-0 -- vault list auth/kubernetes/role
```

**Restart Application Pods**
If you updated the secret in Vault, restart the application pods to sync changes immediately (otherwise, wait for the rotation interval).

```bash
oc delete pod -l app.kubernetes.io/component=agent -n cluster-iq
```
