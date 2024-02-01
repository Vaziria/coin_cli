import os
import re


PRIV_NAME = "private_key.pem"
PUB_NAME = "public_key.hex"

def get_public_key(pub_location):
    with open(pub_location, "r") as out:
        data = out.read().replace("\n", "")
        data = re.findall('pub:(.*?)ASN', data)
        data = data[0]
        data = data.strip()
        data = data.replace(":", "").replace(" ", "").replace("\t", "")

        return data


priv_location = f"/root/secretkey/{PRIV_NAME}"
if not os.path.exists(priv_location):
    os.system(f"openssl ecparam -genkey -name secp256k1 -out {priv_location}")
    
pub_location = f"/root/secretkey/{PUB_NAME}"
if not os.path.exists(pub_location):
    os.system(f"openssl ec -in {priv_location} -text > {pub_location}")



print(get_public_key(pub_location))


