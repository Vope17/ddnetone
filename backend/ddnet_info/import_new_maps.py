import psycopg2
import os
from pathlib import Path
from dotenv import load_dotenv


# 讀取 .env
env_path = Path(__file__).resolve().parent.parent.parent / '.env'
load_dotenv(dotenv_path=env_path)

try:
    host = "localhost"
    user = os.environ['DB_USER']
    password = os.environ['DB_PASSWORD']
    dbname = os.environ['DB_NAME']
    port = "5433"
    print(f"成功從 {env_path} 讀取設定")
except KeyError as e:
    print(f"找不到 DB 環境變數: {e}")

DB_CONFIG = {
    "host": host,
    "user": user,
    "password": password,
    "dbname": dbname,
    "port": port
}

INPUT_FILE = './insane_maps.txt'
TARGET_DIFFICULTY = 'INSANE'

def calculate_points(star):
    """ 
    保留計算邏輯，僅供新增地圖時使用 
    """
    try:
        s = int(star)
        return 15 + (s * 3)
    except:
        return 0

def main():
    if not os.path.exists(INPUT_FILE):
        print(f"❌ 找不到檔案: {INPUT_FILE}")
        return

    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()
        print("🔌 資料庫連線成功")

        file_map_data = {} 
        with open(INPUT_FILE, 'r', encoding='utf-8') as f:
            lines = f.readlines()


        print("📖 正在解析地圖資料...")
        for line in lines:
            line = line.strip()
            if not line or "|" not in line or "[source" in line:
                continue

            parts = line.split('|')
            if len(parts) >= 2 and parts[0].isdigit():
                star_val = int(parts[0])

                map_name = parts[1].strip()
                points_val = calculate_points(star_val)
                
                file_map_data[map_name] = {
                    'points': points_val,
                    'stars': star_val
                }

        to_insert = []
        to_update = []
        
        cursor.execute("SELECT map_name FROM map_records")
        db_maps = {row[0] for row in cursor.fetchall()}


        for name, data in file_map_data.items():
            if name not in db_maps:
                # 新增地圖時，還是會寫入初始星級與分數
                to_insert.append((TARGET_DIFFICULTY, name, data['stars'], data['points']))
            else:
                # ★ 修改處：更新現有地圖時，只放入 stars，不放 points
                to_update.append((data['stars'], name))

        # 執行新增
        if to_insert:

            print(f"🚀 正在新增 {len(to_insert)} 張新地圖...")
            insert_query = """
            INSERT INTO map_records 
            (difficulty, map_name, stars, runner, points, score, note, status) 
            VALUES (%s, %s, %s, '', %s, 0, '', 0)
            """
            cursor.executemany(insert_query, to_insert)

        # 執行更新
        if to_update:

            print(f"🔧 正在更新 {len(to_update)} 張地圖的星級 (不更動分數)...")
            # ★ 修改處：SQL 語句移除 points = %s
            update_query = """
            UPDATE map_records 
            SET stars = %s 
            WHERE map_name = %s
            """
            cursor.executemany(update_query, to_update)

        conn.commit()
        print("✅ 資料同步完成！")

    except Exception as e:
        print(f"❌ 發生錯誤: {e}")
        if conn: conn.rollback()
    finally:
        if conn:
            cursor.close()
            conn.close()
            print("🔒 資料庫連線已關閉")

if __name__ == "__main__":
    main()
