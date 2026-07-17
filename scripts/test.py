import requests
import time

BASE_URL = "http://localhost:8080/api/v1"

HEADERS = {
    "Authorization" : "Bearer demo-key"
}

def test_api(url, params = {}, headers=HEADERS):
    if params != {}:
        r = requests.get(url, headers = headers, params = params)
    else:
        r = requests.get(url, headers= headers)
    if r.status_code == 200:
        data = r.json()
        print(f"SUCCESS on GET {url}")
        # print(data)
    else:
        print(f"FAILURE on GET {url}")

# first test
test_api(f"{BASE_URL}/students", headers=HEADERS, params= {"page":1, "limit":3177})
time.sleep(1)

# second test
test_sid = "ec317086-b0a9-432e-adc9-40935eb453a3"
test_api(f"{BASE_URL}/students/{test_sid}")
time.sleep(1)

# third test
test_api(f"{BASE_URL}/students/{test_sid}/schedules")
time.sleep(1)

# fourth test
test_api(f"{BASE_URL}/classes")
time.sleep(1)

# fifth test
test_api(f"{BASE_URL}/students", headers={})
time.sleep(1)