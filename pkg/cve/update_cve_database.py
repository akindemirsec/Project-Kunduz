import requests
import json
import psycopg2
from psycopg2 import sql

DB_HOST = "db"
DB_NAME = "kunduz"
DB_USER = "user"
DB_PASS = "password"

NVD_API_URL = "https://services.nvd.nist.gov/rest/json/cves/1.0?resultsPerPage=2000"

def fetch_cve_data():
    response = requests.get(NVD_API_URL)
    if response.status_code == 200:
        return response.json()["result"]["CVE_Items"]
    else:
        print(f"Failed to fetch CVE data: {response.status_code}")
        return []

def insert_cve_data(cve_data):
    conn = psycopg2.connect(
        dbname=DB_NAME, user=DB_USER, password=DB_PASS, host=DB_HOST)
    cursor = conn.cursor()

    for item in cve_data:
        cve_id = item["cve"]["CVE_data_meta"]["ID"]
        description = item["cve"]["description"]["description_data"][0]["value"]

        cursor.execute(
            sql.SQL("INSERT INTO cves (cve_id, description) VALUES (%s, %s) ON CONFLICT (cve_id) DO NOTHING"),
            [cve_id, description]
        )

    conn.commit()
    cursor.close()
    conn.close()

if __name__ == "__main__":
    cve_data = fetch_cve_data()
    insert_cve_data(cve_data)
    print("CVE veritabanı güncellendi.")
