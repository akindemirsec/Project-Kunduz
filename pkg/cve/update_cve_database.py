import requests
import psycopg2
from psycopg2 import sql
import xml.etree.ElementTree as ET

# PostgreSQL veritabanı bağlantı ayarları
DB_HOST = "db"
DB_NAME = "kunduz"
DB_USER = "user"
DB_PASS = "password"

# MITRE CVE XML dosyası URL'si
MITRE_XML_URL = "https://cve.mitre.org/data/downloads/allitems.xml"

def download_cve_data():
    try:
        response = requests.get(MITRE_XML_URL)
        response.raise_for_status()
        with open('/app/allitems.xml', 'wb') as file:
            file.write(response.content)
        print("CVE data downloaded successfully.")
    except requests.exceptions.RequestException as e:
        print(f"Failed to download CVE data: {e}")

def parse_cve_data(filename):
    tree = ET.parse(filename)
    root = tree.getroot()
    cve_items = []

    for entry in root.findall('.//entry'):
        cve_id = entry.get('id')
        description = entry.findtext('desc', default='No description available')
        cve_items.append((cve_id, description))

    return cve_items

def insert_cve_data(cve_data):
    conn = psycopg2.connect(
        dbname=DB_NAME, user=DB_USER, password=DB_PASS, host=DB_HOST)
    cursor = conn.cursor()

    for cve_id, description in cve_data:
        cursor.execute(
            sql.SQL("INSERT INTO cves (cve_id, description) VALUES (%s, %s) ON CONFLICT (cve_id) DO NOTHING"),
            [cve_id, description]
        )

    conn.commit()
    cursor.close()
    conn.close()

if __name__ == "__main__":
    download_cve_data()
    cve_data = parse_cve_data('/app/allitems.xml')
    insert_cve_data(cve_data)
    print("CVE veritabanı güncellendi.")
