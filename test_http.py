import requests
import uuid
from requests.auth import HTTPBasicAuth

url = "http://localhost:8080/payment-request"
headers = {
    "Content-Type": "application/json",
}

auth = HTTPBasicAuth("BOLWERK", "xxxx")  # Replace 'xxxx' with the actual password

for _ in range(10):
    data = {
        "debtor_iban": "FR1112739000504482744411A64",
        "debtor_name": "company1",
        "creditor_iban": "DE65500105179799248552",
        "creditor_name": "beneficiary",
        "ammount": 42.99,
        "idempotency_unique_key": str(uuid.uuid4())[:10]  # Generate a unique UUID
    }

    response = requests.post(url, json=data, headers=headers, auth=auth)

print("Status Code:", response.status_code)
print("Response Body:", response.text)