import requests

# note: will fail, use real username/password
BASE_URL = "http://localhost:8080/api/v1"
r = requests.post(f"{BASE_URL}/auth/login", json = {"username":"demo.user", "password": "demo_password"})

print(r.json())