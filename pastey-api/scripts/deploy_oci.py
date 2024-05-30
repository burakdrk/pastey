import oci
import sys
import os
import time
import requests

#############################
# ARGV[1] = compartment_id  #
# ARGV[2] = subnet_id       #
# ARGV[3] = image link      #
# ARGV[4] = a domain        #
# ARGV[5] = github sha      #
#############################

DISPLAY_NAME = "pastey-api"
SHAPE = "CI.Standard.A1.Flex"

# Load the default configuration or read from environment variables
try:
    config = oci.config.from_file()
except:
    print("Error loading OCI config file")
    config = {
        "user": os.getenv("OCI_CLI_USER"),
        "key_content": os.getenv("OCI_CLI_KEY_CONTENT"),
        "fingerprint": os.getenv("OCI_CLI_FINGERPRINT"),
        "tenancy": os.getenv("OCI_CLI_TENANCY"),
        "region": os.getenv("OCI_CLI_REGION")
    }
    oci.config.validate_config(config)

# Create a service client
container_instances_client = oci.container_instances.ContainerInstanceClient(config)

# List all container instances
list_container_instances_response = container_instances_client.list_container_instances(sys.argv[1])

# Delete all container instances
deleted = []
for container_instance in list_container_instances_response.data.items:
    if container_instance.lifecycle_state != "DELETED" and container_instance.display_name == DISPLAY_NAME:
        res = container_instances_client.delete_container_instance(container_instance.id)
        print("Requested deletion of container instance")
        deleted.append(container_instance.id)

time.sleep(5)

# Wait for the container instances to be deleted
while True:
    for container_instance_id in deleted:
        container_instance = container_instances_client.get_container_instance(container_instance_id).data
        if container_instance.lifecycle_state == "DELETED":
            print("Container instance deleted")
            deleted.remove(container_instance_id)

    if len(deleted) == 0:
        break

    time.sleep(10)

print("All container instances deleted")

# Create a new container instance
create_container_instance_details = oci.container_instances.models.CreateContainerInstanceDetails(
    availability_domain=sys.argv[4],
    compartment_id=sys.argv[1],
    display_name=DISPLAY_NAME,
    containers=[
        oci.container_instances.models.CreateContainerDetails(
            image_url=sys.argv[3]+":"+sys.argv[5],
        )
    ],
    shape=SHAPE,
    shape_config=oci.container_instances.models.CreateContainerInstanceShapeConfigDetails(ocpus=1.0, memory_in_gbs=6.0),
    vnics=[
        oci.container_instances.models.CreateContainerVnicDetails(
            subnet_id=sys.argv[2]
        )
    ]
)

create_container_instance_response = container_instances_client.create_container_instance(create_container_instance_details)

vnic_id = ""

while True:
    container_instance = container_instances_client.get_container_instance(create_container_instance_response.data.id).data
    if container_instance.lifecycle_state == "ACTIVE":
        print("Container instance running")
        vnic_id = container_instance.vnics[0].vnic_id
        break

    time.sleep(10)


core_client = oci.core.VirtualNetworkClient(config)
get_vnic_response = core_client.get_vnic(vnic_id=vnic_id)

public_ip = get_vnic_response.data.public_ip

# Modify A record in DNS zone (Cloudflare API)
headers = {
    "Authorization": "Bearer " + os.getenv("CF_API_KEY"),
    "Content-Type": "application/json"
}

data = {
    "type": "A",
    "name": "api",
    "content": public_ip,
    "ttl": 1,
    "proxied": True
}

response = requests.put("https://api.cloudflare.com/client/v4/zones/" + os.getenv("CF_ZONE_ID") + "/dns_records/" + os.getenv("CF_DNS_RECORD_ID"), headers=headers, json=data)

if response.json()["success"] != True:
    print("Error updating DNS record")
else:
    print("DNS record updated")

print("Deployment done")